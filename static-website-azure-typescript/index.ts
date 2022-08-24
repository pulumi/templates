import * as pulumi from "@pulumi/pulumi";
import * as azure_native from "@pulumi/azure-native";
import * as synced_folder from "@pulumi/synced-folder";

const config = new pulumi.Config();
const path = config.get("path") || "./site";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";
const resourceGroup = new azure_native.resources.ResourceGroup("resource-group", {});
const account = new azure_native.storage.StorageAccount("account", {
    resourceGroupName: resourceGroup.name,
    kind: "StorageV2",
    sku: {
        name: "Standard_LRS",
    },
});
const website = new azure_native.storage.StorageAccountStaticWebsite("website", {
    resourceGroupName: resourceGroup.name,
    accountName: account.name,
    indexDocument: indexDocument,
    error404Document: errorDocument,
});
const syncedFolder = new synced_folder.AzureBlobFolder("synced-folder", {
    path: path,
    resourceGroupName: resourceGroup.name,
    storageAccountName: account.name,
    containerName: website.containerName,
});
const profile = new azure_native.cdn.Profile("profile", {
    resourceGroupName: resourceGroup.name,
    sku: {
        name: "Standard_Microsoft",
    },
});
const originHostname = account.primaryEndpoints.apply(endpoints => new URL(endpoints.web)).hostname;
const endpoint = new azure_native.cdn.Endpoint("endpoint", {
    resourceGroupName: resourceGroup.name,
    profileName: profile.name,
    isHttpAllowed: false,
    isHttpsAllowed: true,
    isCompressionEnabled: true,
    contentTypesToCompress: [
        "text/html",
        "text/css",
        "application/javascript",
        "application/json",
        "image/svg+xml",
        "font/woff",
        "font/woff2",
    ],
    originHostHeader: originHostname,
    origins: [{
        name: account.name,
        hostName: originHostname,
    }],
});
export const originURL = account.primaryEndpoints.apply(primaryEndpoints => primaryEndpoints.web);
export const cdnURL = pulumi.interpolate`https://${endpoint.hostName}`;
