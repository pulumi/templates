import * as pulumi from "@pulumi/pulumi";
import * as alicloud from "@pulumi/alicloud";

// Create an AliCloud resource (OSS Bucket)
const bucket = new alicloud.oss.Bucket("my-bucket");

// Export the name of the bucket
export const bucketName = bucket.id;
