package main

import (
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi-eks/sdk/go/eks"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or set default values
		cfg := config.New(ctx, "")
		minClusterSize, err := cfg.TryInt("minClusterSize")
		if err != nil {
			minClusterSize = 3
		}
		maxClusterSize, err := cfg.TryInt("maxClusterSize")
		if err != nil {
			maxClusterSize = 6
		}
		desiredClusterSize, err := cfg.TryInt("desiredClusterSize")
		if err != nil {
			desiredClusterSize = 3
		}
		eksNodeInstanceType, err := cfg.Try("eksNodeInstanceType")
		if err != nil {
			eksNodeInstanceType = "t2.medium"
		}
		vpcNetworkCidr, err := cfg.Try("vpcNetworkCidr")
		if err != nil {
			vpcNetworkCidr = "10.0.0.0/16"
		}

		// Create a new VPC, subnets, and associated infrastructure
		eksVpc, err := ec2.NewVpc(ctx, "eks-vpc", &ec2.VpcArgs{
			EnableDnsHostnames: pulumi.Bool(true),
			CidrBlock:          &vpcNetworkCidr,
		})
		if err != nil {
			return err
		}

		// Create a new EKS cluster
		eksCluster, err := eks.NewCluster(ctx, "eks-cluster", &eks.ClusterArgs{
			VpcId:                        eksVpc.VpcId,
			PublicSubnetIds:              eksVpc.PublicSubnetIds,
			PrivateSubnetIds:             eksVpc.PrivateSubnetIds,
			InstanceType:                 pulumi.String(eksNodeInstanceType),
			DesiredCapacity:              pulumi.Int(desiredClusterSize),
			MinSize:                      pulumi.Int(minClusterSize),
			MaxSize:                      pulumi.Int(maxClusterSize),
			NodeAssociatePublicIpAddress: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		// Export some values in case they are needed
		ctx.Export("kubeconfig", eksCluster.Kubeconfig)
		ctx.Export("vpcId", eksVpc.VpcId)
		return nil
	})
}
