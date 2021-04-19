import * as storage from "@pulumi/google-native/storage/v1";
import * as pulumi from "@pulumi/pulumi";

const config = new pulumi.Config("google-native");
const project = config.require("project");
const bucketName = "pulumi-goog-native-ts-01";

// Create a Google Cloud resource (Storage Bucket)
const bucket = new storage.Bucket("my-bucket", {
    name:bucketName,
    bucket:bucketName,
    project: project,
});

// Export the bucket self-link
export const bucketSelfLink = bucket.selfLink;

