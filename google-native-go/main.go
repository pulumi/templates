package main

import (
	"github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a Google Cloud resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, "bucket", nil)
		if err != nil {
			return err
		}

		// Export the bucket self-link
		ctx.Export("bucketSelfLink", bucket.SelfLink)

		return nil
	})
}
