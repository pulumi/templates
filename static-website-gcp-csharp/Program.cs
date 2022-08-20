using System.Collections.Generic;
using Pulumi;
using Gcp = Pulumi.Gcp;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var gcpProject = config.Get("gcpProject") ?? "pulumi-development";
    var path = config.Get("path") ?? "./site";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";
    var bucket = new Gcp.Storage.Bucket("bucket", new()
    {
        Location = "US",
        Website = new Gcp.Storage.Inputs.BucketWebsiteArgs
        {
            MainPageSuffix = indexDocument,
            NotFoundPage = errorDocument,
        },
    });

    var bucketIamBinding = new Gcp.Storage.BucketIAMBinding("bucket-iam-binding", new()
    {
        Bucket = bucket.Name,
        Role = "roles/storage.objectViewer",
        Members = new[]
        {
            "allUsers",
        },
    });

    var syncedFolder = new SyncedFolder.GoogleCloudFolder("synced-folder", new()
    {
        Path = path,
        BucketName = bucket.Name,
    });

    return new Dictionary<string, object?>
    {
        ["url"] = bucket.Name.Apply(name => $"https://storage.googleapis.com/{name}/index.html"),
    };
});

