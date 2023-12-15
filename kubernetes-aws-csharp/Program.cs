using Pulumi;
using Awsx = Pulumi.Awsx;
using Eks = Pulumi.Eks;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Grab some values from the Pulumi configuration (or use default values)
    var config = new Config();
    var minClusterSize = config.GetInt32("minClusterSize") ?? 3;
    var maxClusterSize = config.GetInt32("maxClusterSize") ?? 6;
    var desiredClusterSize = config.GetInt32("desiredClusterSize") ?? 3;
    var eksNodeInstanceType = config.Get("eksNodeInstanceType") ?? "t3.medium";
    var vpcNetworkCidr = config.Get("vpcNetworkCidr") ?? "10.0.0.0/16";

    // Create a new VPC
    var eksVpc = new Awsx.Ec2.Vpc("eks-vpc", new()
    {
        EnableDnsHostnames = true,
        CidrBlock = vpcNetworkCidr,
    });

    // Create the EKS cluster
    var eksCluster = new Eks.Cluster("eks-cluster", new()
    {
        // Put the cluster in the new VPC created earlier
        VpcId = eksVpc.VpcId,

        // Public subnets will be used for load balancers
        PublicSubnetIds = eksVpc.PublicSubnetIds,

        // Private subnets will be used for cluster nodes
        PrivateSubnetIds = eksVpc.PrivateSubnetIds,

        // Change configuration values to change any of the following settings
        InstanceType = eksNodeInstanceType,
        DesiredCapacity = desiredClusterSize,
        MinSize = minClusterSize,
        MaxSize = maxClusterSize,

        // Do not give the worker nodes public IP addresses
        NodeAssociatePublicIpAddress = false,

        // Change these values for a private cluster (VPN access required)
        EndpointPrivateAccess = false,
        EndpointPublicAccess = true,
    });

    // Export some values for use elsewhere
    return new Dictionary<string, object?>
    {
        ["kubeconfig"] = eksCluster.Kubeconfig,
        ["vpcId"] = eksVpc.VpcId,
    };
});
