package main

import (
	"fmt"
	"strings"

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
		cpu := 1.0
		if param := cfg.GetFloat64("cpu"); param != 0 {
			cpu = param
		}
		memory := 1.0
		if param := cfg.GetFloat64("memory"); param != 0 {
			memory = param
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

		dnsNameSuffix, err := random.NewRandomString(ctx, "dns-name-suffix", &random.RandomStringArgs{
			Length:  pulumi.Int(8),
			Special: pulumi.Bool(false),
		})
		if err != nil {
			return err
		}

		dnsName := dnsNameSuffix.Result.ApplyT(func(result string) string {
			return fmt.Sprintf("%s-%s", imageName, strings.ToLower(result))
		}).(pulumi.StringOutput)

		containerGroup, err := containerinstance.NewContainerGroup(ctx, "container-group", &containerinstance.ContainerGroupArgs{
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
							Cpu:        pulumi.Float64(cpu),
							MemoryInGB: pulumi.Float64(memory),
						},
					},
				},
			},
			IpAddress: containerinstance.IpAddressArgs{
				Type:         pulumi.String("public"),
				DnsNameLabel: dnsName,
				Ports: containerinstance.PortArray{
					containerinstance.PortArgs{
						Port:     pulumi.Int(containerPort),
						Protocol: pulumi.String("tcp"),
					},
				},
			},
		})

		ctx.Export("ipAddress", containerGroup.IpAddress.Elem().Ip())
		ctx.Export("hostname", containerGroup.IpAddress.Elem().Fqdn())
		ctx.Export("url", pulumi.Sprintf("http://%s:%d", containerGroup.IpAddress.Elem().Fqdn(), containerPort))

		return nil
	})
}
