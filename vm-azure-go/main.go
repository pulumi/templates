package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-azure-native-sdk/compute/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/network/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	tls "github.com/pulumi/pulumi-tls/sdk/v4/go/tls"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings
		cfg := config.New(ctx, "")
		vmName, err := cfg.Try("vmName")
		if err != nil {
			vmName = "my-server"
		}
		vmSize, err := cfg.Try("vmSize")
		if err != nil {
			vmSize = "Standard_A1_v2"
		}
		osImage, err := cfg.Try("osImage")
		if err != nil {
			osImage = "Debian:debian-11:11:latest"
		}
		adminUsername, err := cfg.Try("adminUsername")
		if err != nil {
			adminUsername = "pulumiuser"
		}
		servicePort, err := cfg.Try("servicePort")
		if err != nil {
			servicePort = "80"
		}

		osImageArgs := strings.Split(osImage, ":")
		osImagePublisher := osImageArgs[0]
		osImageOffer := osImageArgs[1]
		osImageSku := osImageArgs[2]
		osImageVersion := osImageArgs[3]

		// Create an SSH key
		sshKey, err := tls.NewPrivateKey(ctx, "ssh-key", &tls.PrivateKeyArgs{
			Algorithm: pulumi.String("RSA"),
			RsaBits:   pulumi.Int(4096),
		})
		if err != nil {
			return err
		}

		// Create a resource group
		resourceGroup, err := resources.NewResourceGroup(ctx, "resource-group", nil)
		if err != nil {
			return err
		}

		// Create a virtual network
		virtualNetwork, err := network.NewVirtualNetwork(ctx, "network", &network.VirtualNetworkArgs{
			ResourceGroupName: resourceGroup.Name,
			AddressSpace: network.AddressSpaceArgs{
				AddressPrefixes: pulumi.ToStringArray([]string{
					"10.0.0.0/16",
				}),
			},
			Subnets: network.SubnetTypeArray{
				network.SubnetTypeArgs{
					Name:          pulumi.Sprintf("%s-subnet", vmName),
					AddressPrefix: pulumi.String("10.0.1.0/24"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Use a random string to give the VM a unique DNS name
		domainLabelSuffix, err := random.NewRandomString(ctx, "domain-label", &random.RandomStringArgs{
			Length:  pulumi.Int(8),
			Upper:   pulumi.Bool(false),
			Special: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}
		domainLabel := domainLabelSuffix.Result.ApplyT(func(result string) string {
			return fmt.Sprintf("%s-%s", vmName, result)
		}).(pulumi.StringOutput)

		// Create a public IP address for the VM
		publicIp, err := network.NewPublicIPAddress(ctx, "public-ip", &network.PublicIPAddressArgs{
			ResourceGroupName:        resourceGroup.Name,
			PublicIPAllocationMethod: pulumi.StringPtr("Dynamic"),
			DnsSettings: network.PublicIPAddressDnsSettingsArgs{
				DomainNameLabel: domainLabel,
			},
		})
		if err != nil {
			return err
		}

		// Create a security group allowing inbound access over ports 80 (for HTTP) and 22 (for SSH)
		securityGroup, err := network.NewNetworkSecurityGroup(ctx, "security-group", &network.NetworkSecurityGroupArgs{
			ResourceGroupName: resourceGroup.Name,
			SecurityRules: network.SecurityRuleTypeArray{
				network.SecurityRuleTypeArgs{
					Name:                     pulumi.StringPtr(fmt.Sprintf("%s-securityrule", vmName)),
					Priority:                 pulumi.Int(1000),
					Direction:                pulumi.String("Inbound"),
					Access:                   pulumi.String("Allow"),
					Protocol:                 pulumi.String("Tcp"),
					SourcePortRange:          pulumi.StringPtr("*"),
					SourceAddressPrefix:      pulumi.StringPtr("*"),
					DestinationAddressPrefix: pulumi.StringPtr("*"),
					DestinationPortRanges: pulumi.ToStringArray([]string{
						servicePort,
						"22",
					}),
				},
			},
		})
		if err != nil {
			return err
		}

		// Create a network interface with the virtual network, IP address, and security group
		networkInterface, err := network.NewNetworkInterface(ctx, "network-interface", &network.NetworkInterfaceArgs{
			ResourceGroupName: resourceGroup.Name,
			NetworkSecurityGroup: &network.NetworkSecurityGroupTypeArgs{
				Id: securityGroup.ID(),
			},
			IpConfigurations: network.NetworkInterfaceIPConfigurationArray{
				&network.NetworkInterfaceIPConfigurationArgs{
					Name:                      pulumi.String(fmt.Sprintf("%s-ipconfiguration", vmName)),
					PrivateIPAllocationMethod: pulumi.String("Dynamic"),
					Subnet: &network.SubnetTypeArgs{
						Id: virtualNetwork.Subnets.ApplyT(func(subnets []network.SubnetResponse) (string, error) {
							return *subnets[0].Id, nil
						}).(pulumi.StringOutput),
					},
					PublicIPAddress: &network.PublicIPAddressTypeArgs{
						Id: publicIp.ID(),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Define a script to be run when the VM starts up
		initScript := fmt.Sprintf(`#!/bin/bash
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
			sudo python3 -m http.server %s &`, servicePort)

		// Create the virtual machine
		vm, err := compute.NewVirtualMachine(ctx, "vm", &compute.VirtualMachineArgs{
			ResourceGroupName: resourceGroup.Name,
			NetworkProfile: &compute.NetworkProfileArgs{
				NetworkInterfaces: compute.NetworkInterfaceReferenceArray{
					&compute.NetworkInterfaceReferenceArgs{
						Id:      networkInterface.ID(),
						Primary: pulumi.Bool(true),
					},
				},
			},
			HardwareProfile: &compute.HardwareProfileArgs{
				VmSize: pulumi.String(vmSize),
			},
			OsProfile: &compute.OSProfileArgs{
				ComputerName:  pulumi.String(vmName),
				AdminUsername: pulumi.String(adminUsername),
				CustomData:    pulumi.String(base64.StdEncoding.EncodeToString([]byte(initScript))),
				LinuxConfiguration: &compute.LinuxConfigurationArgs{
					DisablePasswordAuthentication: pulumi.Bool(true),
					Ssh: &compute.SshConfigurationArgs{
						PublicKeys: compute.SshPublicKeyTypeArray{
							&compute.SshPublicKeyTypeArgs{
								KeyData: sshKey.PublicKeyOpenssh,
								Path:    pulumi.String(fmt.Sprintf("/home/%v/.ssh/authorized_keys", adminUsername)),
							},
						},
					},
				},
			},
			StorageProfile: &compute.StorageProfileArgs{
				OsDisk: &compute.OSDiskArgs{
					Name:         pulumi.String(fmt.Sprintf("%v-osdisk", vmName)),
					CreateOption: pulumi.String("FromImage"),
				},
				ImageReference: &compute.ImageReferenceArgs{
					Publisher: pulumi.String(osImagePublisher),
					Offer:     pulumi.String(osImageOffer),
					Sku:       pulumi.String(osImageSku),
					Version:   pulumi.String(osImageVersion),
				},
			},
		})
		if err != nil {
			return err
		}

		// Once the machine is created, fetch its IP address and DNS hostname
		address := vm.ID().ApplyT(func(_ pulumi.ID) network.LookupPublicIPAddressResultOutput {
			return network.LookupPublicIPAddressOutput(ctx, network.LookupPublicIPAddressOutputArgs{
				ResourceGroupName:   resourceGroup.Name,
				PublicIpAddressName: publicIp.Name,
			})
		})

		// Export the VM's hostname, public IP address, HTTP URL, and SSH private key
		ctx.Export("ip", address.ApplyT(func(addr network.LookupPublicIPAddressResult) (string, error) {
			return *addr.IpAddress, nil
		}).(pulumi.StringOutput))

		ctx.Export("hostname", address.ApplyT(func(addr network.LookupPublicIPAddressResult) (string, error) {
			return *addr.DnsSettings.Fqdn, nil
		}).(pulumi.StringOutput))

		ctx.Export("url", address.ApplyT(func(addr network.LookupPublicIPAddressResult) (string, error) {
			return fmt.Sprintf("http://%s:%s", *addr.DnsSettings.Fqdn, servicePort), nil
		}).(pulumi.StringOutput))

		ctx.Export("privatekey", sshKey.PrivateKeyOpenssh)

		return nil
	})
}
