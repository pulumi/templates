using Pulumi;
using AzureNative = Pulumi.AzureNative;
using Docker = Pulumi.Docker;
using Random = Pulumi.Random;
using System.Collections.Generic;

return await Pulumi.Deployment.RunAsync(() =>
{
    var config = new Config();
    var appPath = config.Get("appPath") ?? "./app";
    var containerPort = config.GetInt32("containerPort") ?? 80;
    var imageName = config.Get("imageName") ?? "my-app";

    var resourceGroup = new AzureNative.Resources.ResourceGroup("resourceGroup");

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

    var hostname = new Random.RandomPet("hostname", new()
    {
        Length = 2,
    });

    var group = new AzureNative.ContainerInstance.ContainerGroup("group", new()
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
                        Cpu = 1.0,
                        MemoryInGB = 1.5,
                    },
                },
            },
        },
        IpAddress = new AzureNative.ContainerInstance.Inputs.IpAddressArgs {
            Type = AzureNative.ContainerInstance.ContainerGroupIpAddressType.Public,
            DnsNameLabel = hostname.Id,
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
        ["ipAddress"] = group.IpAddress.Apply(address => address!.Ip),
        ["hostname"] = group.IpAddress.Apply(address => address!.Fqdn),
        ["url"] = group.IpAddress.Apply(address => $"http://{address!.Fqdn}"),
    };
});
