import * as storage from "@pulumi/google-native/storage/v1";
import * as pulumi from "@pulumi/pulumi";
import * as random from "@pulumi/random";

const randomString = new random.RandomString("name", {
    upper: false,
    number: false,
    special: false,
    length: 5,
});

const config = new pulumi.Config("google-native");
const project = config.require("project");
const bucketName = pulumi.interpolate `pulumi-goog-native-ts-${randomString.result}`;

// Create a Google Cloud resource (Storage Bucket)
const bucket = new storage.Bucket("my-bucket", {
    name:bucketName,
    project: project,
});

// Export the bucket self-link
export const bucketSelfLink = bucket.selfLink;

