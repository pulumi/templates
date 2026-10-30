using Pulumi;
using Pulumi.Gcp.Storage;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Create a GCP resource (Storage Bucket)
    var bucket = new Bucket("my-bucket", new BucketArgs
    {
        Location = "US"
    });

    // Export the bucket's gs:// URL
    return new Dictionary<string, object?>
    {
        ["bucketName"] = bucket.Url
    };
});
