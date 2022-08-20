import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";
import * as synced_folder from "@pulumi/synced-folder";

const config = new pulumi.Config();
const gcpProject = config.get("gcpProject") || "pulumi-development";
const path = config.get("path") || "./site";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";
const bucket = new gcp.storage.Bucket("bucket", {
    location: "US",
    website: {
        mainPageSuffix: indexDocument,
        notFoundPage: errorDocument,
    },
});
const bucketIamBinding = new gcp.storage.BucketIAMBinding("bucket-iam-binding", {
    bucket: bucket.name,
    role: "roles/storage.objectViewer",
    members: ["allUsers"],
});
const syncedFolder = new synced_folder.GoogleCloudFolder("synced-folder", {
    path: path,
    bucketName: bucket.name,
});
export const url = pulumi.interpolate`https://storage.googleapis.com/${bucket.name}/index.html`;
