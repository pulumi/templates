import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";
import * as synced_folder from "@pulumi/synced-folder";

const config = new pulumi.Config();
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
const backendBucket = new gcp.compute.BackendBucket("backend-bucket", {
    bucketName: bucket.name,
    enableCdn: true,
});
const ip = new gcp.compute.GlobalAddress("ip", {});
const urlMap = new gcp.compute.URLMap("url-map", {defaultService: backendBucket.selfLink});
const httpProxy = new gcp.compute.TargetHttpProxy("http-proxy", {urlMap: urlMap.selfLink});
const httpForwardingRule = new gcp.compute.GlobalForwardingRule("http-forwarding-rule", {
    ipAddress: ip.address,
    ipProtocol: "TCP",
    portRange: "80",
    target: httpProxy.selfLink,
});
export const originURL = pulumi.interpolate`https://storage.googleapis.com/${bucket.name}/index.html`;
export const cdnURL = pulumi.interpolate`http://${ip.address}`;
