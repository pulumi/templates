"use strict";
const pulumi = require("@pulumi/pulumi");
const aws = require("@pulumi/aws");
const awsx = require("@pulumi/awsx");

// Create an AWS resource (S3 Bucket)
const bucket = new aws.s3.BucketV2("my-bucket");

// Export the name of the bucket
exports.bucketName = bucket.id;
