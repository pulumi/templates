import * as storage from "@pulumi/google-native/storage/v1";

// Create a Google Cloud resource (Storage Bucket)
const bucket = new storage.Bucket("my-bucket");

// Export the bucket self-link
export const bucketSelfLink = bucket.selfLink;

