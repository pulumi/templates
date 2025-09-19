package main

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		machineType, err := cfg.Try("machineType")
		if err != nil {
			machineType = "f1-micro"
		}

		osImage, err := cfg.Try("osImage")
		if err != nil {
			osImage = "debian-11"
		}

		instanceTag, err := cfg.Try("instanceTag")
		if err != nil {
			instanceTag = "webserver"
		}

		servicePort, err := cfg.Try("servicePort")
		if err != nil {
			servicePort = "80"
		}

		// Create a new network for the virtual machine.
		network, err := compute.NewNetwork(ctx, "network", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Create a subnet on the network.
		subnet, err := compute.NewSubnetwork(ctx, "subnet", &compute.SubnetworkArgs{
			IpCidrRange: pulumi.String("10.0.1.0/24"),
			Network:     network.ID(),
		})
		if err != nil {
			return err
		}

		// Create a firewall allowing inbound access over ports 80 (for HTTP) and 22 (for SSH).
		firewall, err := compute.NewFirewall(ctx, "firewall", &compute.FirewallArgs{
			Network: network.SelfLink,
			Allows: compute.FirewallAllowArray{
				compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.ToStringArray([]string{
						"22",
						servicePort,
					}),
				},
			},
			Direction: pulumi.String("INGRESS"),
			SourceRanges: pulumi.ToStringArray([]string{
				"0.0.0.0/0",
			}),
			TargetTags: pulumi.ToStringArray([]string{
				instanceTag,
			}),
		})
		if err != nil {
			return err
		}

		// Define a script to be run when the VM starts up.
		metadataStartupScript := fmt.Sprintf(`#!/bin/bash
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

		// Create the virtual machine.
		instance, err := compute.NewInstance(ctx, "instance", &compute.InstanceArgs{
			MachineType: pulumi.String(machineType),
			BootDisk: compute.InstanceBootDiskArgs{
				InitializeParams: compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String(osImage),
				},
			},
			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				compute.InstanceNetworkInterfaceArgs{
					Network:    network.ID(),
					Subnetwork: subnet.ID(),
					AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
						compute.InstanceNetworkInterfaceAccessConfigArgs{
							// NatIp:       nil,
							// NetworkTier: nil,
						},
					},
				},
			},
			ServiceAccount: compute.InstanceServiceAccountArgs{
				Scopes: pulumi.ToStringArray([]string{
					"https://www.googleapis.com/auth/cloud-platform",
				}),
			},
			AllowStoppingForUpdate: pulumi.Bool(true),
			MetadataStartupScript:  pulumi.String(metadataStartupScript),
			Tags: pulumi.ToStringArray([]string{
				instanceTag,
			}),
		}, pulumi.DependsOn([]pulumi.Resource{firewall}))
		if err != nil {
			return err
		}

		instanceIp := instance.NetworkInterfaces.Index(pulumi.Int(0)).AccessConfigs().Index(pulumi.Int(0)).NatIp()

		// Export the instance's name, public IP address, and HTTP URL.
		ctx.Export("name", instance.Name)
		ctx.Export("ip", instanceIp)
		ctx.Export("url", pulumi.Sprintf("http://%s:%s", instanceIp.Elem(), servicePort))
		return nil
	})
}
