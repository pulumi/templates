using Pulumi;
using Pulumi.GoogleNative.Storage.V1;

class MyStack : Stack
{
    public MyStack()
    {
        var config = new Config("google-native");
        var project = config.Require("project");
        var bucketName = "pulumi-goog-native-bucket-cs-01";
        // Create a Google Cloud resource (Storage Bucket)
        var bucket = new Bucket("my-bucket", new BucketArgs
        {
            Name = bucketName,
            Bucket = bucketName,
            Project = project,
        });

        // Export the DNS name of the bucket
        this.BucketSelfLink = bucket.SelfLink;
    }

    [Output]
    public Output<string> BucketSelfLink { get; set; }
}
