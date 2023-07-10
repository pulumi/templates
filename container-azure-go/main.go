package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-azure-native-sdk/containerinstance/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/containerregistry/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		appPath := "./app"
		if param := cfg.Get("appPath"); param != "" {
			appPath = param
		}
		imageName := "my-app"
		if param := cfg.Get("imageName"); param != "" {
			imageName = param
		}
		imageTag := "latest"
		if param := cfg.Get("imageTag"); param != "" {
			imageName = param
		}
		containerPort := 80
		if param := cfg.GetInt("containerPort"); param != 0 {
			containerPort = param
		}
		cpu := 1
		if param := cfg.GetInt("cpu"); param != 0 {
			cpu = param
		}
		memory := 2
		if param := cfg.GetInt("memory"); param != 0 {
			memory = param
		}

		// Create a resource group for the container registry.
		resourceGroup, err := resources.NewResourceGroup(ctx, "resource-group", nil)
		if err != nil {
			return err
		}

		// Create a container registry.
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

		// Fetch login credentials for the registry.
		credentials := containerregistry.ListRegistryCredentialsOutput(ctx, containerregistry.ListRegistryCredentialsOutputArgs{
			ResourceGroupName: resourceGroup.Name,
			RegistryName:      registry.Name,
		})
		registryUsername := credentials.Username().Elem()
		registryPassword := credentials.Passwords().Index(pulumi.Int(0)).Value().Elem()

		// Create a container image for the service.
		image, err := docker.NewImage(ctx, "image", &docker.ImageArgs{
			ImageName: pulumi.Sprintf("%s/%s:%s", registry.LoginServer, imageName, imageTag),
			Build: docker.DockerBuildArgs{
				Context: pulumi.String(appPath),
				Platform: pulumi.String("linux/amd64"),
			},
			Registry: docker.RegistryArgs{
				Server:   registry.LoginServer,
				Username: registryUsername,
				Password: registryPassword,
			},
		})
		if err != nil {
			return err
		}

		// Use a random string to give the service a unique DNS name.
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

		// Create a container group for the service that makes it publicly accessible.
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
		if err != nil {
			return err
		}

		// Export the service's IP address, hostname, and fully-qualified URL.
		ctx.Export("ip", containerGroup.IpAddress.Elem().Ip())
		ctx.Export("hostname", containerGroup.IpAddress.Elem().Fqdn())
		ctx.Export("url", pulumi.Sprintf("http://%s:%d", containerGroup.IpAddress.Elem().Fqdn(), containerPort))

		return nil
	})
}
