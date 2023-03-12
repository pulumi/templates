using Pulumi;
using Gcp = Pulumi.Gcp;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Pulumi.Config();
    var machineType = config.Get("machineType") ?? "f1-micro";
    var osImage = config.Get("osImage") ?? "debian-11";
    var instanceTag = config.Get("instanceTag") ?? "webserver";
    var servicePort = config.Get("servicePort") ?? "80";

    // Create a new network for the virtual machine.
    var network = new Gcp.Compute.Network("network", new()
    {
        AutoCreateSubnetworks = false,
    });

    // Create a subnet on the network.
    var subnet = new Gcp.Compute.Subnetwork("subnet", new()
    {
        IpCidrRange = "10.0.1.0/24",
        Network = network.Id,
    });

    // Create a firewall allowing inbound access over ports 80 (for HTTP) and 22 (for SSH).
    var firewall = new Gcp.Compute.Firewall("firewall", new()
    {
        Network = network.SelfLink,
        Allows = new[]
        {
            new Gcp.Compute.Inputs.FirewallAllowArgs {
                Protocol = "tcp",
                Ports = new[] {
                    "22",
                    servicePort,
                },
            },
        },
        Direction = "INGRESS",
        SourceRanges = new[]
        {
            "0.0.0.0/0",
        },
        TargetTags = new[]
        {
            instanceTag,
        },
    });

    // Define a script to be run when the VM starts up.
    var metadataStartupScript = $@"#!/bin/bash
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

    // Create the virtual machine.
    var instance = new Gcp.Compute.Instance("instance", new()
    {
        MachineType = machineType,
        BootDisk = new Gcp.Compute.Inputs.InstanceBootDiskArgs
        {
            InitializeParams = new Gcp.Compute.Inputs.InstanceBootDiskInitializeParamsArgs
            {
                Image = osImage,
            }
        },
        NetworkInterfaces = new[]
        {
            new Gcp.Compute.Inputs.InstanceNetworkInterfaceArgs
            {
                Network = network.Id,
                Subnetwork = subnet.Id,
                AccessConfigs = new[]
                {
                    new Gcp.Compute.Inputs.InstanceNetworkInterfaceAccessConfigArgs
                    {

                    },
                },
            },
        },
        ServiceAccount = new Gcp.Compute.Inputs.InstanceServiceAccountArgs
        {
            Scopes = new[]
            {
                "https://www.googleapis.com/auth/cloud-platform",
            },
        },
        AllowStoppingForUpdate = true,
        MetadataStartupScript = metadataStartupScript,
        Tags = new[]
        {
            instanceTag,
        },
    }, new() { DependsOn = firewall });

    var instanceIP = instance.NetworkInterfaces.Apply(interfaces => {
        return interfaces[0].AccessConfigs[0].NatIp;
    });

    // Export the instance's name, public IP address, and HTTP URL.
    return new Dictionary<string, object?>
    {
        ["name"] = instance.Name,
        ["ip"] = instanceIP,
        ["url"] = Output.Format($"http://{instanceIP}:{servicePort}"),
    };
});
