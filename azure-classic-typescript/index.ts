import * as pulumi from "@pulumi/pulumi";
import * as azure from "@pulumi/azure";

// Create an Azure Resource Group
const resourceGroup = new azure.core.ResourceGroup("resourceGroup");

// Create an Azure resource (Storage Account)
const account = new azure.storage.Account("storage", {
    // The location for the storage account will be derived automatically from the resource group.
    resourceGroupName: resourceGroup.name,
    accountTier: "Standard",
    accountReplicationType: "LRS",
});

// Export the connection string for the storage account
export const connectionString = account.primaryConnectionString;
