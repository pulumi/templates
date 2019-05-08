package main

import (
	"github.com/pulumi/pulumi-linode/sdk/go/linode"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a linode resource (Linode Instance)
		instance, err := linode.NewInstance(ctx, "my-linode", &linode.InstanceArgs{
			Type: "g6-nanode-1", 
			Region: "us-east", 
			Image: "linode/ubuntu18.04",
		})
		if err != nil {
			return err
		}

		// Export the DNS name of the instance
		ctx.Export("instanceIpAddress", instance.IpAddress())
		return nil
	})
}
