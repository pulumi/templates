using System.Collections.Generic;
using System.Text.Json;
using Pulumi;
using Gcp = Pulumi.Gcp;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var sitePath = config.Get("sitePath") ?? "./www";
    var appPath = config.Get("appPath") ?? "./app";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";

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

    // Create a Cloud Function that returns some data.
    var dataFunction = new Gcp.CloudFunctions.Function("data-function", new()
    {
        SourceArchiveBucket = appBucket.Name,
        SourceArchiveObject = appArchive.Name,
        Runtime = "dotnet6",
        EntryPoint = "App.Data",
        TriggerHttp = true,
    });

    // Create an IAM member to invoke the function.
    var invoker = new Gcp.CloudFunctions.FunctionIamMember("data-function-invoker", new()
    {
        Project = dataFunction.Project,
        Region = dataFunction.Region,
        CloudFunction = dataFunction.Name,
        Role = "roles/cloudfunctions.invoker",
        Member = "allUsers",
    });

    // Create a JSON configuration file for the website.
    var siteConfig = new Gcp.Storage.BucketObject("site-config", new()
    {
        Name = "config.json",
        ContentType = "application/json",
        Bucket = siteBucket.Name,
        Source = dataFunction.HttpsTriggerUrl.Apply(url => {
            var config = JsonSerializer.Serialize(new { api = url });
            return new Pulumi.StringAsset(config) as AssetOrArchive;
        }),
    });

    // Export the URLs of the website and serverless endpoint.
    return new Dictionary<string, object?>
    {
        ["siteURL"] = siteBucket.Name.Apply(name => $"https://storage.googleapis.com/{name}/index.html"),
        ["apiURL"] = dataFunction.HttpsTriggerUrl.Apply(url => url),
    };
});
