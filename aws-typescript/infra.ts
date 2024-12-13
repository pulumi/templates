import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";

// Create an AWS resource (S3 Bucket) with tags.
export const bucket = new aws.s3.BucketV2("my-bucket", {
    tags: {
        "Name": "My bucket",
    },
});
