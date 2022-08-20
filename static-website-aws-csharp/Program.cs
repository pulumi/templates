using System.Collections.Generic;
using Pulumi;
using Aws = Pulumi.Aws;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var awsRegion = config.Get("awsRegion") ?? "us-west-2";
    var path = config.Get("path") ?? "./site";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";
    var bucket = new Aws.S3.Bucket("bucket", new()
    {
        Acl = "public-read",
        Website = new Aws.S3.Inputs.BucketWebsiteArgs
        {
            IndexDocument = indexDocument,
            ErrorDocument = errorDocument,
        },
    });

    var bucketFolder = new SyncedFolder.S3BucketFolder("bucket-folder", new()
    {
        Path = path,
        BucketName = bucket.BucketName,
        Acl = "public-read",
    });

    return new Dictionary<string, object?>
    {
        ["url"] = bucket.WebsiteEndpoint.Apply(websiteEndpoint => $"http://{websiteEndpoint}"),
    };
});

