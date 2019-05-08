"use strict";
const pulumi = require("@pulumi/pulumi");
const linode = require("@pulumi/linode");

// Create a Linode resource (Linode Instance)
const instance = new linode.Instance("my-instance", {
    type: "g6-nanode-1",
    region: "us-east",
    image: "linode/ubuntu18.04",
});

// Export the Instance label of the instance
exports.instanceLabel = instance.label;
