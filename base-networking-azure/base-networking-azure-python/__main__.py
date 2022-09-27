import pulumi
import pulumi_azure_native as azure_native

config = pulumi.Config()
azure_native_location = config.require("azureNativeLocation")
mgmt_group_id = config.require("mgmtGroupId")
resource_group = azure_native.resources.ResourceGroup("resourceGroup",
    location=azure_native_location,
    resource_group_name="rg")
virtual_network = azure_native.network.VirtualNetwork("virtualNetwork",
    address_space=azure_native.network.AddressSpaceArgs(
        address_prefixes=["10.0.0.0/16"],
    ),
    location=azure_native_location,
    resource_group_name=resource_group.name,
    virtual_network_name="vnet")
subnet1 = azure_native.network.Subnet("subnet1",
    address_prefix="10.0.0.0/22",
    name="subnet-1",
    resource_group_name=resource_group.name,
    virtual_network_name=virtual_network.name)
subnet2 = azure_native.network.Subnet("subnet2",
    address_prefix="10.0.4.0/22",
    name="subnet-2",
    resource_group_name=resource_group.name,
    virtual_network_name=virtual_network.name)
subnet3 = azure_native.network.Subnet("subnet3",
    address_prefix="10.0.8.0/22",
    name="subnet-3",
    resource_group_name=resource_group.name,
    virtual_network_name=virtual_network.name)
