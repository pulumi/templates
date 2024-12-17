using Pulumi;
using Pulumi.Aws.S3;
using System.Collections.Generic;

return await Deployment.RunAsync(Deploy.Infra);

public static class Deploy
{
    public static Dictionary<string, object?> Infra()
    {
        // Create an AWS resource (S3 Bucket) with tags.
        var bucket = new BucketV2("my-bucket", new()
        {
            Tags =
            {
                { "Name", "My bucket" },
            },
        });

        // Export the name of the bucket.
        return new Dictionary<string, object?>
        {
            ["bucketName"] = bucket.Id
        };
    }
}
