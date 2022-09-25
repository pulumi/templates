package main

import (
	"encoding/base64"

	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerservice"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/network"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or set default values
		cfg := config.New(ctx, "")
		azureLocation := config.Require(ctx, "azure-native:location")
		sshPubKey := cfg.Require("sshPubKey")
		mgmtGroup := cfg.Require("mgmtGroupId")
		prefixForDns, err := cfg.Try("prefixForDns")
		if err != nil {
			prefixForDns = "pulumi"
		}
		kubernetesVersion, err := cfg.Try("kubernetesVersion")
		if err != nil {
			kubernetesVersion = "1.24.3"
		}
		numWorkerNodes, err := cfg.TryInt("numWorkerNodes")
		if err != nil {
			numWorkerNodes = 3
		}
		nodeVmSize, err := cfg.Try("nodeVmSize")
		if err != nil {
			nodeVmSize = "Standard_DS2_v2"
		}

		// Create an Azure Resource Group
		resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{
			Location:          pulumi.String(azureLocation),
			ResourceGroupName: pulumi.String("aks-rg"),
		})
		if err != nil {
			return err
		}
		ctx.Export("resourceGroupName", resourceGroup.Name)

		// Create an Azure Virtual Network
		virtualNetwork, err := network.NewVirtualNetwork(ctx, "aks-network", &network.VirtualNetworkArgs{
			AddressSpace: &network.AddressSpaceArgs{
				AddressPrefixes: pulumi.StringArray{
					pulumi.String("10.0.0.0/16"),
				},
			},
			Location:           pulumi.String(azureLocation),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: pulumi.String("aks-network"),
		})
		if err != nil {
			return err
		}
		ctx.Export("networkName", virtualNetwork.Name)

		// Create three subnets in the virtual network
		subnet1, err := network.NewSubnet(ctx, "subnet-1", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.0.0/22"),
			ResourceGroupName:  resourceGroup.Name,
			SubnetName:         pulumi.String("aks-subnet-1"),
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}
		ctx.Export("subnet1Id", subnet1.ID())

		subnet2, err := network.NewSubnet(ctx, "subnet-2", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.4.0/22"),
			ResourceGroupName:  resourceGroup.Name,
			SubnetName:         pulumi.String("aks-subnet-2"),
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}
		ctx.Export("subnet2Id", subnet2.ID())

		subnet3, err := network.NewSubnet(ctx, "subnet-3", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.8.0/22"),
			ResourceGroupName:  resourceGroup.Name,
			SubnetName:         pulumi.String("aks-subnet-3"),
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}
		ctx.Export("subnet3Id", subnet3.ID())

		// Create a managed AKS cluster
		cluster, err := containerservice.NewManagedCluster(ctx, "aks-cluster", &containerservice.ManagedClusterArgs{
			AadProfile: &containerservice.ManagedClusterAADProfileArgs{
				EnableAzureRBAC: pulumi.Bool(true),
				Managed:         pulumi.Bool(true),
				AdminGroupObjectIDs: pulumi.StringArray{
					pulumi.String(mgmtGroup),
				},
			},
			// Use multiple agent/node pool profiles to distribute nodes across subnets
			AgentPoolProfiles: containerservice.ManagedClusterAgentPoolProfileArray{
				&containerservice.ManagedClusterAgentPoolProfileArgs{
					AvailabilityZones: pulumi.StringArray{
						pulumi.String("1"),
						pulumi.String("2"),
						pulumi.String("3"),
					},
					Count:              pulumi.Int(numWorkerNodes),
					EnableNodePublicIP: pulumi.Bool(false),
					Mode:               pulumi.String("System"),
					Name:               pulumi.String("systempool"),
					OsDiskSizeGB:       pulumi.Int(30),
					OsType:             pulumi.String("Linux"),
					VmSize:             pulumi.String(nodeVmSize),
					// Change next line for additional node pools to distribute across subnets
					VnetSubnetID: subnet1.ID(),
				},
			},
			// Change AuthorizedIPRanges to limit access to API server
			// Changing EnablePrivateCluster requires alternate access to API server (VPN or similar)
			ApiServerAccessProfile: containerservice.ManagedClusterAPIServerAccessProfileArgs{
				AuthorizedIPRanges: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
				EnablePrivateCluster: pulumi.Bool(false),
			},
			DnsPrefix:  pulumi.String(prefixForDns),
			EnableRBAC: pulumi.Bool(true),
			Identity: containerservice.ManagedClusterIdentityArgs{
				Type: containerservice.ResourceIdentityTypeSystemAssigned,
			},
			KubernetesVersion: pulumi.String(kubernetesVersion),
			LinuxProfile: &containerservice.ContainerServiceLinuxProfileArgs{
				AdminUsername: pulumi.String("azureuser"),
				Ssh: containerservice.ContainerServiceSshConfigurationArgs{
					PublicKeys: containerservice.ContainerServiceSshPublicKeyArray{
						containerservice.ContainerServiceSshPublicKeyArgs{
							KeyData: pulumi.String(sshPubKey),
						},
					},
				},
			},
			Location: pulumi.String(azureLocation),
			NetworkProfile: containerservice.ContainerServiceNetworkProfileArgs{
				NetworkPlugin: pulumi.String("azure"),
				NetworkPolicy: pulumi.String("azure"),
				ServiceCidr:   pulumi.String("10.96.0.0/16"),
				DnsServiceIP:  pulumi.String("10.96.0.10"),
			},
			ResourceGroupName: resourceGroup.Name,
			ResourceName:      pulumi.String("aks-cluster"),
		})
		if err != nil {
			return err
		}
		ctx.Export("clusterName", cluster.Name)

		// Build a Kubeconfig for accessing the cluster
		creds := containerservice.ListManagedClusterUserCredentialsOutput(ctx,
			containerservice.ListManagedClusterUserCredentialsOutputArgs{
				ResourceGroupName: resourceGroup.Name,
				ResourceName:      cluster.Name,
			})

		kubeconfig := creds.Kubeconfigs().Index(pulumi.Int(0)).Value().
			ApplyT(func(encoded string) string {
				kubeconfig, err := base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					return ""
				}
				return string(kubeconfig)
			})
		ctx.Export("kubeconfig", kubeconfig)

		return nil
	})
}
