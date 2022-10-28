import * as pulumi from "@pulumi/pulumi";
import * as resources from "@pulumi/azure-native/resources";
import * as network from "@pulumi/azure-native/network";
import * as compute from "@pulumi/azure-native/compute";
import * as random from "@pulumi/random";

const config = new pulumi.Config();
const adminUsername = config.get("adminUsername") || "pulumiUser";
const vmName = config.get("vmName") || "my-server";
const vmSize = config.get("vmSize") || "Standard_A0";
const osImage = config.get("osImage") || "Debian:debian-11:11:latest";
const servicePort = config.getNumber("servicePort") || 80;
const sshPublicKey = config.require("sshPublicKey");

const [ osImagePublisher, osImageOffer, osImageSku, osImageVersion ] = osImage.split(":");

const resourceGroup = new resources.ResourceGroup("resource-group");

const virtualNetwork = new network.VirtualNetwork("virtual-network", {
    resourceGroupName: resourceGroup.name,
    addressSpace: {
        addressPrefixes: ["10.0.0.0/16"],
    },
    subnets: [
        {
            name: "default",
            addressPrefix: "10.0.1.0/24",
        },
    ],
});

// Use a random string to give the server a unique DNS name.
var dnsName = new random.RandomString("dns-name", {
    length: 8,
    special: false,
}).result.apply(result => `${vmName}-${result.toLowerCase()}`);

const publicIp = new network.PublicIPAddress("public-ip", {
    resourceGroupName: resourceGroup.name,
    publicIPAllocationMethod: network.IPAllocationMethod.Dynamic,
    dnsSettings: {
        domainNameLabel: dnsName,
    },
});

const securityGroup = new network.NetworkSecurityGroup("security-group", {
    resourceGroupName: resourceGroup.name,
    securityRules: [
        {
            name: "web",
            priority: 1000,
            direction: network.AccessRuleDirection.Inbound,
            access: "Allow",
            protocol: "Tcp",
            sourcePortRange: "*",
            sourceAddressPrefix: "*",
            destinationPortRanges: [
                "22",
                servicePort.toString(),
            ],
            destinationAddressPrefix: "*",
        }
    ]
});

const networkInterface = new network.NetworkInterface("network-interface", {
    resourceGroupName: resourceGroup.name,
    ipConfigurations: [{
        name: "othername",
        subnet: virtualNetwork.subnets.apply(subnet => subnet![0]),
        privateIPAllocationMethod: network.IPAllocationMethod.Dynamic,
        publicIPAddress: {
            id: publicIp.id,
        },
    }],
    networkSecurityGroup: {
        id: securityGroup.id,
    },
});

const initScript = `#!/bin/bash
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

const vm = new compute.VirtualMachine("vm", {
    resourceGroupName: resourceGroup.name,
    networkProfile: {
        networkInterfaces: [
            {
                id: networkInterface.id,
                primary: true,
            },
        ],
    },
    hardwareProfile: {
        vmSize: vmSize,
    },
    osProfile: {
        computerName: "somename",
        adminUsername: adminUsername,
        customData: Buffer.from(initScript).toString("base64"),
        linuxConfiguration: {
            disablePasswordAuthentication: true,
            ssh: {
                publicKeys: [
                    {
                        keyData: sshPublicKey,
                        path: `/home/${adminUsername}/.ssh/authorized_keys`,
                    },
                ],
            },
        },
    },

    storageProfile: {
        osDisk: {
            name: "myosdisk1",
            createOption: compute.DiskCreateOption.FromImage,
        },
        imageReference: {
            publisher: osImagePublisher,
            offer: osImageOffer,
            sku: osImageSku,
            version: osImageVersion,
        },
    },
});

const address = vm.id.apply(_ => network.getPublicIPAddressOutput({
    resourceGroupName: resourceGroup.name,
    publicIpAddressName: publicIp.name,
}));

export const ip = address.ipAddress;
export const hostname = address.dnsSettings?.apply(settings => settings?.fqdn);
export const url = hostname?.apply(name => `http://${name}`);
