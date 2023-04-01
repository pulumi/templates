package main

import (
	"github.com/lbrlabs/pulumi-scaleway/sdk/go/scaleway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cluster, err := scaleway.NewKubernetesCluster(ctx, "pulumi-kapsule", &scaleway.KubernetesClusterArgs{
			Name:    pulumi.String("pulumi-kapsule"),
			Version: pulumi.String("1.23"),
			Region:  pulumi.String("fr-par"),
			Cni:     pulumi.String("cilium"),
			Tags: pulumi.StringArray{
				pulumi.String("pulumi"),
				pulumi.String("kapsule"),
			},
			AutoUpgrade: &scaleway.KubernetesClusterAutoUpgradeArgs{
				Enable:                     pulumi.Bool(true),
				MaintenanceWindowStartHour: pulumi.Int(3),
				MaintenanceWindowDay:       pulumi.String("sunday"),
			},
		})
		if err != nil {
			return err
		}
		_, err = scaleway.NewKubernetesNodePool(ctx, "pulumi-kapsule-pool", &scaleway.KubernetesNodePoolArgs{
			Zone:        pulumi.String("fr-par-1"),
			Name:        pulumi.String("pulumi-kapsule-pool"),
			NodeType:    pulumi.String("DEV1-L"),
			Size:        pulumi.Int(1),
			Autoscaling: pulumi.Bool(true),
			MinSize:     pulumi.Int(1),
			MaxSize:     pulumi.Int(3),
			Autohealing: pulumi.Bool(true),
			ClusterId:   cluster.ID(),
		})
		if err != nil {
			return err
		}
		ctx.Export("cluster_id", cluster.ID())
		ctx.Export("kubeconfig", pulumi.ToSecret(cluster.Kubeconfigs.Index(pulumi.Int(0)).ConfigFile()))
		return nil
	})
}
