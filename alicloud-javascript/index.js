"use strict";
const pulumi = require("@pulumi/pulumi");
const alicloud = require("@pulumi/alicloud");

// Create an AliCloud resource (OSS Bucket)
const bucket = new alicloud.oss.Bucket("my-bucket");

// Export the name of the bucket
exports.bucketName = bucket.id;
