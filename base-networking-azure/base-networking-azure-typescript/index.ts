import * as pulumi from "@pulumi/pulumi";
import * as azure_native from "@pulumi/azure-native";

const config = new pulumi.Config();
const azureNativeLocation = config.require("azureNativeLocation");
const mgmtGroupId = config.require("mgmtGroupId");
const resourceGroup = new azure_native.resources.ResourceGroup("resourceGroup", {
    location: azureNativeLocation,
    resourceGroupName: "rg",
});
const virtualNetwork = new azure_native.network.VirtualNetwork("virtualNetwork", {
    addressSpace: {
        addressPrefixes: ["10.0.0.0/16"],
    },
    location: azureNativeLocation,
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: "vnet",
});
const subnet1 = new azure_native.network.Subnet("subnet1", {
    addressPrefix: "10.0.0.0/22",
    name: "subnet-1",
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: virtualNetwork.name,
});
const subnet2 = new azure_native.network.Subnet("subnet2", {
    addressPrefix: "10.0.4.0/22",
    name: "subnet-2",
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: virtualNetwork.name,
});
const subnet3 = new azure_native.network.Subnet("subnet3", {
    addressPrefix: "10.0.8.0/22",
    name: "subnet-3",
    resourceGroupName: resourceGroup.name,
    virtualNetworkName: virtualNetwork.name,
});
