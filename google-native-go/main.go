package main

import (
	storage "github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	const bucketName = "pulumi-goog-native-bucket-go-01"
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "google-native")
		project := conf.Require("project")
		// Create a Google Cloud resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, "bucket", &storage.BucketArgs{
			Name:    pulumi.StringPtr(bucketName),
			Bucket:  pulumi.String(bucketName),
			Project: pulumi.String(project),
		})
		if err != nil {
			return err
		}
		// Export the bucket self-link
		ctx.Export("bucketSelfLink", bucket.SelfLink)

		return nil
	})
}

