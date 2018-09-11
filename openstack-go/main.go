package main

import (
	"github.com/pulumi/pulumi-openstack/sdk/go/openstack/compute"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an OpenStack resource (Compute Instance)
		instance, err := compute.NewInstance(ctx, "test", &compute.InstanceArgs{
			FlavorName: "s1-2",
			ImageName:  "Ubuntu 16.04",
		})
		if err != nil {
			return err
		}

		// Export the IP of the instance
		ctx.Export("instanceIP", instance.AccessIpV4())
		return nil
	})
}
