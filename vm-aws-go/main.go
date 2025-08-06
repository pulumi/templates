package main

import (
	"fmt"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or set default values.
		cfg := config.New(ctx, "")
		instanceType := "t3.micro"
		if param := cfg.Get("instanceType"); param != "" {
			instanceType = param
		}
		vpcNetworkCidr := "10.0.0.0/16"
		if param := cfg.Get("vpcNetworkCidr"); param != "" {
			vpcNetworkCidr = param
		}

		// Look up the latest Amazon Linux 2 AMI.
		ami, err := ec2.LookupAmi(ctx, &ec2.LookupAmiArgs{
			Filters: []ec2.GetAmiFilter{
				ec2.GetAmiFilter{
					Name: "name",
					Values: []string{
						"amzn2-ami-hvm-*",
					},
				},
			},
			Owners: []string{
				"amazon",
			},
			MostRecent: pulumi.BoolRef(true),
		}, nil)
		if err != nil {
			return err
		}

		// User data to start a HTTP server in the EC2 instance
		userData := "#!/bin/bash\necho \"Hello, World from Pulumi!\" > index.html\nnohup python -m SimpleHTTPServer 80 &\n"

		// Create VPC.
		vpc, err := ec2.NewVpc(ctx, "vpc", &ec2.VpcArgs{
			CidrBlock:          pulumi.String(vpcNetworkCidr),
			EnableDnsHostnames: pulumi.Bool(true),
			EnableDnsSupport:   pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create an internet gateway.
		gateway, err := ec2.NewInternetGateway(ctx, "gateway", &ec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
		})
		if err != nil {
			return err
		}

		// Create a subnet that automatically assigns new instances a public IP address.
		subnet, err := ec2.NewSubnet(ctx, "subnet", &ec2.SubnetArgs{
			VpcId:               vpc.ID(),
			CidrBlock:           pulumi.String("10.0.1.0/24"),
			MapPublicIpOnLaunch: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Create a route table.
		routeTable, err := ec2.NewRouteTable(ctx, "routeTable", &ec2.RouteTableArgs{
			VpcId: vpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock: pulumi.String("0.0.0.0/0"),
					GatewayId: gateway.ID(),
				},
			},
		})
		if err != nil {
			return err
		}

		// Associate the route table with the public subnet.
		_, err = ec2.NewRouteTableAssociation(ctx, "routeTableAssociation", &ec2.RouteTableAssociationArgs{
			SubnetId:     subnet.ID(),
			RouteTableId: routeTable.ID(),
		})
		if err != nil {
			return err
		}

		// Create a security group allowing inbound access over port 80 and outbound
		// access to anywhere.
		secGroup, err := ec2.NewSecurityGroup(ctx, "secGroup", &ec2.SecurityGroupArgs{
			Description: pulumi.String("Enable HTTP access"),
			VpcId:       vpc.ID(),
			Ingress: ec2.SecurityGroupIngressArray{
				&ec2.SecurityGroupIngressArgs{
					FromPort: pulumi.Int(80),
					ToPort:   pulumi.Int(80),
					Protocol: pulumi.String("tcp"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				&ec2.SecurityGroupEgressArgs{
					FromPort: pulumi.Int(0),
					ToPort:   pulumi.Int(0),
					Protocol: pulumi.String("-1"),
					CidrBlocks: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Create and launch an EC2 instance into the public subnet.
		server, err := ec2.NewInstance(ctx, "server", &ec2.InstanceArgs{
			InstanceType: pulumi.String(instanceType),
			SubnetId:     subnet.ID(),
			VpcSecurityGroupIds: pulumi.StringArray{
				secGroup.ID(),
			},
			UserData: pulumi.String(userData),
			Ami:      pulumi.String(ami.Id),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("webserver"),
			},
		})
		if err != nil {
			return err
		}

		// Export the instance's publicly accessible IP address and hostname.
		ctx.Export("ip", server.PublicIp)
		ctx.Export("hostname", server.PublicDns)
		ctx.Export("url", server.PublicDns.ApplyT(func(publicDns string) (string, error) {
			return fmt.Sprintf("http://%v", publicDns), nil
		}).(pulumi.StringOutput))
		return nil
	})
}
