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

// Create a Cloud Function that returns some data.
const dataFunction = new gcp.cloudfunctions.Function("data-function", {
    sourceArchiveBucket: appBucket.name,
    sourceArchiveObject: appArchive.name,
    runtime: "nodejs16",
    entryPoint: "date",
    triggerHttp: true,
});

// Create an IAM member to invoke the function.
const invoker = new gcp.cloudfunctions.FunctionIamMember("data-function-invoker", {
    project: dataFunction.project,
    region: dataFunction.region,
    cloudFunction: dataFunction.name,
    role: "roles/cloudfunctions.invoker",
    member: "allUsers",
});

// Create a JSON configuration file for the website.
const siteConfig = new gcp.storage.BucketObject("site-config", {
    name: "config.json",
    source: dataFunction.httpsTriggerUrl
        .apply(url => new pulumi.asset.StringAsset(JSON.stringify({ api: url }))),
    contentType: "application/json",
    bucket: siteBucket.name,
});

// Export the URLs of the website and serverless endpoint.
export const siteURL = pulumi.interpolate`https://storage.googleapis.com/${siteBucket.name}/index.html`;
export const apiURL = dataFunction.httpsTriggerUrl;
