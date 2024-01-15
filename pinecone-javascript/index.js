"use strict";
const pulumi = require("@pulumi/pulumi");
const pinecone = require("@pinecone-database/pulumi");

const myExampleIndex = new pinecone.PineconeIndex("my-example-index", {
    name: "my-example-index",
    metric: pinecone.IndexMetric.Cosine,
    spec: {
        serverless: {
            cloud: pinecone.ServerlessSpecCloud.Aws,
            region: "us-west-2",
        }
    }
});

exports.host = myExampleIndex.host;
