using Pulumi;
using Pulumi.GoogleNative.Storage.V1;

class MyStack : Stack
{
    public MyStack()
    {
        // Create a Google Cloud resource (Storage Bucket)
        var bucket = new Bucket("my-bucket");

        // Export the DNS name of the bucket
        this.BucketSelfLink = bucket.SelfLink;
    }

    [Output]
    public Output<string> BucketSelfLink { get; set; }
}
