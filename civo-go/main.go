package main

import (
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		firewall, err := civo.NewFirewall(ctx, "civo-firewall", &civo.FirewallArgs{
			Name:               pulumi.String("myFirstFirewall"),
			Region:             pulumi.StringPtr("LON1"),
			CreateDefaultRules: pulumi.BoolPtr(true),
		})
		if err != nil {
			return err
		}
		cluster, err := civo.NewKubernetesCluster(ctx, "civo-k3s-cluster", &civo.KubernetesClusterArgs{
			Name:            pulumi.StringPtr("myFirstCivoCluster"),
			NumTargetNodes:  pulumi.IntPtr(3),
			TargetNodesSize: pulumi.StringPtr("g3.k3s.medium"),
			Region:          pulumi.StringPtr("LON1"),
			FirewallId:      firewall.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("name", cluster.Name)
		return nil
	})
}
