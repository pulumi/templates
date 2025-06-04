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
			ovhRegion = "GRA"
		}

		planName, err := cfg.Try("planName")
		if err != nil {
			planName = "SMALL"
		}

		registryName, err := cfg.Try("registryName")
		if err != nil {
			registryName = "my-registry"
		}

		registryUserName, err := cfg.Try("registryUserName")
		if err != nil {
			registryUserName = "user"
		}

		registryUserEmail, err := cfg.Try("registryUserEmail")
		if err != nil {
			registryUserEmail = "myuser@ovh.com"
		}

		registryUserLogin, err := cfg.Try("registryUserLogin")
		if err != nil {
			registryUserLogin = "myuser"
		}

		// Initiate the configuration of the registry
		regcap, err := cloudproject.GetCapabilitiesContainerFilter(ctx, &cloudproject.GetCapabilitiesContainerFilterArgs{
			ServiceName: ovhServiceName,
			PlanName:    planName,
			Region:      ovhRegion,
		}, nil)
		if err != nil {
			return err
		}

		// Deploy a new Managed private registry
		myRegistry, err := cloudproject.NewContainerRegistry(ctx, registryName, &cloudproject.ContainerRegistryArgs{
			ServiceName: pulumi.String(regcap.ServiceName),
			PlanId:      pulumi.String(regcap.Id),
			Region:      pulumi.String(regcap.Region),
		})
		if err != nil {
			return err
		}

		// Create a Private Registry User
		myRegistryUser, err := cloudproject.NewContainerRegistryUser(ctx, registryUserName, &cloudproject.ContainerRegistryUserArgs{
			ServiceName: pulumi.String(regcap.ServiceName),
			RegistryId:  myRegistry.ID(),
			Email:       pulumi.String(registryUserEmail),
			Login:       pulumi.String(registryUserLogin),
		})
		if err != nil {
			return err
		}

		// Add as an output registry information
		ctx.Export("registryURL", myRegistry.Url)
		ctx.Export("registryUser", myRegistryUser.User)
		ctx.Export("registryPassword", myRegistryUser.Password)

		return nil
	})
}
