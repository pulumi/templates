using System;
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

    var profile = new AzureNative.Cdn.Profile("profile", new()
    {
        ResourceGroupName = resourceGroup.Name,
        Sku = new AzureNative.Cdn.Inputs.SkuArgs
        {
            Name = "Standard_Microsoft",
        },
    });

    var originHostname = account.PrimaryEndpoints.Apply(endpoints => new Uri(endpoints.Web).Host);

    var endpoint = new AzureNative.Cdn.Endpoint("endpoint", new()
    {
        ResourceGroupName = resourceGroup.Name,
        ProfileName = profile.Name,
        IsHttpAllowed = false,
        IsHttpsAllowed = true,
        IsCompressionEnabled = true,
        ContentTypesToCompress = new[]
        {
            "text/html",
            "text/css",
            "application/javascript",
            "application/json",
            "image/svg+xml",
            "font/woff",
            "font/woff2",
        },
        OriginHostHeader = originHostname,
        Origins = new[]
        {
            new AzureNative.Cdn.Inputs.DeepCreatedOriginArgs
            {
                Name = account.Name,
                HostName = originHostname,
            },
        },
    });

    return new Dictionary<string, object?>
    {
        ["originURL"] = account.PrimaryEndpoints.Apply(primaryEndpoints => primaryEndpoints.Web),
        ["originHostname"] = originHostname,
        ["cdnURL"] = endpoint.HostName.Apply(hostName => $"https://{hostName}"),
        ["cdnHostname"] = endpoint.HostName,
    };
});
