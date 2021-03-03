"use strict";
const pulumi = require("@pulumi/pulumi");
const azure = require("@pulumi/azure-native");

// Create an Azure Resource Group
const resourceGroup = new azure.resources.ResourceGroup("resourceGroup");

// Create an Azure resource (Storage Account)
const storageAccount = new azure.storage.StorageAccount("sa", {
    resourceGroupName: resourceGroup.name,
    sku: {
        name: "Standard_LRS",
    },
    kind: "StorageV2",
});

// Export the primary key of the Storage Account
const storageAccountKeys = pulumi.all([resourceGroup.name, storageAccount.name]).apply(([resourceGroupName, accountName]) =>
    azure.storage.listStorageAccountKeys({ resourceGroupName, accountName }));

// Export the primary storage key for the storage account
exports.primaryStorageKey = storageAccountKeys.keys[0].value;
