using Pulumi;
using AzureNative = Pulumi.AzureNative;
using System.Collections.Generic;
using System.Text;
using System;

return await Pulumi.Deployment.RunAsync(() =>
{
    // Grab some values from the Pulumi stack configuration (or use defaults)
    var projCfg = new Config();
    var numWorkerNodes = projCfg.GetInt32("numWorkerNodes") ?? 3;
    var k8sVersion = projCfg.Get("kubernetesVersion") ?? "1.32";
    var prefixForDns = projCfg.Get("prefixForDns") ?? "pulumi";
    var nodeVmSize = projCfg.Get("nodeVmSize") ?? "Standard_DS2_v2";

    // The next two configuration values are required (no default can be provided)
    var mgmtGroupId = projCfg.Require("mgmtGroupId");
    var sshPubKey = projCfg.Require("sshPubKey");

    // Create a new Azure Resource Group
    var resourceGroup = new AzureNative.Resources.ResourceGroup("resourceGroup");

    // Create a new Azure Virtual Network
    var virtualNetwork = new AzureNative.Network.VirtualNetwork("virtualNetwork", new()
    {
        AddressSpace = new AzureNative.Network.Inputs.AddressSpaceArgs
        {
            AddressPrefixes = new[]
            {
                "10.0.0.0/16",
            },
        },
        ResourceGroupName = resourceGroup.Name,
    });

    // Create three subnets in the virtual network
    var subnet1 = new AzureNative.Network.Subnet("subnet1", new()
    {
        AddressPrefix = "10.0.0.0/22",
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = virtualNetwork.Name,
    });

    var subnet2 = new AzureNative.Network.Subnet("subnet2", new()
    {
        AddressPrefix = "10.0.4.0/22",
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = virtualNetwork.Name,
    });

    var subnet3 = new AzureNative.Network.Subnet("subnet3", new()
    {
        AddressPrefix = "10.0.8.0/22",
        ResourceGroupName = resourceGroup.Name,
        VirtualNetworkName = virtualNetwork.Name,
    });

    // Create an Azure Kubernetes Cluster
    var managedCluster = new AzureNative.ContainerService.ManagedCluster("managedCluster", new()
    {
        AadProfile =new AzureNative.ContainerService.Inputs.ManagedClusterAADProfileArgs
        {
            EnableAzureRBAC = true,
            Managed = true,
            AdminGroupObjectIDs = new[]
            {
                mgmtGroupId,
            },
        },
        AddonProfiles = {},
        // Use multiple agent/node pool profiles to distribute nodes across subnets
        AgentPoolProfiles = new AzureNative.ContainerService.Inputs.ManagedClusterAgentPoolProfileArgs
        {
            AvailabilityZones = new[]
            {
                "1", "2", "3",
            },
            Count = numWorkerNodes,
            EnableNodePublicIP = false,
            Mode = "System",
            Name = "systempool",
            OsType = "Linux",
            OsDiskSizeGB = 30,
            Type = "VirtualMachineScaleSets",
            VmSize = nodeVmSize,
            // Change next line for additional node pools to distribute across subnets
            VnetSubnetID = subnet1.Id,
        },

        // Change authorizedIPRanges to limit access to API server
        // Changing enablePrivateCluster requires alternate access to API server (VPN or similar)
        ApiServerAccessProfile = new AzureNative.ContainerService.Inputs.ManagedClusterAPIServerAccessProfileArgs
        {
            AuthorizedIPRanges = new[]
            {
                "0.0.0.0/0",
            },
            EnablePrivateCluster = false,
        },
        DnsPrefix = prefixForDns,
        EnableRBAC = true,
        Identity = new AzureNative.ContainerService.Inputs.ManagedClusterIdentityArgs
        {
            Type = AzureNative.ContainerService.ResourceIdentityType.SystemAssigned,
        },
        KubernetesVersion = k8sVersion,
        LinuxProfile = new AzureNative.ContainerService.Inputs.ContainerServiceLinuxProfileArgs
        {
            AdminUsername = "azureuser",
            Ssh = new AzureNative.ContainerService.Inputs.ContainerServiceSshConfigurationArgs
            {
                PublicKeys = new[]
                {
                    new AzureNative.ContainerService.Inputs.ContainerServiceSshPublicKeyArgs
                    {
                        KeyData = sshPubKey,
                    },
                },
            },
        },
        NetworkProfile = new AzureNative.ContainerService.Inputs.ContainerServiceNetworkProfileArgs
        {
            NetworkPlugin = "azure",
            NetworkPolicy = "azure",
            ServiceCidr = "10.96.0.0/16",
            DnsServiceIP = "10.96.0.10",
        },
        ResourceGroupName = resourceGroup.Name,
    });

    // Build a user Kubeconfig
    // This SHOULD NOT be used for an explicit provider
    // This SHOULD be used for user logins to the cluster
    var creds = AzureNative.ContainerService.ListManagedClusterUserCredentials.Invoke(new()
    {
        ResourceGroupName = resourceGroup.Name,
        ResourceName = managedCluster.Name,
    });
    var encoded = creds.Apply(result => result.Kubeconfigs[0]!.Value);
    var decoded = encoded.Apply(enc => {
        var bytes = Convert.FromBase64String(enc);
        return Encoding.UTF8.GetString(bytes);
    });

    // Build an admin Kubeconfig
    // This SHOULD be used for an explicit provider
    // This SHOULD NOT be used for user logins to the cluster
    var adminCreds = AzureNative.ContainerService.ListManagedClusterAdminCredentials.Invoke(new()
    {
        ResourceGroupName = resourceGroup.Name,
        ResourceName = managedCluster.Name,
    });
    var adminEncoded = adminCreds.Apply(result => result.Kubeconfigs[0]!.Value);
    var adminDecoded = adminEncoded.Apply(enc => {
        var bytes = Convert.FromBase64String(enc);
        return Encoding.UTF8.GetString(bytes);
    });


    // Export some values for use elsewhere
    return new Dictionary<string, object?>
    {
        ["rgName"] = resourceGroup.Name,
        ["networkName"] = virtualNetwork.Name,
        ["clusterName"] = managedCluster.Name,
        ["kubeconfig"] = decoded,
        ["adminKubeconfig"] = Pulumi.Output.CreateSecret(adminDecoded),
    };
});
