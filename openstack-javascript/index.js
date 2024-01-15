"use strict";
const pulumi = require("@pulumi/pulumi");
const os = require("@pulumi/openstack");

// Create an OpenStack resource (Compute Instance)
const instance = new os.compute.Instance("test", {
	flavorName: "s1-2",
	imageName: "Ubuntu 22.04",
});

// Export the IP of the instance
exports.instanceIP = instance.accessIpV4;
