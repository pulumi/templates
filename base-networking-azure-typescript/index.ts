import * as azure_native from "@pulumi/azure-native";

// Create a new resource group
const resourceGroup = new azure_native.resources.ResourceGroup("resourceGroup", {});

// Create a new virtual network
const virtualNetwork = new azure_native.network.VirtualNetwork("virtualNetwork", {
    addressSpace: {
        addressPrefixes: ["10.0.0.0/16"],
    },
    resourceGroupName: resourceGroup.name,
});

// Create three subnets in the new virtual network
const subnet1 = new azure_native.network.Subnet("subnet1", {
    addressPrefix: "10.0.0.0/22",
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: virtualNetwork.name,
});
const subnet2 = new azure_native.network.Subnet("subnet2", {
    addressPrefix: "10.0.4.0/22",
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: virtualNetwork.name,
});
const subnet3 = new azure_native.network.Subnet("subnet3", {
    addressPrefix: "10.0.8.0/22",
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: virtualNetwork.name,
});

// Create a security group to allow HTTPS traffic
const securityGroup = new azure_native.network.NetworkSecurityGroup("securityGroup", {
    resourceGroupName: resourceGroup.name,
    securityRules: [{
        access: "Allow",
        destinationAddressPrefix: "*",
        destinationPortRange: "443",
        direction: "Inbound",
        name: "allow-inbound-https",
        priority: 200,
        protocol: "TCP",
        sourceAddressPrefix: "*",
        sourcePortRange: "*",
    }],
});

// Export some values for use elsewhere
export const virtualNetworkId = virtualNetwork.id;
export const securityGroupId = securityGroup.id;
