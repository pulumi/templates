using Pulumi;
using Pulumi.AliCloud.Oss;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Create an AliCloud resource (OSS Bucket)
    var bucket = new Bucket("my-bucket");

    // Export the name of the bucket
    return new Dictionary<string, object?>
    {
        ["bucketName"] = bucket.Id
    };
});
