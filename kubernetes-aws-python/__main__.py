import pulumi
import pulumi_awsx as awsx
import pulumi_eks as eks

# Get some values from the Pulumi configuration (or use defaults)
config = pulumi.Config()
min_cluster_size = 3 if config.get_float("minClusterSize") is None else config.get_float("minClusterSize")
max_cluster_size = 6 if config.get_float("maxClusterSize") is None else config.get_float("maxClusterSize")
desired_cluster_size = 3 if config.get_float("desiredClusterSize") is None else config.get_float("desiredClusterSize")
eks_node_instance_type = "t2.medium" if config.get("eksNodeInstanceType") is None else config.get("eksNodeInstanceType")
vpc_network_cidr = "10.0.0.0/16" if config.get("vpcNetworkCidr") is None else config.get("vpcNetworkCidr")

# Create a VPC for the EKS cluster
eks_vpc = awsx.ec2.Vpc("eks-vpc",
    enable_dns_hostnames=True,
    cidr_block=vpc_network_cidr)

# Create the EKS cluster
eks_cluster = eks.Cluster("eks-cluster",
    # Put the cluster in the new VPC created earlier
    vpc_id=eks_vpc.vpc_id,
    # Public subnets will be used for load balancers
    public_subnet_ids=eks_vpc.public_subnet_ids,
    # Private subnets will be used for cluster nodes
    private_subnet_ids=eks_vpc.private_subnet_ids,
    # Change configuration values to change any of the following settings
    instance_type=eks_node_instance_type,
    desired_capacity=desired_cluster_size,
    min_size=min_cluster_size,
    max_size=max_cluster_size,
    # Do not give worker nodes a public IP address
    node_associate_public_ip_address=False,
    # Uncomment the next two lines for private cluster (VPN access required)
    # endpoint_private_access=true,
    # endpoint_public_access=false
    )

# Export values to use elsewhere
pulumi.export("kubeconfig", eks_cluster.kubeconfig)
pulumi.export("vpcId", eks_vpc.vpc_id)
