"""A minimal Pinecone Python Pulumi program"""

import pulumi
import pinecone_pulumi as pinecone

my_pinecone_index = pinecone.PineconeIndex(
    "myPineconeIndex",
    name="example-index-python",
    metric=pinecone.IndexMetric.COSINE,
    spec={
        "serverless": {
            "cloud": pinecone.ServerlessSpecCloud.AWS,
            "region": "us-west-2",
        },
    },
)

pulumi.export("host", my_pinecone_index.host)
