using System.Collections.Generic;
using System.Text.Json;
using Pulumi;
using Gcp = Pulumi.Gcp;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var sitePath = config.Get("sitePath") ?? "www";
    var appPath = config.Get("appPath") ?? "app";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";
    var region = new Config("gcp").Get("region") ?? "us-central1";

    // Create a storage bucket and configure it as a website.
    var siteBucket = new Gcp.Storage.Bucket("site-bucket", new()
    {
        Location = "US",
        Website = new Gcp.Storage.Inputs.BucketWebsiteArgs
        {
            MainPageSuffix = indexDocument,
            NotFoundPage = errorDocument,
        },
    });

    // Create an IAM binding to allow public read access to the bucket.
    var siteBucketIAMBinding = new Gcp.Storage.BucketIAMBinding("site-bucket-iam-binding", new()
    {
        Bucket = siteBucket.Name,
        Role = "roles/storage.objectViewer",
        Members = new[]
        {
            "allUsers",
        },
    });

    // Use a synced folder to manage the files of the website.
    var syncedFolder = new SyncedFolder.GoogleCloudFolder("synced-folder", new()
    {
        Path = sitePath,
        BucketName = siteBucket.Name,
    });

    // Create another storage bucket for the serverless app.
    var appBucket = new Gcp.Storage.Bucket("app-bucket", new()
    {
        Location = "US",
    });

    // Upload the serverless app to the storage bucket.
    var appArchive = new Gcp.Storage.BucketObject("app-archive", new()
    {
        Bucket = appBucket.Name,
        Source = new Pulumi.FileArchive(appPath),
    });

    // Create a Cloud Function (Gen 2) that returns some data.
    var dataFunction = new Gcp.CloudFunctionsV2.Function("data-function", new()
    {
        Location = region,
        BuildConfig = new Gcp.CloudFunctionsV2.Inputs.FunctionBuildConfigArgs
        {
            Runtime = "dotnet8",
            EntryPoint = "App.Data",
            Source = new Gcp.CloudFunctionsV2.Inputs.FunctionBuildConfigSourceArgs
            {
                StorageSource = new Gcp.CloudFunctionsV2.Inputs.FunctionBuildConfigSourceStorageSourceArgs
                {
                    Bucket = appBucket.Name,
                    Object = appArchive.Name,
                },
            },
        },
        ServiceConfig = new Gcp.CloudFunctionsV2.Inputs.FunctionServiceConfigArgs
        {
            AvailableMemory = "256M",
            TimeoutSeconds = 60,
        },
    });

    // Allow public, unauthenticated invocations of the underlying Cloud Run service.
    var invoker = new Gcp.CloudRun.IamMember("data-function-invoker", new()
    {
        Location = dataFunction.Location,
        Service = dataFunction.Name,
        Role = "roles/run.invoker",
        Member = "allUsers",
    });

    // Create a JSON configuration file for the website.
    var siteConfig = new Gcp.Storage.BucketObject("site-config", new()
    {
        Name = "config.json",
        ContentType = "application/json",
        Bucket = siteBucket.Name,
        Source = dataFunction.Url.Apply(url => {
            var config = JsonSerializer.Serialize(new { api = url });
            return new Pulumi.StringAsset(config) as AssetOrArchive;
        }),
    });

    // Export the URLs of the website and serverless endpoint.
    return new Dictionary<string, object?>
    {
        ["siteURL"] = siteBucket.Name.Apply(name => $"https://storage.googleapis.com/{name}/index.html"),
        ["apiURL"] = dataFunction.Url,
    };
});
