package main

import (
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		machineType := "f1-micro"
		if param := cfg.Get("machineType"); param != "" {
			machineType = param
		}
		osImage := "debian-11"
		if param := cfg.Get("osImage"); param != "" {
			osImage = param
		}
		instanceTag := "webserver"
		if param := cfg.Get("instanceTag"); param != "" {
			instanceTag = param
		}
		servicePort := 80
		if param := cfg.GetInt("servicePort"); param != 0 {
			servicePort = param
		}

		network, err := compute.NewNetwork(ctx, "network", &compute.NetworkArgs{
			AutoCreateSubnetworks: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		subnet, err := compute.NewSubnetwork(ctx, "subnet", &compute.SubnetworkArgs{
			IpCidrRange: pulumi.String("10.0.1.0/24"),
			Network:     network.ID(),
		})
		if err != nil {
			return err
		}

		firewall, err := compute.NewFirewall(ctx, "firewall", &compute.FirewallArgs{
			Network: network.SelfLink,
			Allows: compute.FirewallAllowArray{
				compute.FirewallAllowArgs{
					Protocol: pulumi.String("tcp"),
					Ports: pulumi.ToStringArray([]string{
						"22",
						strconv.Itoa(servicePort),
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
			sudo python3 -m http.server %d &`, servicePort)

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

		ctx.Export("name", instance.Name)
		ctx.Export("ip", instanceIp)
		ctx.Export("url", pulumi.Sprintf("http://%s", instanceIp.Elem()))
		return nil
	})
}
