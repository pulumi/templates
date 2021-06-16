package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	ppkg "github.com/pulumi/pulumi-package/sdk/go/pulumipackage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Create a bucket and expose a website index document.
		bucket, err := s3.NewBucket(ctx, "pluginServer", &s3.BucketArgs{
			Website: s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String("index.html"),
			},
			ForceDestroy: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Set the access policy for the bucket so all objects are readable.
		if _, err := s3.NewBucketPolicy(ctx, "bucketPolicy", &s3.BucketPolicyArgs{
			Bucket: bucket.ID(),
			Policy: pulumi.Any(map[string]interface{}{
				"Version": "2012-10-17",
				"Statement": []map[string]interface{}{
					{
						"Effect":    "Allow",
						"Principal": "*",
						"Action": []interface{}{
							"s3:GetObject",
						},
						"Resource": []interface{}{
							pulumi.Sprintf("arn:aws:s3:::%s/*", bucket.ID()), // policy refers to bucket name explicitly
						},
					},
				},
			}),
		}); err != nil {
			return err
		}

		language := "go"
		packageName := "${PROJECT}"

		pkg, err := ppkg.NewPackage(ctx, "mypkg", &ppkg.PackageArgs{
			Language:                    &language,
			Name:                        &packageName,
			ServerBucketName:            bucket.Bucket,
			ServerBucketWebsiteEndpoint: bucket.WebsiteEndpoint,
		})
		if err != nil {
			return err
		}

		ctx.Export("releases", pkg.Releases)

		return nil
	})
}
