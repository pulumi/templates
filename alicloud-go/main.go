package main

import (
	"github.com/pulumi/pulumi-alicloud/sdk/v2/go/alicloud/oss"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an AliCloud resource (OSS Bucket)
		bucket, err := oss.NewBucket(ctx, "my-bucket", nil)
		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		return nil
	})
}
