using Pulumi;
using Pulumi.AliCloud.Oss;

class MyStack : Stack
{
    public MyStack()
    {
        // Create an AliCloud resource (OSS Bucket)
        var bucket = new Bucket("my-bucket");

        // Export the name of the bucket
        this.BucketName = bucket.Id;
    }

    [Output]
    public Output<string> BucketName { get; set; }
}
