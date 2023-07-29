"use strict";
const pulumi = require("@pulumi/pulumi");
const scaleway = require("@lbrlabs/pulumi-scaleway");

// Create a Scaleway resource (Object Bucket).
const bucket = new scaleway.ObjectBucket("my-bucket");

// Export the name of the bucket.
exports.bucketName = bucket.id;
