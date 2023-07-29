package main

import (
	"github.com/lbrlabs/pulumi-scaleway/sdk/go/scaleway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a Scaleway resource (Object Bucket).
		bucket, err := scaleway.NewObjectBucket(ctx, "my-bucket", nil)
		if err != nil {
			return err
		}

		// Export the name of the bucket.
		ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
