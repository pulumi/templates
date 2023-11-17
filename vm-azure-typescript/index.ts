import * as pulumi from "@pulumi/pulumi";
import * as resources from "@pulumi/azure-native/resources";
import * as network from "@pulumi/azure-native/network";
import * as compute from "@pulumi/azure-native/compute";
import * as random from "@pulumi/random";
import * as tls from "@pulumi/tls";

// Import the program's configuration settings
const config = new pulumi.Config();
const vmName = config.get("vmName") || "my-server";
const vmSize = config.get("vmSize") || "Standard_A1_v2";
const osImage = config.get("osImage") || "Debian:debian-11:11:latest";
const adminUsername = config.get("adminUsername") || "pulumiuser";
const servicePort = config.get("servicePort") || "80";

const [ osImagePublisher, osImageOffer, osImageSku, osImageVersion ] = osImage.split(":");

// Create an SSH key
const sshKey = new tls.PrivateKey("ssh-key", {
    algorithm: "RSA",
    rsaBits: 4096,
});

// Create a resource group
const resourceGroup = new resources.ResourceGroup("resource-group");

// Create a virtual network
const virtualNetwork = new network.VirtualNetwork("network", {
    resourceGroupName: resourceGroup.name,
    addressSpace: {
        addressPrefixes: ["10.0.0.0/16"],
    },
    subnets: [
        {
            name: `${vmName}-subnet`,
            addressPrefix: "10.0.1.0/24",
        },
    ],
});

// Use a random string to give the VM a unique DNS name
var domainNameLabel = new random.RandomString("domain-label", {
    length: 8,
    upper: false,
    special: false,
}).result.apply(result => `${vmName}-${result}`);

// Create a public IP address for the VM
const publicIp = new network.PublicIPAddress("public-ip", {
    resourceGroupName: resourceGroup.name,
    publicIPAllocationMethod: network.IPAllocationMethod.Dynamic,
    dnsSettings: {
        domainNameLabel: domainNameLabel,
    },
});

// Create a security group allowing inbound access over ports 80 (for HTTP) and 22 (for SSH)
const securityGroup = new network.NetworkSecurityGroup("security-group", {
    resourceGroupName: resourceGroup.name,
    securityRules: [
        {
            name: `${vmName}-securityrule`,
            priority: 1000,
            direction: network.AccessRuleDirection.Inbound,
            access: "Allow",
            protocol: "Tcp",
            sourcePortRange: "*",
            sourceAddressPrefix: "*",
            destinationAddressPrefix: "*",
            destinationPortRanges: [
                servicePort,
                "22",
            ],
        }
    ]
});

// Create a network interface with the virtual network, IP address, and security group
const networkInterface = new network.NetworkInterface("network-interface", {
    resourceGroupName: resourceGroup.name,
    networkSecurityGroup: {
        id: securityGroup.id,
    },
    ipConfigurations: [{
        name: `${vmName}-ipconfiguration`,
        privateIPAllocationMethod: network.IPAllocationMethod.Dynamic,
        subnet: {
            id: virtualNetwork.subnets.apply(subnets => subnets![0].id!),
        },
        publicIPAddress: {
            id: publicIp.id,
        },
    }],
});

// Define a script to be run when the VM starts up
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

// Create the virtual machine
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
        computerName: vmName,
        adminUsername: adminUsername,
        customData: Buffer.from(initScript).toString("base64"),
        linuxConfiguration: {
            disablePasswordAuthentication: true,
            ssh: {
                publicKeys: [
                    {
                        keyData: sshKey.publicKeyOpenssh,
                        path: `/home/${adminUsername}/.ssh/authorized_keys`,
                    },
                ],
            },
        },
    },
    storageProfile: {
        osDisk: {
            name: `${vmName}-osdisk`,
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

// Once the machine is created, fetch its IP address and DNS hostname
const vmAddress = vm.id.apply(_ => network.getPublicIPAddressOutput({
    resourceGroupName: resourceGroup.name,
    publicIpAddressName: publicIp.name,
}));

// Export the VM's hostname, public IP address, HTTP URL, and SSH private key
export const ip = vmAddress.ipAddress;
export const hostname = vmAddress.dnsSettings?.apply(settings => settings?.fqdn);
export const url = hostname?.apply(name => `http://${name}:${servicePort}`);
export const privatekey = sshKey.privateKeyOpenssh;
