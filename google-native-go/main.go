package main

import (
	storage "github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "google-native")
		project := conf.Require("project")

		// Generate random bucket name
		randStr, err := random.NewRandomString(ctx, "suffix", &random.RandomStringArgs{
			Length:  pulumi.Int(5),
			Number:  pulumi.Bool(false),
			Special: pulumi.Bool(false),
			Upper:   pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		bucketName := pulumi.Sprintf("pulumi-goog-native-go-%s", randStr.Result)

		// Create a Google Cloud resource (Storage Bucket)
		bucket, err := storage.NewBucket(ctx, "bucket", &storage.BucketArgs{
			Name:    pulumi.StringPtr(bucketName),
			Bucket:  bucketName,
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
