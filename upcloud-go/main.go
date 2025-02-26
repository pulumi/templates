package main

import (
	"github.com/UpCloudLtd/pulumi-upcloud/sdk/go/upcloud"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Load configuration from Pulumi config
		cfg := config.New(ctx, "")

		objectStorageName := cfg.Get("object_storage_name")
		if objectStorageName == "" {
			objectStorageName = "bucket-example-objstov2"
		}

		region := cfg.Get("region")
		if region == "" {
			region = "europe-1"
		}

		bucketName := cfg.Get("bucket_name")
		if bucketName == "" {
			bucketName = "bucket"
		}

		// Create an UpCloud Managed Object Storage
		objectStorage, err := upcloud.NewManagedObjectStorage(ctx, "objectStorage", &upcloud.ManagedObjectStorageArgs{
			Name:             pulumi.String(objectStorageName),
			Region:           pulumi.String(region),
			ConfiguredStatus: pulumi.String("started"),
		})
		if err != nil {
			return err
		}

		// Create a Bucket inside the Object Storage
		bucket, err := upcloud.NewManagedObjectStorageBucket(ctx, "storageBucket", &upcloud.ManagedObjectStorageBucketArgs{
			ServiceUuid: objectStorage.ID(),
			Name:        pulumi.String(bucketName),
		})
		if err != nil {
			return err
		}

		// Export outputs
		ctx.Export("object_storage_uuid", objectStorage.ID())
		ctx.Export("bucket_name", bucket.Name)

		return nil
	})
}
