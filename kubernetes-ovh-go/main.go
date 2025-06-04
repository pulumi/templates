package main

import (
	"github.com/ovh/pulumi-ovh/sdk/go/ovh/cloudproject"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values or use defaults
		cfg := config.New(ctx, "")
		ovhServiceName := cfg.Require("ovhServiceName")
		ovhRegion, err := cfg.Try("ovhRegion")
		if err != nil {
			ovhRegion = "GRA9"
		}

		clusterName, err := cfg.Try("clusterName")
		if err != nil {
			clusterName = "my_desired_cluster"
		}

		nodePoolName, err := cfg.Try("nodePoolName")
		if err != nil {
			nodePoolName = "my-desired-pool"
		}

		nodePoolDesiredNodes, err := cfg.TryInt("nodePoolDesiredNodes")
		if err != nil {
			nodePoolDesiredNodes = 1
		}

		nodePoolMaxNodes, err := cfg.TryInt("nodePoolMaxNodes")
		if err != nil {
			nodePoolMaxNodes = 3
		}

		nodePoolMinNodes, err := cfg.TryInt("nodePoolMinNodes")
		if err != nil {
			nodePoolMinNodes = 1
		}

		flavorName, err := cfg.Try("flavorName")
		if err != nil {
			flavorName = "b3-8"
		}

		// Deploy a new Kubernetes cluster
		myKube, err := cloudproject.NewKube(ctx, clusterName, &cloudproject.KubeArgs{
			ServiceName: pulumi.String(ovhServiceName),
			Name:        pulumi.String(clusterName),
			Region:      pulumi.String(ovhRegion),
		})
		if err != nil {
			return err
		}

		// Export kubeconfig file to a secret
		ctx.Export("kubeconfig", pulumi.ToSecret(myKube.Kubeconfig))

		//Create a Node Pool
		nodePool, err := cloudproject.NewKubeNodePool(ctx, nodePoolName, &cloudproject.KubeNodePoolArgs{
			ServiceName:  pulumi.String(ovhServiceName),
			KubeId:       myKube.ID(),
			Name:         pulumi.String(nodePoolName),
			DesiredNodes: pulumi.Int(nodePoolDesiredNodes),
			MaxNodes:     pulumi.Int(nodePoolMaxNodes),
			MinNodes:     pulumi.Int(nodePoolMinNodes),
			FlavorName:   pulumi.String(flavorName),
		})
		if err != nil {
			return err
		}

		ctx.Export("nodePoolID", nodePool.ID())

		return nil
	})
}

