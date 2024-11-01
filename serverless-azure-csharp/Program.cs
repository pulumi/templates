using System.Collections.Generic;
using System.Text.Json;
using System.IO;
using Pulumi;
using Pulumi.Command.Local;
using AzureNative = Pulumi.AzureNative;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var sitePath = config.Get("sitePath") ?? "./www";
    var appPath = config.Get("appPath") ?? "./app";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";

    // Create a resource group for the website.
    var resourceGroup = new AzureNative.Resources.ResourceGroup("resource-group");

    // Create a blob storage account.
    var account = new AzureNative.Storage.StorageAccount("account", new()
    {
        ResourceGroupName = resourceGroup.Name,
        Kind = AzureNative.Storage.Kind.StorageV2,
        Sku = new AzureNative.Storage.Inputs.SkuArgs
        {
            Name = AzureNative.Storage.SkuName.Standard_LRS,
        },
    });

    // Create a storage container for the pages of the website.
    var website = new AzureNative.Storage.StorageAccountStaticWebsite("website", new()
    {
        AccountName = account.Name,
        ResourceGroupName = resourceGroup.Name,
        IndexDocument = indexDocument,
        Error404Document = errorDocument,
    });

    // Use a synced folder to manage the files of the website.
    var syncedFolder = new SyncedFolder.AzureBlobFolder("synced-folder", new()
    {
        Path = sitePath,
        ResourceGroupName = resourceGroup.Name,
        StorageAccountName = account.Name,
        ContainerName = website.ContainerName,
    });

    // Create a storage container for the serverless app.
    var appContainer = new AzureNative.Storage.BlobContainer("app-container", new()
    {
        AccountName = account.Name,
        ResourceGroupName = resourceGroup.Name,
        PublicAccess = AzureNative.Storage.PublicAccess.None,
    });

    // Compile the the app.
    var outputPath = "publish";
    var publishCommand = Run.Invoke(new()
    {
        Command = $"dotnet publish --output {outputPath}",
        Dir = appPath,
    });

    // Upload the serverless app to the storage container.
    var appBlob = new AzureNative.Storage.Blob("app-blob", new()
    {
        AccountName = account.Name,
        ResourceGroupName = resourceGroup.Name,
        ContainerName = appContainer.Name,
        Source = publishCommand.Apply(result => new FileArchive(Path.Combine(appPath, outputPath)) as AssetOrArchive),
    });

    // Create a shared access signature to give the Function App access to the code.
    var signature = AzureNative.Storage.ListStorageAccountServiceSAS.Invoke(new()
    {
        ResourceGroupName = resourceGroup.Name,
        AccountName = account.Name,
        Protocols = AzureNative.Storage.HttpProtocol.Https,
        SharedAccessStartTime = "2022-01-01",
        SharedAccessExpiryTime = "2030-01-01",
        Resource = AzureNative.Storage.SignedResource.C,
        Permissions = AzureNative.Storage.Permissions.R,
        ContentType = "application/json",
        CacheControl = "max-age=5",
        ContentDisposition = "inline",
        ContentEncoding = "deflate",
        CanonicalizedResource = Output.Tuple(account.Name, appContainer.Name).Apply(values => $"/blob/{values.Item1}/{values.Item2}"),
    }).Apply(result => result.ServiceSasToken);

    // Create an App Service plan for the Function App.
    var plan = new AzureNative.Web.AppServicePlan("plan", new()
    {
        ResourceGroupName = resourceGroup.Name,
        Sku = new AzureNative.Web.Inputs.SkuDescriptionArgs
        {
            Name = "Y1",
            Tier = "Dynamic",
        },
    });

    // Create the Function App.
    var app = new AzureNative.Web.WebApp("app", new()
    {
        ResourceGroupName = resourceGroup.Name,
        ServerFarmId = plan.Id,
        Kind = "FunctionApp",
        SiteConfig = new AzureNative.Web.Inputs.SiteConfigArgs
        {
            NetFrameworkVersion = "v8.0",
            DetailedErrorLoggingEnabled = true,
            HttpLoggingEnabled = true,
            AppSettings = new[]
            {
                new AzureNative.Web.Inputs.NameValuePairArgs
                {
                    Name = "FUNCTIONS_WORKER_RUNTIME",
                    Value = "dotnet",
                },
                new AzureNative.Web.Inputs.NameValuePairArgs
                {
                    Name = "FUNCTIONS_EXTENSION_VERSION",
                    Value = "~4",
                },
                new AzureNative.Web.Inputs.NameValuePairArgs
                {
                    Name = "WEBSITE_RUN_FROM_PACKAGE",
                    Value = Output.Tuple(account.Name, appContainer.Name, appBlob.Name, signature).Apply(values =>
                    {
                        var accountName = values.Item1;
                        var containerName = values.Item2;
                        var blobName = values.Item3;
                        var token = values.Item4;
                        return $"https://{accountName}.blob.core.windows.net/{containerName}/{blobName}?{token}";
                    }),
                },
            },
            Cors = new AzureNative.Web.Inputs.CorsSettingsArgs
            {
                AllowedOrigins = new[]
                {
                    "*",
                },
            },
        },
    });

    // Create a JSON configuration file for the website.
    var siteConfig = new AzureNative.Storage.Blob("config.json", new()
    {
        AccountName = account.Name,
        ResourceGroupName = resourceGroup.Name,
        ContainerName = website.ContainerName,
        ContentType = "application/json",
        Source = app.DefaultHostName.Apply(hostname => {
            var config = JsonSerializer.Serialize(new { api = $"https://{hostname}/api" });
            return new Pulumi.StringAsset(config) as AssetOrArchive;
        }),
    });

    // Export the URLs of the website and serverless endpoint.
    return new Dictionary<string, object?>
    {
        ["siteURL"] = account.PrimaryEndpoints.Apply(primaryEndpoints => primaryEndpoints.Web),
        ["apiURL"] = app.DefaultHostName.Apply(defaultHostName => $"https://{defaultHostName}/api"),
    };
});
