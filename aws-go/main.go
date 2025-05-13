package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type infrastructure struct {
	bucket *s3.BucketV2
}

func createInfrastructure(ctx *pulumi.Context) (*infrastructure, error) {
	// Create an AWS resource (S3 Bucket) with tags.
	bucket, err := s3.NewBucketV2(ctx, "my-bucket", &s3.BucketV2Args{
		Tags: pulumi.StringMap{
			"Name": pulumi.String("My bucket"),
		},
	})
	if err != nil {
		return nil, err
	}

	return &infrastructure{
		bucket: bucket,
	}, nil
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		infra, err := createInfrastructure(ctx)
		if err != nil {
			return err
		}

		// Export the name of the bucket.
		ctx.Export("bucketName", infra.bucket.ID())
		return nil
	})
}
