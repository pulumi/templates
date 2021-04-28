using Pulumi;
using Pulumi.GoogleNative.Storage.V1;
using Pulumi.Random;

class MyStack : Stack
{
    public MyStack()
    {
        var config = new Config("google-native");
        var project = config.Require("project");

        // Generate random bucket name
        var suffix = new RandomString("suffix", new RandomStringArgs
        {
            Length = 5,
            Number = false,
            Special = false,
            Upper = false,
        });
        var bucketName = Output.Format($"pulumi-goog-native-bucket-cs-{suffix.Result}");

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
