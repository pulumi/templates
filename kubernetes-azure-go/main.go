package main

import (
	"encoding/base64"

	"github.com/pulumi/pulumi-azure-native-sdk/containerservice/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/network/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or set default values
		cfg := config.New(ctx, "")
		prefixForDns, err := cfg.Try("prefixForDns")
		if err != nil {
			prefixForDns = "pulumi"
		}
		kubernetesVersion, err := cfg.Try("kubernetesVersion")
		if err != nil {
			kubernetesVersion = "1.32"
		}
		numWorkerNodes, err := cfg.TryInt("numWorkerNodes")
		if err != nil {
			numWorkerNodes = 3
		}
		nodeVmSize, err := cfg.Try("nodeVmSize")
		if err != nil {
			nodeVmSize = "Standard_DS2_v2"
		}
		// The next two configuration values are required (no default can be provided)
		sshPubKey := cfg.Require("sshPubKey")
		mgmtGroup := cfg.Require("mgmtGroupId")

		// Create an Azure Resource Group
		resourceGroup, err := resources.NewResourceGroup(ctx, "resourceGroup", &resources.ResourceGroupArgs{})
		if err != nil {
			return err
		}

		// Create an Azure Virtual Network
		virtualNetwork, err := network.NewVirtualNetwork(ctx, "aks-network", &network.VirtualNetworkArgs{
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
		subnet1, err := network.NewSubnet(ctx, "subnet-1", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.0.0/22"),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}

		subnet2, err := network.NewSubnet(ctx, "subnet-2", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.4.0/22"),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}

		subnet3, err := network.NewSubnet(ctx, "subnet-3", &network.SubnetArgs{
			AddressPrefix:      pulumi.String("10.0.8.0/22"),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
		})
		if err != nil {
			return err
		}

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
			NetworkProfile: containerservice.ContainerServiceNetworkProfileArgs{
				NetworkPlugin: pulumi.String("azure"),
				NetworkPolicy: pulumi.String("azure"),
				ServiceCidr:   pulumi.String("10.96.0.0/16"),
				DnsServiceIP:  pulumi.String("10.96.0.10"),
			},
			ResourceGroupName: resourceGroup.Name,
		})
		if err != nil {
			return err
		}

		// Build a user Kubeconfig
		// This Kubeconfig SHOULD NOT be used for an explicit provider
		// This Kubeconfig SHOULD be used for user logins to the cluster
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

		// Build an admin Kubeconfig
		// This Kubeconfig SHOULD be used for an explicit provider
		// This Kubeconfig SHOULD NOT be used for user logins to the cluster
		adminCreds := containerservice.ListManagedClusterAdminCredentialsOutput(ctx,
			containerservice.ListManagedClusterAdminCredentialsOutputArgs{
				ResourceGroupName: resourceGroup.Name,
				ResourceName:      cluster.Name,
			})

		adminKubeconfig := adminCreds.Kubeconfigs().Index(pulumi.Int(0)).Value().
			ApplyT(func(encoded string) string {
				adminKubeconfig, err := base64.StdEncoding.DecodeString(encoded)
				if err != nil {
					return ""
				}
				return string(adminKubeconfig)
			})

		// Export some values for use elsewhere
		ctx.Export("resourceGroupName", resourceGroup.Name)
		ctx.Export("networkName", virtualNetwork.Name)
		ctx.Export("subnet1Name", subnet1.Name)
		ctx.Export("subnet2Name", subnet2.Name)
		ctx.Export("subnet3Name", subnet3.Name)
		ctx.Export("clusterName", cluster.Name)
		ctx.Export("kubeconfig", kubeconfig)
		ctx.Export("adminKubeconfig", pulumi.ToSecret(adminKubeconfig))

		return nil
	})
}
