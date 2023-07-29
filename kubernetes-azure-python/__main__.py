"""An Azure RM Python Pulumi program"""

import base64
import pulumi
from pulumi_azure_native import resources
from pulumi_azure_native import network
from pulumi_azure_native import containerservice

# Get some project-namespaced configuration values or use default values
proj_cfg = pulumi.Config()
num_worker_nodes = proj_cfg.get_float("numWorkerNodes", 3)
k8s_version = proj_cfg.get("kubernetesVersion", "1.26.3")
prefix_for_dns = proj_cfg.get("prefixForDns", "pulumi")
node_vm_size = proj_cfg.get("nodeVmSize", "Standard_DS2_v2")
# The next two configuration values are required (no default can be provided)
mgmt_group_id = proj_cfg.require("mgmtGroupId")
ssh_pub_key = proj_cfg.require("sshPubKey")

# Create an Azure Resource Group
resource_group = resources.ResourceGroup(
    "resource_group"
)

# Create an Azure Virtual Network
virtual_network = network.VirtualNetwork(
    "virtual_network",
    address_space=network.AddressSpaceArgs(
        address_prefixes=["10.0.0.0/16"],
    ),
    resource_group_name=resource_group.name
)

# Create three subnets in the virtual network
subnet1 = network.Subnet(
    "subnet-1",
    address_prefix="10.0.0.0/22",
    resource_group_name=resource_group.name,
    virtual_network_name=virtual_network.name
)
subnet2 = network.Subnet(
    "subnet-2",
    address_prefix="10.0.4.0/22",
    resource_group_name=resource_group.name,
    virtual_network_name=virtual_network.name
)
subnet3 = network.Subnet(
    "subnet-3",
    address_prefix="10.0.8.0/22",
    resource_group_name=resource_group.name,
    virtual_network_name=virtual_network.name
)

# Create an Azure Kubernetes Service cluster
managed_cluster = containerservice.ManagedCluster(
    "managed_cluster",
    aad_profile=containerservice.ManagedClusterAADProfileArgs(
        enable_azure_rbac=True,
        managed=True,
        admin_group_object_ids=[mgmt_group_id],
    ),
    # Use multiple agent/node pools to distribute nodes across subnets
    agent_pool_profiles=[containerservice.ManagedClusterAgentPoolProfileArgs(
        availability_zones=["1","2","3",],
        count=3,
        enable_node_public_ip=False,
        mode="System",
        name="systempool",
        os_type="Linux",
        os_disk_size_gb=30,
        type="VirtualMachineScaleSets",
        vm_size=node_vm_size,
        # Change next line for additional node pools to distribute across subnets
        vnet_subnet_id=subnet1.id
    )],
    # Change authorized_ip_ranges to limit access to API server
    # Changing enable_private_cluster requires alternate access to API server (VPN or similar)
    api_server_access_profile=containerservice.ManagedClusterAPIServerAccessProfileArgs(
        authorized_ip_ranges=["0.0.0.0/0"],
        enable_private_cluster=False
    ),
    dns_prefix=prefix_for_dns,
    enable_rbac=True,
    identity=containerservice.ManagedClusterIdentityArgs(
        type=containerservice.ResourceIdentityType.SYSTEM_ASSIGNED,
    ),
    kubernetes_version=k8s_version,
    linux_profile=containerservice.ContainerServiceLinuxProfileArgs(
        admin_username="azureuser",
        ssh=containerservice.ContainerServiceSshConfigurationArgs(
            public_keys=[containerservice.ContainerServiceSshPublicKeyArgs(
                key_data=ssh_pub_key,
            )],
        ),
    ),
    network_profile=containerservice.ContainerServiceNetworkProfileArgs(
        network_plugin="azure",
        network_policy="azure",
        service_cidr="10.96.0.0/16",
        dns_service_ip="10.96.0.10",
    ),
    resource_group_name=resource_group.name
)

# Build a Kubeconfig to access the cluster
creds = containerservice.list_managed_cluster_user_credentials_output(
    resource_group_name=resource_group.name,
    resource_name=managed_cluster.name,
)
encoded = creds.kubeconfigs[0].value
kubeconfig = encoded.apply(lambda enc: base64.b64decode(enc).decode())

# Export some values for use elsewhere
pulumi.export("rgname", resource_group.name)
pulumi.export("vnetName", virtual_network.name)
pulumi.export("clusterName", managed_cluster.name)
pulumi.export("kubeconfig", kubeconfig)
