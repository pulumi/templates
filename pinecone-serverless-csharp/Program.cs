using System.Collections.Generic;
using PineconeDatabase.Pinecone.Inputs;
using Pulumi;
using Pinecone = PineconeDatabase.Pinecone;

return await Deployment.RunAsync(() =>
{
    var myExampleIndex = new Pinecone.PineconeIndex("myExampleIndex", new Pinecone.PineconeIndexArgs
    {
        Name = "example-index-csharp",
        Metric= Pinecone.IndexMetric.Cosine,
        Spec= new Pinecone.Inputs.PineconeSpecArgs {
            Serverless= new PineconeServerlessSpecArgs{
                Cloud= Pinecone.ServerlessSpecCloud.Aws,
                Region= "us-west-2",
        }
    },
    });

    return new Dictionary<string, object?>
    {
        ["myPineconeIndexHost"] = myExampleIndex.Host
    };
});
