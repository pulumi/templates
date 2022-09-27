using System.Collections.Generic;
using Pulumi;
using AzureNative = Pulumi.AzureNative;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var azureNativeLocation = config.Require("azureNativeLocation");
    var mgmtGroupId = config.Require("mgmtGroupId");
    var resourceGroup = new AzureNative.Resources.ResourceGroup("resourceGroup", new()
    {
        Location = azureNativeLocation,
        ResourceGroupName = "rg",
    });

    var virtualNetwork = new AzureNative.Network.VirtualNetwork("virtualNetwork", new()
    {
        AddressSpace = new AzureNative.Network.Inputs.AddressSpaceArgs
        {
            AddressPrefixes = new[]
            {
                "10.0.0.0/16",
            },
        },
        Location = azureNativeLocation,
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = "vnet",
    });

    var subnet1 = new AzureNative.Network.Subnet("subnet1", new()
    {
        AddressPrefix = "10.0.0.0/22",
        Name = "subnet-1",
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = virtualNetwork.Name,
    });

    var subnet2 = new AzureNative.Network.Subnet("subnet2", new()
    {
        AddressPrefix = "10.0.4.0/22",
        Name = "subnet-2",
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = virtualNetwork.Name,
    });

    var subnet3 = new AzureNative.Network.Subnet("subnet3", new()
    {
        AddressPrefix = "10.0.8.0/22",
        Name = "subnet-3",
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = virtualNetwork.Name,
    });

});

