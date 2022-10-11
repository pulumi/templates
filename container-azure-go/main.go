package main

import (
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerinstance"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/containerregistry"
	"github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	"github.com/pulumi/pulumi-docker/sdk/v3/go/docker"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		cfg := config.New(ctx, "")
		imageName := "my-app"
		if param := cfg.Get("imageName"); param != "" {
			imageName = param
		}
		appPath := "./app"
		if param := cfg.Get("appPath"); param != "" {
			appPath = param
		}
		containerPort := 80
		if param := cfg.GetInt("containerPort"); param != 0 {
			containerPort = param
		}

		resourceGroup, err := resources.NewResourceGroup(ctx, "resource-group", nil)
		if err != nil {
			return err
		}

		registry, err := containerregistry.NewRegistry(ctx, "registry", &containerregistry.RegistryArgs{
			ResourceGroupName: resourceGroup.Name,
			AdminUserEnabled:  pulumi.Bool(true),
			Sku: &containerregistry.SkuArgs{
				Name: pulumi.String(containerregistry.SkuNameBasic),
			},
		})
		if err != nil {
			return err
		}

		credentials := containerregistry.ListRegistryCredentialsOutput(ctx, containerregistry.ListRegistryCredentialsOutputArgs{
			ResourceGroupName: resourceGroup.Name,
			RegistryName:      registry.Name,
		})
		registryUsername := credentials.Username().Elem()
		registryPassword := credentials.Passwords().Index(pulumi.Int(0)).Value().Elem()

		image, err := docker.NewImage(ctx, "image", &docker.ImageArgs{
			ImageName: pulumi.Sprintf("%s/%s", registry.LoginServer, imageName),
			Build: docker.DockerBuildArgs{
				Context: pulumi.String(appPath),
			},
			Registry: docker.ImageRegistryArgs{
				Server:   registry.LoginServer,
				Username: registryUsername,
				Password: registryPassword,
			},
		})
		if err != nil {
			return err
		}

		hostname, err := random.NewRandomPet(ctx, "hostname", &random.RandomPetArgs{
			Length: pulumi.Int(2),
		})
		if err != nil {
			return err
		}

		group, err := containerinstance.NewContainerGroup(ctx, "group", &containerinstance.ContainerGroupArgs{
			ResourceGroupName: resourceGroup.Name,
			OsType:            pulumi.String("linux"),
			RestartPolicy:     pulumi.String("always"),
			ImageRegistryCredentials: containerinstance.ImageRegistryCredentialArray{
				containerinstance.ImageRegistryCredentialArgs{
					Server:   registry.LoginServer,
					Username: registryUsername,
					Password: registryPassword,
				},
			},
			Containers: containerinstance.ContainerArray{
				containerinstance.ContainerArgs{
					Name:  pulumi.String(imageName),
					Image: image.ImageName,
					Ports: containerinstance.ContainerPortArray{
						containerinstance.ContainerPortArgs{
							Port:     pulumi.Int(containerPort),
							Protocol: pulumi.String("tcp"),
						},
					},
					EnvironmentVariables: containerinstance.EnvironmentVariableArray{
						containerinstance.EnvironmentVariableArgs{
							Name:  pulumi.String("PORT"),
							Value: pulumi.Sprintf("%d", containerPort),
						},
					},
					Resources: containerinstance.ResourceRequirementsArgs{
						Requests: containerinstance.ResourceRequestsArgs{
							Cpu:        pulumi.Float64(1.0),
							MemoryInGB: pulumi.Float64(1.5),
						},
					},
				},
			},
			IpAddress: containerinstance.IpAddressArgs{
				Type:         pulumi.String("public"),
				DnsNameLabel: hostname.ID(),
				Ports: containerinstance.PortArray{
					containerinstance.PortArgs{
						Port:     pulumi.Int(containerPort),
						Protocol: pulumi.String("tcp"),
					},
				},
			},
		})

		ctx.Export("ipAddress", group.IpAddress.Elem().Ip())
		ctx.Export("hostname", group.IpAddress.Elem().Fqdn())
		ctx.Export("url", pulumi.Sprintf("http://%s:%d", group.IpAddress.Elem().Fqdn(), containerPort))

		return nil
	})
}
