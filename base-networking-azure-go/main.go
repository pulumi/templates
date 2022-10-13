package main

import (
	"fmt"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a new resource group
		resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", nil)
		if err != nil {
			return err
		}

		// Create a new virtual network
		virtualNetwork, err := network.NewVirtualNetwork(ctx, "virtualNetwork", &network.VirtualNetworkArgs{
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
			ResourceGroupName: resourceGroup.Name,
		})
		if err != nil {
			return err
		}

		// Create three subnets in the virtual network
		for i := 0; i < 3; i++ {
			_, err = network.NewSubnet(ctx, fmt.Sprintf("subnet%d", i+1), &network.SubnetArgs{
				AddressPrefix:      pulumi.String(fmt.Sprintf("10.0.%d.0/22", i*4)),
				ResourceGroupName:  resourceGroup.Name,
				VirtualNetworkName: virtualNetwork.Name,
			})
			if err != nil {
				return err
			}
		}

		// Create a security group to allow HTTPS traffic
		securityGroup, err := network.NewNetworkSecurityGroup(ctx, "securityGroup", &network.NetworkSecurityGroupArgs{
			ResourceGroupName: resourceGroup.Name,
			SecurityRules: network.SecurityRuleTypeArray{
				&network.SecurityRuleTypeArgs{
					Access:                   pulumi.String("Allow"),
					DestinationAddressPrefix: pulumi.String("*"),
					DestinationPortRange:     pulumi.String("443"),
					Direction:                pulumi.String("Inbound"),
					Name:                     pulumi.String("allow-inbound-https"),
					Priority:                 pulumi.Int(200),
					Protocol:                 pulumi.String("TCP"),
					SourceAddressPrefix:      pulumi.String("*"),
					SourcePortRange:          pulumi.String("*"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export some values for use elsewhere
		ctx.Export("virtualNetworkId", virtualNetwork.ID())
		ctx.Export("securityGroupId", securityGroup.ID())
		return nil
	})
}
