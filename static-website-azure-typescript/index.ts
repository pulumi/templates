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
export const url = account.primaryEndpoints.apply(primaryEndpoints => primaryEndpoints.web);
