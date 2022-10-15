using System;
using System.Collections.Generic;
using Pulumi;
using AzureNative = Pulumi.AzureNative;
using Docker = Pulumi.Docker;
using Random = Pulumi.Random;

return await Pulumi.Deployment.RunAsync(() =>
{
    var config = new Config();
    var appPath = config.Get("appPath") ?? "./app";
    var imageName = config.Get("imageName") ?? "my-app";
    var containerPort = config.GetInt32("containerPort") ?? 80;
    var cpu = Math.Max(config.GetObject<double>("cpu"), 1.0);
    var memory = Math.Max(config.GetObject<double>("memory"), 1.5);

    var resourceGroup = new AzureNative.Resources.ResourceGroup("resource-group");

    var registry = new AzureNative.ContainerRegistry.Registry("registry", new()
    {
        ResourceGroupName = resourceGroup.Name,
        AdminUserEnabled = true,
        Sku = new AzureNative.ContainerRegistry.Inputs.SkuArgs {
            Name = AzureNative.ContainerRegistry.SkuName.Basic,
        },
    });

    var credentials = AzureNative.ContainerRegistry.ListRegistryCredentials.Invoke(new()
    {
        ResourceGroupName = resourceGroup.Name,
        RegistryName = registry.Name,
    });

    var registryUsername = credentials.Apply(result => result.Username!);
    var registryPassword = credentials.Apply(result => result.Passwords[0]!.Value!);

    var image = new Docker.Image("image", new()
    {
        ImageName = Pulumi.Output.Format($"{registry.LoginServer}/{imageName}"),
        Build = new Docker.DockerBuild {
            Context = appPath,
        },
        Registry = new Docker.ImageRegistry {
            Server = registry.LoginServer,
            Username = registryUsername,
            Password = registryPassword,
        },
    });

    var dnsName = new Random.RandomString("dns-name", new()
    {
        Length = 8,
        Special = false,
    }).Result.Apply(result => $"{imageName}-{result.ToLower()}");

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

    return new Dictionary<string, object?>
    {
        ["ipAddress"] = containerGroup.IpAddress.Apply(addr => addr!.Ip),
        ["hostname"] = containerGroup.IpAddress.Apply(addr => addr!.Fqdn),
        ["url"] = containerGroup.IpAddress.Apply(addr => $"http://{addr!.Fqdn}:{containerPort}"),
    };
});
