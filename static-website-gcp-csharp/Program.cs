using System.Collections.Generic;
using Pulumi;
using Gcp = Pulumi.Gcp;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var path = config.Get("path") ?? "./www";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";

    // Create a storage bucket and configure it as a website.
    var bucket = new Gcp.Storage.Bucket("bucket", new()
    {
        Location = "US",
        Website = new Gcp.Storage.Inputs.BucketWebsiteArgs
        {
            MainPageSuffix = indexDocument,
            NotFoundPage = errorDocument,
        },
    });

    // Create an IAM binding to allow public read access to the bucket.
    var bucketIamBinding = new Gcp.Storage.BucketIAMBinding("bucket-iam-binding", new()
    {
        Bucket = bucket.Name,
        Role = "roles/storage.objectViewer",
        Members = new[]
        {
            "allUsers",
        },
    });

    // Use a synced folder to manage the files of the website.
    var syncedFolder = new SyncedFolder.GoogleCloudFolder("synced-folder", new()
    {
        Path = path,
        BucketName = bucket.Name,
    });

    // Enable the storage bucket as a CDN.
    var backendBucket = new Gcp.Compute.BackendBucket("backend-bucket", new()
    {
        BucketName = bucket.Name,
        EnableCdn = true,
    });

    // Provision a global IP address for the CDN.
    var ip = new Gcp.Compute.GlobalAddress("ip");

    // Create a URLMap to route requests to the storage bucket.
    var urlMap = new Gcp.Compute.URLMap("url-map", new()
    {
        DefaultService = backendBucket.SelfLink,
    });

    // Create an HTTP proxy to route requests to the URLMap.
    var httpProxy = new Gcp.Compute.TargetHttpProxy("http-proxy", new()
    {
        UrlMap = urlMap.SelfLink,
    });

    // Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
    var httpForwardingRule = new Gcp.Compute.GlobalForwardingRule("http-forwarding-rule", new()
    {
        IpAddress = ip.Address,
        IpProtocol = "TCP",
        PortRange = "80",
        Target = httpProxy.SelfLink,
    });

    // Export the URLs and hostnames of the bucket and CDN.
    return new Dictionary<string, object?>
    {
        ["originURL"] = bucket.Name.Apply(name => $"https://storage.googleapis.com/{name}/index.html"),
        ["originHostname"] = bucket.Name.Apply(name => $"storage.googleapis.com/{name}"),
        ["cdnURL"] = ip.Address.Apply(address => $"http://{address}"),
        ["cdnHostname"] = ip.Address,
    };
});
