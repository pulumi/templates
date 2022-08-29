using System.Collections.Generic;
using Pulumi;
using Gcp = Pulumi.Gcp;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
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

    var backendBucket = new Gcp.Compute.BackendBucket("backend-bucket", new()
    {
        BucketName = bucket.Name,
        EnableCdn = true,
    });

    var ip = new Gcp.Compute.GlobalAddress("ip");

    var urlMap = new Gcp.Compute.URLMap("url-map", new()
    {
        DefaultService = backendBucket.SelfLink,
    });

    var httpProxy = new Gcp.Compute.TargetHttpProxy("http-proxy", new()
    {
        UrlMap = urlMap.SelfLink,
    });

    var httpForwardingRule = new Gcp.Compute.GlobalForwardingRule("http-forwarding-rule", new()
    {
        IpAddress = ip.Address,
        IpProtocol = "TCP",
        PortRange = "80",
        Target = httpProxy.SelfLink,
    });

    return new Dictionary<string, object?>
    {
        ["originURL"] = bucket.Name.Apply(name => $"https://storage.googleapis.com/{name}/index.html"),
        ["originHostname"] = bucket.Name.Apply(name => $"storage.googleapis.com/{name}"),
        ["cdnURL"] = ip.Address.Apply(address => $"http://{address}"),
        ["cdnHostname"] = ip.Address,
    };
});

