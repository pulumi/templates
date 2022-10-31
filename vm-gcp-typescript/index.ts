import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";

// Import the program's configuration settings.
const config = new pulumi.Config();
const machineType = config.get("machineType") || "f1-micro";
const osImage = config.get("osImage") || "debian-11";
const instanceTag = config.get("instanceTag") || "webserver";
const servicePort = config.get("servicePort") || "80";

// Create a new network for the virtual machine.
const network = new gcp.compute.Network("network", {
    autoCreateSubnetworks: false,
});

// Create a subnet on the network.
const subnet = new gcp.compute.Subnetwork("subnet", {
    ipCidrRange: "10.0.1.0/24",
    network: network.id,
});

// Create a firewall allowing inbound access over ports 80 (for HTTP) and 22 (for SSH).
const firewall = new gcp.compute.Firewall("firewall", {
    network: network.selfLink,
    allows: [
        {
            protocol: "tcp",
            ports: [
                "22",
                servicePort,
            ],
        },
    ],
    direction: "INGRESS",
    sourceRanges: [
        "0.0.0.0/0",
    ],
    targetTags: [
        instanceTag,
    ],
});

// Define a script to be run when the VM starts up.
const metadataStartupScript = `#!/bin/bash
    echo '<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Hello, world!</title>
    </head>
    <body>
        <h1>Hello, world! 👋</h1>
        <p>Deployed with 💜 by <a href="https://pulumi.com/">Pulumi</a>.</p>
    </body>
    </html>' > index.html
    sudo python3 -m http.server ${servicePort} &`;

// Create the virtual machine.
const instance = new gcp.compute.Instance("instance", {
    machineType,
    bootDisk: {
        initializeParams: {
            image: osImage,
        },
    },
    networkInterfaces: [
        {
            network: network.id,
            subnetwork: subnet.id,
            accessConfigs: [
                {},
            ],
        },
    ],
    serviceAccount: {
        scopes: [
            "https://www.googleapis.com/auth/cloud-platform",
        ],
    },
    allowStoppingForUpdate: true,
    metadataStartupScript,
    tags: [
        instanceTag,
    ],
}, { dependsOn: firewall });

// Export the instance's name, public IP address, and HTTP URL.
export const name = instance.name;
export const ip = instance.networkInterfaces.apply(interfaces => {
    const configs = interfaces[0].accessConfigs;
    return configs && configs[0] && configs[0].natIp;
});
export const url = pulumi.interpolate`http://${ip}:${servicePort}`;
