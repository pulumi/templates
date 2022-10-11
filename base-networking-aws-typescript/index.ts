import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as awsx from "@pulumi/awsx";

// Grab some values from the Pulumi configuration (or use default values)
const config = new pulumi.Config();
const vpcNetworkCidr = config.get("vpcNetworkCidr") || "10.0.0.0/16";

// Create a new VPC
const vpc = new awsx.ec2.Vpc("vpc", {
    enableDnsHostnames: true,
    cidrBlock: vpcNetworkCidr,
    tags: {
        https: "allow",
    },
});

// Create a security group to allow HTTPS traffic
const securityGroup = new aws.ec2.SecurityGroup("securityGroup", {
    description: "Allow HTTPS inbound traffic",
    vpcId: vpc.vpcId,
    ingress: [{
        description: "HTTPS from VPC",
        fromPort: 443,
        toPort: 443,
        protocol: "tcp",
        cidrBlocks: [vpcNetworkCidr],
    }],
    egress: [{
        fromPort: 0,
        toPort: 0,
        protocol: "-1",
        cidrBlocks: ["0.0.0.0/0"],
    }],
    tags: {
        Name: "https",
    },
});

// Export some values for use elsewhere
export const vpcId = vpc.vpcId;
export const securityGroupId = securityGroup.id;
