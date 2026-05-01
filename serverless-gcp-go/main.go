package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudfunctionsv2"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/storage"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		sitePath := "www"
		if param := cfg.Get("sitePath"); param != "" {
			sitePath = param
		}
		appPath := "app"
		if param := cfg.Get("appPath"); param != "" {
			appPath = param
		}
		indexDocument := "index.html"
		if param := cfg.Get("indexDocument"); param != "" {
			indexDocument = param
		}
		errorDocument := "error.html"
		if param := cfg.Get("errorDocument"); param != "" {
			errorDocument = param
		}

		gcpCfg := config.New(ctx, "gcp")
		region := gcpCfg.Get("region")
		if region == "" {
			region = "us-central1"
		}

		// Create a storage bucket and configure it as a website.
		siteBucket, err := storage.NewBucket(ctx, "site-bucket", &storage.BucketArgs{
			Location: pulumi.String("US"),
			Website: &storage.BucketWebsiteArgs{
				MainPageSuffix: pulumi.String(indexDocument),
				NotFoundPage:   pulumi.String(errorDocument),
			},
		})
		if err != nil {
			return err
		}

		// Create an IAM binding to allow public read access to the bucket.
		_, err = storage.NewBucketIAMBinding(ctx, "site-bucket-iam-binding", &storage.BucketIAMBindingArgs{
			Bucket: siteBucket.Name,
			Role:   pulumi.String("roles/storage.objectViewer"),
			Members: pulumi.StringArray{
				pulumi.String("allUsers"),
			},
		})
		if err != nil {
			return err
		}

		// Use a synced folder to manage the files of the website.
		_, err = synced.NewGoogleCloudFolder(ctx, "synced-folder", &synced.GoogleCloudFolderArgs{
			Path:       pulumi.String(sitePath),
			BucketName: siteBucket.Name,
		})
		if err != nil {
			return err
		}

		// Create another storage bucket for the serverless app.
		appBucket, err := storage.NewBucket(ctx, "app-bucket", &storage.BucketArgs{
			Location: pulumi.String("US"),
		})
		if err != nil {
			return err
		}

		// Upload the serverless app to the storage bucket.
		appArchive, err := storage.NewBucketObject(ctx, "app-archive", &storage.BucketObjectArgs{
			Bucket: appBucket.Name,
			Source: pulumi.NewFileArchive(appPath),
		})
		if err != nil {
			return err
		}

		// Create a Cloud Function (Gen 2) that returns some data.
		dataFunction, err := cloudfunctionsv2.NewFunction(ctx, "data-function", &cloudfunctionsv2.FunctionArgs{
			Location: pulumi.String(region),
			BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
				Runtime:    pulumi.String("go122"),
				EntryPoint: pulumi.String("Data"),
				Source: &cloudfunctionsv2.FunctionBuildConfigSourceArgs{
					StorageSource: &cloudfunctionsv2.FunctionBuildConfigSourceStorageSourceArgs{
						Bucket: appBucket.Name,
						Object: appArchive.Name,
					},
				},
			},
			ServiceConfig: &cloudfunctionsv2.FunctionServiceConfigArgs{
				AvailableMemory: pulumi.String("256M"),
				TimeoutSeconds:  pulumi.Int(60),
			},
		})
		if err != nil {
			return err
		}

		// Allow public, unauthenticated invocations of the underlying Cloud Run service.
		_, err = cloudrun.NewIamMember(ctx, "data-function-invoker", &cloudrun.IamMemberArgs{
			Location: dataFunction.Location,
			Service:  dataFunction.Name,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		// Create a JSON configuration file for the website.
		_, err = storage.NewBucketObject(ctx, "site-config", &storage.BucketObjectArgs{
			Name:        pulumi.String("config.json"),
			Bucket:      siteBucket.Name,
			ContentType: pulumi.String("application/json"),
			Source: dataFunction.Url.ApplyT(func(url string) pulumi.AssetOrArchiveOutput {
				config := fmt.Sprintf(`{ "api": "%s" }`, url)
				return pulumi.NewStringAsset(config).ToAssetOrArchiveOutput()
			}).(pulumi.AssetOrArchiveOutput),
		})
		if err != nil {
			return err
		}

		// Export the URLs of the website and serverless endpoint.
		ctx.Export("siteURL", pulumi.Sprintf("https://storage.googleapis.com/%s/index.html", siteBucket.Name))
		ctx.Export("apiURL", dataFunction.Url)

		return nil
	})
}
