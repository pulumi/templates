using System.Collections.Generic;
using System.Threading.Tasks;

using Pulumi;
using Pulumi.Gcp.Storage;

class Program
{
    static Task<int> Main()
    {
        return Deployment.RunAsync(() => {

            // Create a GCP resource (Storage Bucket)
            var bucket = new Bucket("my-bucket");

            // Export the DNS name of the bucket
            return new Dictionary<string, object>
            {
                { "bucketName", bucket.Url },
            };
        });
    }
}
