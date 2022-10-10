package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudfunctions"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		sitePath := "./www"
		if param := cfg.Get("sitePath"); param != "" {
			sitePath = param
		}
		appPath := "./app"
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

		// Create a Cloud Function that returns some data.
		dataFunction, err := cloudfunctions.NewFunction(ctx, "data-function", &cloudfunctions.FunctionArgs{
			SourceArchiveBucket: appBucket.Name,
			SourceArchiveObject: appArchive.Name,
			Runtime:             pulumi.String("go116"),
			EntryPoint:          pulumi.String("Data"),
			TriggerHttp:         pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create an IAM member to invoke the function.
		_, err = cloudfunctions.NewFunctionIamMember(ctx, "data-function-invoker", &cloudfunctions.FunctionIamMemberArgs{
			Project:       dataFunction.Project,
			Region:        dataFunction.Region,
			CloudFunction: dataFunction.Name,
			Role:          pulumi.String("roles/cloudfunctions.invoker"),
			Member:        pulumi.String("allUsers"),
		})
		if err != nil {
			return err
		}

		// Create a JSON configuration file for the website.
		_, err = storage.NewBucketObject(ctx, "site-config", &storage.BucketObjectArgs{
			Name:        pulumi.String("config.json"),
			Bucket:      siteBucket.Name,
			ContentType: pulumi.String("application/json"),
			Source: dataFunction.HttpsTriggerUrl.ApplyT(func(url string) pulumi.AssetOrArchiveOutput {
				config := fmt.Sprintf(`{ "api": "%s" }`, url)
				return pulumi.NewStringAsset(config).ToAssetOrArchiveOutput()
			}).(pulumi.AssetOrArchiveOutput),
		})
		if err != nil {
			return err
		}

		// Export the URLs of the website and serverless endpoint.
		ctx.Export("siteURL", pulumi.Sprintf("https://storage.googleapis.com/%s/index.html", siteBucket.Name))
		ctx.Export("apiURL", dataFunction.HttpsTriggerUrl)

		return nil
	})
}
