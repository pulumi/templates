import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";
import * as synced from "@pulumi/synced-folder";

// Import the program's configuration settings.
const config = new pulumi.Config();
const sitePath = config.get("path") || "./www";
const appPath = config.get("appPath") || "./app";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";

// Create a storage bucket and configure it as a website.
const siteBucket = new gcp.storage.Bucket("site-bucket", {
    location: "US",
    website: {
        mainPageSuffix: indexDocument,
        notFoundPage: errorDocument,
    },
});

// Create an IAM binding to allow public read access to the bucket.
const siteBucketIAMBinding = new gcp.storage.BucketIAMBinding("site-bucket-iam-binding", {
    bucket: siteBucket.name,
    role: "roles/storage.objectViewer",
    members: ["allUsers"],
});

// Use a synced folder to manage the files of the website.
const syncedFolder = new synced.GoogleCloudFolder("synced-folder", {
    path: sitePath,
    bucketName: siteBucket.name,
});

// Create another storage bucket for the serverless app.
const appBucket = new gcp.storage.Bucket("app-bucket", {
    location: "US",
});

// Upload the serverless app to the storage bucket.
const appArchive = new gcp.storage.BucketObject("app-archive", {
    bucket: appBucket.name,
    source: new pulumi.asset.FileArchive(appPath),
});

// Create a Cloud Function (Gen 2) that returns some data.
const dataFunction = new gcp.cloudfunctionsv2.Function("data-function", {
    location: gcp.config.region || "us-central1",
    buildConfig: {
        runtime: "nodejs22",
        entryPoint: "date",
        source: {
            storageSource: {
                bucket: appBucket.name,
                object: appArchive.name,
            },
        },
    },
    serviceConfig: {
        availableMemory: "256M",
        timeoutSeconds: 60,
    },
});

// Allow public, unauthenticated invocations of the underlying Cloud Run service.
const invoker = new gcp.cloudrun.IamMember("data-function-invoker", {
    location: dataFunction.location,
    service: dataFunction.name,
    role: "roles/run.invoker",
    member: "allUsers",
});

// Create a JSON configuration file for the website.
const siteConfig = new gcp.storage.BucketObject("site-config", {
    name: "config.json",
    source: dataFunction.url.apply(url =>
        new pulumi.asset.StringAsset(JSON.stringify({ api: url })),
    ),
    contentType: "application/json",
    bucket: siteBucket.name,
});

// Export the URLs of the website and serverless endpoint.
export const siteURL = pulumi.interpolate`https://storage.googleapis.com/${siteBucket.name}/index.html`;
export const apiURL = dataFunction.url;
