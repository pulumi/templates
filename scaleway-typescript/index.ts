import * as pulumi from "@pulumi/pulumi";
import * as scaleway from "@lbrlabs/pulumi-scaleway";

// Create a Scaleway resource (Object Bucket).
const bucket = new scaleway.ObjectBucket("my-bucket",);

// Export the name of the bucket.
export const bucketName = bucket.id;
