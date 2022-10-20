using System.Collections.Generic;
using Pulumi;
using Gcp = Pulumi.Gcp;

return await Deployment.RunAsync(() => 
{
    // Get some provider-namespaced configuration values (or use defaults)
    var provCfg = new Config("gcp");
    var gcpProject = provCfg.Require("project");
    var gcpRegion = provCfg.Get("region") ?? "us-central1";
    // Get some additional configuration values (or use defaults)
    var config = new Config();
    var nodesPerZone = config.GetInt32("nodesPerZone") ?? 1;

    // Create a new network
    var gkeNetwork = new Gcp.Compute.Network("gke-network", new()
    {
        AutoCreateSubnetworks = false,
        Description = "A virtual network for your GKE cluster(s)",
    });

    // Create a new subnet within the network
    var gkeSubnet = new Gcp.Compute.Subnetwork("gke-subnet", new()
    {
        IpCidrRange = "10.128.0.0/12",
        Network = gkeNetwork.Id,
        PrivateIpGoogleAccess = true,
    });

    // Create a new GKE cluster
    var gkeCluster = new Gcp.Container.Cluster("gke-cluster", new()
    {
        AddonsConfig = new Gcp.Container.Inputs.ClusterAddonsConfigArgs
        {
            DnsCacheConfig = new Gcp.Container.Inputs.ClusterAddonsConfigDnsCacheConfigArgs
            {
                Enabled = true,
            },
        },
        BinaryAuthorization = new Gcp.Container.Inputs.ClusterBinaryAuthorizationArgs
        {
            EvaluationMode = "PROJECT_SINGLETON_POLICY_ENFORCE",
        },
        DatapathProvider = "ADVANCED_DATAPATH",
        Description = "A GKE cluster",
        InitialNodeCount = 1,
        IpAllocationPolicy = new Gcp.Container.Inputs.ClusterIpAllocationPolicyArgs
        {
            ClusterIpv4CidrBlock = "/14",
            ServicesIpv4CidrBlock = "/20",
        },
        Location = gcpRegion,
        MasterAuthorizedNetworksConfig = new Gcp.Container.Inputs.ClusterMasterAuthorizedNetworksConfigArgs
        {
            CidrBlocks = new[]
            {
                new Gcp.Container.Inputs.ClusterMasterAuthorizedNetworksConfigCidrBlockArgs
                {
                    CidrBlock = "0.0.0.0/0",
                    DisplayName = "All networks",
                },
            },
        },
        Network = gkeNetwork.Name,
        NetworkingMode = "VPC_NATIVE",
        PrivateClusterConfig = new Gcp.Container.Inputs.ClusterPrivateClusterConfigArgs
        {
            EnablePrivateNodes = true,
            EnablePrivateEndpoint = false,
            MasterIpv4CidrBlock = "10.100.0.0/28",
        },
        RemoveDefaultNodePool = true,
        ReleaseChannel = new Gcp.Container.Inputs.ClusterReleaseChannelArgs
        {
            Channel = "STABLE",
        },
        Subnetwork = gkeSubnet.Name,
        WorkloadIdentityConfig = new Gcp.Container.Inputs.ClusterWorkloadIdentityConfigArgs
        {
            WorkloadPool = $"{gcpProject}.svc.id.goog",
        },
    });

    // Create a GCP service account for the node pool
    var gkeNodepoolSa = new Gcp.ServiceAccount.Account("gke-nodepool-sa", new()
    {
        AccountId = gkeCluster.Name.Apply(name => $"{name}-np-1-sa"),
        DisplayName = "Nodepool 1 Service Account",
    });

    // Create a new node pool
    var gkeNodepool = new Gcp.Container.NodePool("gke-nodepool", new()
    {
        Cluster = gkeCluster.Id,
        NodeCount = nodesPerZone,
        NodeConfig = new Gcp.Container.Inputs.NodePoolNodeConfigArgs
        {
            OauthScopes = new[]
            {
                "https://www.googleapis.com/auth/cloud-platform",
            },
            ServiceAccount = gkeNodepoolSa.Email,
        },
    });

    // Build Kubeconfig for accessing the cluster
    var clusterKubeconfig = Output.CreateSecret(Output.Tuple(gkeCluster.MasterAuth, gkeCluster.Endpoint, gkeCluster.Name).Apply(values =>
    {
        var masterAuth = values.Item1;
        var endpoint = values.Item2;
        var gkeClusterName = values.Item3;
        return @$"apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {masterAuth.ClusterCaCertificate}
    server: https://{endpoint}
  name: {gkeClusterName}
contexts:
- context:
    cluster: {gkeClusterName}
    user: {gkeClusterName}
  name: {gkeClusterName}
current-context: {gkeClusterName}
kind: Config
preferences: {{}}
users:
- name: {gkeClusterName}
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
      installHint: Install gke-gcloud-auth-plugin for use with kubectl by following
        https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke
      provideClusterInfo: true
";
    }));

    // Export some values for use elsewhere
    return new Dictionary<string, object?>
    {
        ["networkName"] = gkeNetwork.Name,
        ["networkId"] = gkeNetwork.Id,
        ["clusterName"] = gkeCluster.Name,
        ["clusterId"] = gkeCluster.Id,
        ["kubeconfig"] = clusterKubeconfig,
    };
});
