using System;
using System.Collections.Generic;
using Pulumi;
using AzureNative = Pulumi.AzureNative;
using Docker = Pulumi.Docker;
using Random = Pulumi.Random;

return await Pulumi.Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var appPath = config.Get("appPath") ?? "./app";
    var imageName = config.Get("imageName") ?? "my-app";
    var imageTag = config.Get("imageTag") ?? "latest";
    var containerPort = config.GetInt32("containerPort") ?? 80;
    var cpu = config.GetInt32("cpu") ?? 1;
    var memory = config.GetInt32("memory") ?? 2;

    // Create a resource group for the container registry.
    var resourceGroup = new AzureNative.Resources.ResourceGroup("resource-group");

    // Create a container registry.
    var registry = new AzureNative.ContainerRegistry.Registry("registry", new()
    {
        ResourceGroupName = resourceGroup.Name,
        AdminUserEnabled = true,
        Sku = new AzureNative.ContainerRegistry.Inputs.SkuArgs {
            Name = AzureNative.ContainerRegistry.SkuName.Basic,
        },
    });

    // Fetch login credentials for the registry.
    var credentials = AzureNative.ContainerRegistry.ListRegistryCredentials.Invoke(new()
    {
        ResourceGroupName = resourceGroup.Name,
        RegistryName = registry.Name,
    });
    var registryUsername = credentials.Apply(result => result.Username!);
    var registryPassword = credentials.Apply(result => result.Passwords[0]!.Value!);

    // Create a container image for the service.
    var image = new Docker.Image("image", new()
    {
        ImageName = Pulumi.Output.Format($"{registry.LoginServer}/{imageName}:{imageTag}"),
        Build = new Docker.Inputs.DockerBuildArgs {
            Context = appPath,
            Platform = "linux/amd64",
        },
        Registry = new Docker.Inputs.RegistryArgs {
            Server = registry.LoginServer,
            Username = registryUsername,
            Password = registryPassword,
        },
    });

    // Use a random string to give the service a unique DNS name.
    var dnsName = new Random.RandomString("dns-name", new()
    {
        Length = 8,
        Special = false,
    }).Result.Apply(result => $"{imageName}-{result.ToLower()}");

    // Create a container group for the service that makes it publicly accessible.
    var containerGroup = new AzureNative.ContainerInstance.ContainerGroup("container-group", new()
    {
        ResourceGroupName = resourceGroup.Name,
        OsType = "linux",
        RestartPolicy = "always",
        ImageRegistryCredentials = new AzureNative.ContainerInstance.Inputs.ImageRegistryCredentialArgs {
            Server = registry.LoginServer,
            Username = registryUsername,
            Password = registryPassword,
        },
        Containers = new[]
        {
            new AzureNative.ContainerInstance.Inputs.ContainerArgs {
                Name = imageName,
                Image = image.ImageName,
                Ports = new[]
                {
                    new AzureNative.ContainerInstance.Inputs.ContainerPortArgs {
                        Port = containerPort,
                        Protocol = "tcp",
                    },
                },
                EnvironmentVariables = new[]
                {
                    new AzureNative.ContainerInstance.Inputs.EnvironmentVariableArgs {
                        Name = "ASPNETCORE_URLS",
                        Value = $"http://0.0.0.0:{containerPort}",
                    },
                },
                Resources = new AzureNative.ContainerInstance.Inputs.ResourceRequirementsArgs {
                    Requests = new AzureNative.ContainerInstance.Inputs.ResourceRequestsArgs {
                        Cpu = cpu,
                        MemoryInGB = memory,
                    },
                },
            },
        },
        IpAddress = new AzureNative.ContainerInstance.Inputs.IpAddressArgs {
            Type = AzureNative.ContainerInstance.ContainerGroupIpAddressType.Public,
            DnsNameLabel = dnsName,
            Ports = new[]
            {
                new AzureNative.ContainerInstance.Inputs.PortArgs {
                    Port = containerPort,
                    Protocol = "tcp",
                },
            },
        }
    });

    // Export the service's IP address, hostname, and fully-qualified URL.
    return new Dictionary<string, object?>
    {
        ["hostname"] = containerGroup.IpAddress.Apply(addr => addr!.Fqdn),
        ["ip"] = containerGroup.IpAddress.Apply(addr => addr!.Ip),
        ["url"] = containerGroup.IpAddress.Apply(addr => $"http://{addr!.Fqdn}:{containerPort}"),
    };
});
