package main

import (
	network "github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	resources "github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		azureNativeLocation := cfg.Require("azureNativeLocation")
		mgmtGroupId := cfg.Require("mgmtGroupId")
		resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
			Location:          pulumi.String(azureNativeLocation),
			ResourceGroupName: pulumi.String("rg"),
		})
		if err != nil {
			return err
		}
		virtualNetwork, err := network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
			Location:           pulumi.String(azureNativeLocation),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: pulumi.String("vnet"),
		})
		if err != nil {
			return err
		}
		_, err = network.NewSubnet(ctx, "subnet1", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.0.0/22"),
			Name:               pulumi.String("subnet-1"),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}
		_, err = network.NewSubnet(ctx, "subnet2", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.4.0/22"),
			Name:               pulumi.String("subnet-2"),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}
		_, err = network.NewSubnet(ctx, "subnet3", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.8.0/22"),
			Name:               pulumi.String("subnet-3"),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}
		return nil
	})
}
