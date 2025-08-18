import * as pulumi from "@pulumi/pulumi";
import * as upcloud from "@upcloud/pulumi-upcloud";

// Load Pulumi config values
const config = new pulumi.Config();

const objectStorageName = config.require("object_storage_name");
const region = config.require("region");
const bucketName = config.require("bucket_name");

// Create an UpCloud Managed Object Storage
const objectStorage = new upcloud.ManagedObjectStorage("objectStorage", {
    name: objectStorageName,
    region: region,
    configuredStatus: "started"
});

// Create a Bucket inside the Object Storage
const bucket = new upcloud.ManagedObjectStorageBucket("storageBucket", {
    serviceUuid: objectStorage.id,
    name: bucketName
});

// Export outputs
export const objectStorageUuid = objectStorage.id;
export const bucketNameOutput = bucket.name;
