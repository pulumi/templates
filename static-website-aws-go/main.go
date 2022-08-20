package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		awsRegion := "us-west-2"
		if param := cfg.Get("awsRegion"); param != "" {
			awsRegion = param
		}
		path := "./site"
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
		bucket, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{
			Acl: pulumi.String("public-read"),
			Website: &s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String(indexDocument),
				ErrorDocument: pulumi.String(errorDocument),
			},
		})
		if err != nil {
			return err
		}
		_, err = synced - folder.NewS3BucketFolder(ctx, "bucket-folder", &synced-folder.S3BucketFolderArgs{
			Path:       pulumi.String(path),
			BucketName: bucket.Bucket,
			Acl:        pulumi.String("public-read"),
		})
		if err != nil {
			return err
		}
		ctx.Export("url", bucket.WebsiteEndpoint.ApplyT(func(websiteEndpoint string) (string, error) {
			return fmt.Sprintf("http://%v", websiteEndpoint), nil
		}).(pulumi.StringOutput))
		return nil
	})
}
