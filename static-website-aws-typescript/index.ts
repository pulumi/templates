import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as synced_folder from "@pulumi/synced-folder";

const config = new pulumi.Config();
const awsRegion = config.get("awsRegion") || "us-west-2";
const path = config.get("path") || "./site";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";
const bucket = new aws.s3.Bucket("bucket", {
    acl: "public-read",
    website: {
        indexDocument: indexDocument,
        errorDocument: errorDocument,
    },
});
const bucketFolder = new synced_folder.S3BucketFolder("bucket-folder", {
    path: path,
    bucketName: bucket.bucket,
    acl: "public-read",
});
export const url = pulumi.interpolate`http://${bucket.websiteEndpoint}`;
