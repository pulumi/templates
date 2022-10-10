using Pulumi;
using Pulumi.GoogleNative.Storage.V1;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Create a Google Cloud resource (Storage Bucket)
    var bucket = new Bucket("my-bucket");

    // Export the DNS name of the bucket
    return new Dictionary<string, object?>
    {
        ["bucketSelfLink"] = bucket.SelfLink
    };
});
