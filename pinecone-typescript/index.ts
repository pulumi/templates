import * as pulumi from "@pulumi/pulumi";
import * as pinecone from "@pinecone-database/pulumi";

const myExampleIndex = new pinecone.PineconeIndex("my-example-index", {
    name: "example-index-ts",
    metric: pinecone.IndexMetric.Cosine,
    spec: {
        serverless: {
            cloud: pinecone.ServerlessSpecCloud.Aws,
            region: "us-west-2",
        },
    },
});
export const host = myExampleIndex.host;
