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
for (let i = 0; i < 3; i++) {
    var subnet = new azure_native.network.Subnet(`subnet${i + 1}`, {
        addressPrefix: `10.0.${i * 4}.0/22`,
        resourceGroupName: resourceGroup.name,
        virtualNetworkName: virtualNetwork.name,
    })
};

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
