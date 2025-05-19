import pulumi
import pulumi_ovh as ovh

# Get some configuration values
config = pulumi.Config();

ovh_service_name = config.require('ovhServiceName')
ovh_region = config.get("ovhRegion", "GRA9")
cluster_name = config.get("clusterName", "my_desired_cluster")

node_pool_name = config.get("nodePoolName", "my-desired-pool")
node_pool_desired_nodes = config.get_int("nodePoolDesiredNodes", 1)
node_pool_max_nodes = config.get_int("nodePoolMaxNodes", 3)
node_pool_min_nodes = config.get_int("nodePoolMinNodes", 1)
flavor_name = config.get("flavorName", "b3-8")
		
# Deploy a new Kubernetes cluster
my_kube_cluster = ovh.cloudproject.Kube(cluster_name,
    region=ovh_region,
    service_name=ovh_service_name,
    name=cluster_name)

# Export kubeconfig file
pulumi.export("kubeconfig", my_kube_cluster.kubeconfig)

# Create a Node Pool
node_pool = ovh.cloudproject.KubeNodePool(node_pool_name,
    desired_nodes=node_pool_desired_nodes,
    flavor_name=flavor_name,
    kube_id=my_kube_cluster.id,
    max_nodes=node_pool_max_nodes,
    min_nodes=node_pool_min_nodes,
    service_name=ovh_service_name)

pulumi.export("nodePoolID", node_pool.id)