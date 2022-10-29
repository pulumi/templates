using Pulumi;
using AzureNative = Pulumi.AzureNative;
using Random = Pulumi.Random;
using System.Collections.Generic;
using System;
using System.Text;

return await Pulumi.Deployment.RunAsync(() =>
{
    var config = new Pulumi.Config();
    var vmName = config.Get("vmName") ?? "my-server";
    var vmSize = config.Get("vmSize") ?? "Standard_A0";
    var osImage = config.Get("osImage") ?? "Debian:debian-11:11:latest";
    var adminUsername = config.Get("adminUsername") ?? "pulumiUser";
    var servicePort = config.Get("servicePort") ?? "80";
    var sshPublicKey = config.Require("sshPublicKey");

    string[] osImageArgs = osImage.Split(":");
    var osImagePublisher = osImageArgs[0];
    var osImageOffer = osImageArgs[1];
    var osImageSku = osImageArgs[2];
    var osImageVersion = osImageArgs[3];

    var resourceGroup = new AzureNative.Resources.ResourceGroup("resource-group");

    var virtualNetwork = new AzureNative.Network.VirtualNetwork("network", new()
    {
        ResourceGroupName = resourceGroup.Name,
        AddressSpace = new AzureNative.Network.Inputs.AddressSpaceArgs {
            AddressPrefixes = new[]
            {
                "10.0.0.0/16",
            },
        },
        Subnets = new[] {
            new AzureNative.Network.Inputs.SubnetArgs {
                Name = $"{vmName}-subnet",
                AddressPrefix = "10.0.1.0/24",
            },
        },
    });

    var domainNameLabel = new Random.RandomString("domain-label", new()
    {
        Length = 8,
        Upper= false,
        Special = false,
    }).Result.Apply(result => $"{vmName}-{result}");

    var publicIp = new AzureNative.Network.PublicIPAddress("public-ip", new()
    {
        ResourceGroupName = resourceGroup.Name,
        PublicIPAllocationMethod = AzureNative.Network.IPAllocationMethod.Dynamic,
        DnsSettings = new AzureNative.Network.Inputs.PublicIPAddressDnsSettingsArgs {
            DomainNameLabel = domainNameLabel,
        },
    });

    var securityGroup = new AzureNative.Network.NetworkSecurityGroup("security-group", new()
    {
        ResourceGroupName = resourceGroup.Name,
        SecurityRules = new[]
        {
            new AzureNative.Network.Inputs.SecurityRuleArgs {
                Name = $"{vmName}-securityrule",
                Priority = 1000,
                Direction = AzureNative.Network.SecurityRuleDirection.Inbound,
                Access = "Allow",
                Protocol = "Tcp",
                SourcePortRange = "*",
                SourceAddressPrefix = "*",
                DestinationAddressPrefix = "*",
                DestinationPortRanges = new[]
                {
                    servicePort,
                    "22",
                },
            },
        },
    });

    var networkInterface = new AzureNative.Network.NetworkInterface("network-interface", new()
    {
        ResourceGroupName = resourceGroup.Name,
        NetworkSecurityGroup = new AzureNative.Network.Inputs.NetworkSecurityGroupArgs {
            Id = securityGroup.Id
        },
        IpConfigurations = new[]
        {
            new AzureNative.Network.Inputs.NetworkInterfaceIPConfigurationArgs {
                Name = $"{vmName}-ipconfiguration",
                PrivateIPAllocationMethod = AzureNative.Network.IPAllocationMethod.Dynamic,
                Subnet = new AzureNative.Network.Inputs.SubnetArgs {
                    Id = virtualNetwork.Subnets.GetAt(0).Apply(subnet => subnet.Id!)
                },
                PublicIPAddress = new AzureNative.Network.Inputs.PublicIPAddressArgs {
                    Id = publicIp.Id,
                },
            }
        }
    });

    var initScript = $@"#!/bin/bash
        echo '<!DOCTYPE html>
        <html lang=""en"">
        <head>
            <meta charset=""utf-8"">
            <title>Hello, world!</title>
        </head>
        <body>
            <h1>Hello, world! 👋</h1>
            <p>Deployed with 💜 by <a href=""https://pulumi.com/"">Pulumi</a>.</p>
        </body>
        </html>' > index.html
        sudo python3 -m http.server {servicePort} &";

    var vm = new AzureNative.Compute.VirtualMachine("vm", new()
    {
        ResourceGroupName = resourceGroup.Name,
        NetworkProfile = new AzureNative.Compute.Inputs.NetworkProfileArgs {
            NetworkInterfaces = new[] {
                new AzureNative.Compute.Inputs.NetworkInterfaceReferenceArgs {
                    Id = networkInterface.Id,
                    Primary = true,
                },
            },
        },
        HardwareProfile = new AzureNative.Compute.Inputs.HardwareProfileArgs {
            VmSize = vmSize,
        },
        OsProfile = new AzureNative.Compute.Inputs.OSProfileArgs {
            ComputerName = vmName,
            AdminUsername = adminUsername,
            CustomData = Convert.ToBase64String(Encoding.UTF8.GetBytes(initScript)),
            LinuxConfiguration = new AzureNative.Compute.Inputs.LinuxConfigurationArgs {
                DisablePasswordAuthentication = true,
                Ssh = new AzureNative.Compute.Inputs.SshConfigurationArgs {
                    PublicKeys = new[]
                    {
                        new AzureNative.Compute.Inputs.SshPublicKeyArgs {
                            KeyData = sshPublicKey,
                            Path = $"/home/{adminUsername}/.ssh/authorized_keys",
                        },
                    },
                },
            },
        },
        StorageProfile = new AzureNative.Compute.Inputs.StorageProfileArgs {
            OsDisk = new AzureNative.Compute.Inputs.OSDiskArgs {
                Name = $"{vmName}-osdisk",
                CreateOption = AzureNative.Compute.DiskCreateOptionTypes.FromImage,
            },
            ImageReference = new AzureNative.Compute.Inputs.ImageReferenceArgs {
                Publisher = osImagePublisher,
                Offer = osImageOffer,
                Sku = osImageSku,
                Version = osImageVersion,
            },
        },
    });

    var vmAddress = vm.Id.Apply(_ => {
        return AzureNative.Network.GetPublicIPAddress.Invoke(new()
        {
            ResourceGroupName = resourceGroup.Name,
            PublicIpAddressName = publicIp.Name,
        });
    });

    return new Dictionary<string, object?>
    {
        ["hostname"] = vmAddress.Apply(addr => addr.DnsSettings!.Fqdn),
        ["ip"] = vmAddress.Apply(addr => addr.IpAddress),
        ["url"] = vmAddress.Apply(addr => $"http://{addr.DnsSettings!.Fqdn}"),
    };
});