import * as pulumi from "@pulumi/pulumi";
import * as azure_native from "@pulumi/azure-native";
import * as synced_folder from "@pulumi/synced-folder";

// Import the program's configuration settings.
const config = new pulumi.Config();
const path = config.get("path") || "www";
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

// Pull the hostname out of the storage-account endpoint.
const originHostname = account.primaryEndpoints.apply(endpoints => new URL(endpoints.web)).hostname;

// Export the URL and hostname of the storage account's website.
export const originURL = account.primaryEndpoints.apply(endpoints => endpoints.web);
export { originHostname };
