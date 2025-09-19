package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/storage"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		path := "./www"
		if param := cfg.Get("path"); param != "" {
			path = param
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
		bucket, err := storage.NewBucket(ctx, "bucket", &storage.BucketArgs{
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
		_, err = storage.NewBucketIAMBinding(ctx, "bucket-iam-binding", &storage.BucketIAMBindingArgs{
			Bucket: bucket.Name,
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
			Path:       pulumi.String(path),
			BucketName: bucket.Name,
		})
		if err != nil {
			return err
		}

		// Enable the storage bucket as a CDN.
		backendBucket, err := compute.NewBackendBucket(ctx, "backend-bucket", &compute.BackendBucketArgs{
			BucketName: bucket.Name,
			EnableCdn:  pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Provision a global IP address for the CDN.
		ip, err := compute.NewGlobalAddress(ctx, "ip", nil)
		if err != nil {
			return err
		}

		// Create a URLMap to route requests to the storage bucket.
		urlMap, err := compute.NewURLMap(ctx, "url-map", &compute.URLMapArgs{
			DefaultService: backendBucket.SelfLink,
		})
		if err != nil {
			return err
		}

		// Create an HTTP proxy to route requests to the URLMap.
		httpProxy, err := compute.NewTargetHttpProxy(ctx, "http-proxy", &compute.TargetHttpProxyArgs{
			UrlMap: urlMap.SelfLink,
		})
		if err != nil {
			return err
		}

		// Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
		_, err = compute.NewGlobalForwardingRule(ctx, "http-forwarding-rule", &compute.GlobalForwardingRuleArgs{
			IpAddress:  ip.Address,
			IpProtocol: pulumi.String("TCP"),
			PortRange:  pulumi.String("80"),
			Target:     httpProxy.SelfLink,
		})
		if err != nil {
			return err
		}

		// Export the URLs and hostnames of the bucket and CDN.
		ctx.Export("originURL", pulumi.Sprintf("https://storage.googleapis.com/%v/index.html", bucket.Name))
		ctx.Export("originHostname", pulumi.Sprintf("storage.googleapis.com/%v", bucket.Name))
		ctx.Export("cdnURL", pulumi.Sprintf("http://%v", ip.Address))
		ctx.Export("cdnHostname", ip.Address)
		return nil
	})
}
