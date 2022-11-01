import * as pulumi from "@pulumi/pulumi";
import * as azure_native from "@pulumi/azure-native";
import * as synced_folder from "@pulumi/synced-folder";

// Import the program's configuration settings.
const config = new pulumi.Config();
const path = config.get("path") || "./www";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";

// Create a resource group for the website.
const resourceGroup = new azure_native.resources.ResourceGroup("resource-group", {});

// Create a blob storage account.
const account = new azure_native.storage.StorageAccount("account", {
    resourceGroupName: resourceGroup.name,
    kind: "StorageV2",
    sku: {
        name: "Standard_LRS",
    },
});

// Configure the storage account as a website.
const website = new azure_native.storage.StorageAccountStaticWebsite("website", {
    resourceGroupName: resourceGroup.name,
    accountName: account.name,
    indexDocument: indexDocument,
    error404Document: errorDocument,
});

// Use a synced folder to manage the files of the website.
const syncedFolder = new synced_folder.AzureBlobFolder("synced-folder", {
    path: path,
    resourceGroupName: resourceGroup.name,
    storageAccountName: account.name,
    containerName: website.containerName,
});

// Create a CDN profile.
const profile = new azure_native.cdn.Profile("profile", {
    resourceGroupName: resourceGroup.name,
    sku: {
        name: "Standard_Microsoft",
    },
});

// Pull the hostname out of the storage-account endpoint.
const originHostname = account.primaryEndpoints.apply(endpoints => new URL(endpoints.web)).hostname;

// Create a CDN endpoint to distribute and cache the website.
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

// Export the URLs and hostnames of the storage account and CDN.
export const originURL = account.primaryEndpoints.apply(endpoints => endpoints.web);
export { originHostname };
export const cdnURL = pulumi.interpolate`https://${endpoint.hostName}`;
export const cdnHostname = endpoint.hostName;
