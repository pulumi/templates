import pulumi
import pulumi_awsx as awsx
import pulumi_eks as eks

config = pulumi.Config()
min_cluster_size = config.get_float("minClusterSize")
if min_cluster_size is None:
    min_cluster_size = 3
max_cluster_size = config.get_float("maxClusterSize")
if max_cluster_size is None:
    max_cluster_size = 6
desired_cluster_size = config.get_float("desiredClusterSize")
if desired_cluster_size is None:
    desired_cluster_size = 3
eks_node_instance_type = config.get("eksNodeInstanceType")
if eks_node_instance_type is None:
    eks_node_instance_type = "t2.medium"
vpc_network_cidr = config.get("vpcNetworkCidr")
if vpc_network_cidr is None:
    vpc_network_cidr = "10.0.0.0/16"
eks_vpc = awsx.ec2.Vpc("eks-vpc",
    enable_dns_hostnames=True,
    cidr_block=vpc_network_cidr)
eks_cluster = eks.Cluster("eks-cluster",
    vpc_id=eks_vpc.vpc_id,
    public_subnet_ids=eks_vpc.public_subnet_ids,
    private_subnet_ids=eks_vpc.private_subnet_ids,
    instance_type=eks_node_instance_type,
    desired_capacity=desired_cluster_size,
    min_size=min_cluster_size,
    max_size=max_cluster_size,
    node_associate_public_ip_address=False)
pulumi.export("kubeconfig", eks_cluster.kubeconfig)
pulumi.export("vpcId", eks_vpc.vpc_id)
