using System.Collections.Generic;
using System.Threading.Tasks;

using Pulumi;
using Pulumi.Aws.S3;

class Program
{
    static Task<int> Main()
    {
        return Deployment.RunAsync(() => {

            // Create an AWS resource (S3 Bucket)
            var bucket = new Bucket("my-bucket");

            // Export the name of the bucket
            return new Dictionary<string, object>
            {
                { "bucketName", bucket.Id },
            };
        });
    }
}
