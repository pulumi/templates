package main

import (
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or set default values
		cfg := config.New(ctx, "")
		vpcNetworkCidr, err := cfg.Try("vpcNetworkCidr")
		if err != nil {
			vpcNetworkCidr = "10.0.0.0/16"
		}

		// Create a new VPC, subnets, and associated components
		vpc, err := ec2.NewVpc(ctx, "vpc", &ec2.VpcArgs{
			EnableDnsHostnames: pulumi.Bool(true),
			CidrBlock:          &vpcNetworkCidr,
		})
		if err != nil {
			return err
		}

		// Create a security group to allow HTTPS traffic
		securityGroup, err := awsec2.NewSecurityGroup(ctx, "securityGroup", &awsec2.SecurityGroupArgs{
			VpcId:       vpc.VpcId,
			Description: pulumi.String("Allow HTTPS traffic"),
			Ingress: awsec2.SecurityGroupIngressArray{
				awsec2.SecurityGroupIngressArgs{
					Protocol:    pulumi.String("tcp"),
					ToPort:      pulumi.Int(443),
					FromPort:    pulumi.Int(443),
					Description: pulumi.String("Allow inbound HTTPS"),
					CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: awsec2.SecurityGroupEgressArray{
				awsec2.SecurityGroupEgressArgs{
					Protocol:    pulumi.String("-1"),
					ToPort:      pulumi.Int(0),
					FromPort:    pulumi.Int(0),
					Description: pulumi.String("Allow all outbound traffic"),
					CidrBlocks:  pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Export some values for use elsewhere
		ctx.Export("vpcId", vpc.VpcId)
		ctx.Export("securityGroupId", securityGroup.ID())
		return nil
	})
}
