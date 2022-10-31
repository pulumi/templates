import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";

const config = new pulumi.Config();
const machineType = config.get("machineType") || "f1-micro";
const osImage = config.get("osImage") || "debian-11";
const instanceTag = config.get("instanceTag") || "webserver";
const servicePort = config.get("servicePort") || "80";

// const address = new gcp.compute.Address("address");

const network = new gcp.compute.Network("network", {
    autoCreateSubnetworks: false,
});

const subnet = new gcp.compute.Subnetwork("subnet", {
    ipCidrRange: "10.0.1.0/24",
    network: network.id,
});

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

const instance = new gcp.compute.Instance("instance", {

    // https://cloud.google.com/compute/docs/machine-types
    machineType,
    bootDisk: {
        initializeParams: {
            // https://gcloud-compute.com/images.html
            image: osImage,
        },
    },
    networkInterfaces: [
        {
            network: network.id,
            subnetwork: subnet.id,
            accessConfigs: [
                {
                    // natIp: address.address,
                },
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

export const name = instance.name;
export const ip = instance.networkInterfaces.apply(interfaces => {
    const configs = interfaces[0].accessConfigs;
    return configs && configs[0] && configs[0].natIp;
});
export const url = pulumi.interpolate`http://${ip}:${servicePort}`;

// gcloud compute ssh $(pulumi stack output name) --zone $(pulumi config get gcp:zone) --project $(pulumi config get gcp:project)