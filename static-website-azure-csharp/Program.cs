using System.Collections.Generic;
using Pulumi;
using AzureNative = Pulumi.AzureNative;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var path = config.Get("path") ?? "./site";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";
    var resourceGroup = new AzureNative.Resources.ResourceGroup("resource-group");

    var account = new AzureNative.Storage.StorageAccount("account", new()
    {
        ResourceGroupName = resourceGroup.Name,
        Kind = "StorageV2",
        Sku = new AzureNative.Storage.Inputs.SkuArgs
        {
            Name = "Standard_LRS",
        },
    });

    var website = new AzureNative.Storage.StorageAccountStaticWebsite("website", new()
    {
        ResourceGroupName = resourceGroup.Name,
        AccountName = account.Name,
        IndexDocument = indexDocument,
        Error404Document = errorDocument,
    });

    var syncedFolder = new SyncedFolder.AzureBlobFolder("synced-folder", new()
    {
        Path = path,
        ResourceGroupName = resourceGroup.Name,
        StorageAccountName = account.Name,
        ContainerName = website.ContainerName,
    });

    return new Dictionary<string, object?>
    {
        ["url"] = account.PrimaryEndpoints.Apply(primaryEndpoints => primaryEndpoints.Web),
    };
});

