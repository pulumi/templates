import * as pulumi from "@pulumi/pulumi";
import * as azure from "@pulumi/azure-native";
import * as synced from "@pulumi/synced-folder";

// Import the program's configuration settings.
const config = new pulumi.Config();
const sitePath = config.get("sitePath") || "./www";
const appPath = config.get("appPath") || "./app";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";

// Create a resource group for the website.
const resourceGroup = new azure.resources.ResourceGroup("resource-group", {});

// Create a blob storage account.
const account = new azure.storage.StorageAccount("account", {
    resourceGroupName: resourceGroup.name,
    kind: azure.storage.Kind.StorageV2,
    sku: {
        name: azure.storage.SkuName.Standard_LRS,
    },
});

// Create a storage container for the pages of the website.
const website = new azure.storage.StorageAccountStaticWebsite("website", {
    accountName: account.name,
    resourceGroupName: resourceGroup.name,
    indexDocument: indexDocument,
    error404Document: errorDocument,
});

// Use a synced folder to manage the files of the website.
const syncedFolder = new synced.AzureBlobFolder("synced-folder", {
    path: sitePath,
    resourceGroupName: resourceGroup.name,
    storageAccountName: account.name,
    containerName: website.containerName,
});

// Create a storage container for the serverless app.
const appContainer = new azure.storage.BlobContainer("app-container", {
    accountName: account.name,
    resourceGroupName: resourceGroup.name,
    publicAccess: azure.storage.PublicAccess.None,
});

// Upload the serverless app to the storage container.
const appBlob = new azure.storage.Blob("app-blob", {
    accountName: account.name,
    resourceGroupName: resourceGroup.name,
    containerName: appContainer.name,
    source: new pulumi.asset.FileArchive(appPath),
});

// Create a shared access signature to give the Function App access to the code.
const signature = azure.storage.listStorageAccountServiceSASOutput({
    resourceGroupName: resourceGroup.name,
    accountName: account.name,
    protocols: azure.storage.HttpProtocol.Https,
    sharedAccessStartTime: "2022-01-01",
    sharedAccessExpiryTime: "2030-01-01",
    resource: azure.storage.SignedResource.C,
    permissions: azure.storage.Permissions.R,
    contentType: "application/json",
    cacheControl: "max-age=5",
    contentDisposition: "inline",
    contentEncoding: "deflate",
    canonicalizedResource: pulumi.interpolate`/blob/${account.name}/${appContainer.name}`,
});

// Create an App Service plan for the Function App.
const plan = new azure.web.AppServicePlan("plan", {
    resourceGroupName: resourceGroup.name,
    sku: {
        name: "Y1",
        tier: "Dynamic",
    },
});

// Create the Function App.
const functionApp = new azure.web.WebApp("function-app", {
    resourceGroupName: resourceGroup.name,
    serverFarmId: plan.id,
    kind: "FunctionApp",
    siteConfig: {
        appSettings: [
            {
                name: "FUNCTIONS_WORKER_RUNTIME",
                value: "node",
            },
            {
                name: "WEBSITE_NODE_DEFAULT_VERSION",
                value: "~14",
            },
            {
                name: "FUNCTIONS_EXTENSION_VERSION",
                value: "~3",
            },
            {
                name: "WEBSITE_RUN_FROM_PACKAGE",
                value: pulumi.all([account.name, appContainer.name, appBlob.name, signature])
                    .apply(([accountName, containerName, blobName, signature]) => `https://${accountName}.blob.core.windows.net/${containerName}/${blobName}?${signature.serviceSasToken}`),
            },
        ],
        cors: {
            allowedOrigins: [
                "*"
            ],
        },
    },
});

// Create a JSON configuration file for the website.
const configFile = new azure.storage.Blob("config.json", {
    source: functionApp.defaultHostName
        .apply(host => new pulumi.asset.StringAsset(JSON.stringify({ api: `https://${host}/api` }))),
    contentType: "application/json",
    accountName: account.name,
    resourceGroupName: resourceGroup.name,
    containerName: website.containerName,
});

// Export the URLs of the website and serverless endpoint.
export const siteURL = account.primaryEndpoints.apply(primaryEndpoints => primaryEndpoints.web);
export const apiURL = pulumi.interpolate`https://${functionApp.defaultHostName}/api`;
