terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
  }
}

variable "node_count" {
  description = "Number of worker nodes in the cluster"
  type        = number
  default     = 3
}

variable "dns_prefix" {
  description = "DNS prefix for the cluster"
  type        = string
  default     = "pulumi"
}

variable "node_vm_size" {
  description = "VM size to use for worker nodes in the cluster"
  type        = string
  default     = "Standard_DS2_v2"
}

# Create a new resource group
resource "azure-native_resources_resource_group" "resource-group" {}

# Create a managed cluster
resource "azure-native_containerservice_managed_cluster" "cluster" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  dns_prefix          = var.dns_prefix
  enable_rbac         = true

  identity = {
    type = "SystemAssigned"
  }

  agent_pool_profiles {
    name    = "systempool"
    count   = var.node_count
    mode    = "System"
    os_type = "Linux"
    vm_size = var.node_vm_size
  }
}

# Create a user Kubeconfig
# This SHOULD NOT be used for an explicit provider
# This SHOULD be used for user logins to the cluster
data "azure-native_containerservice_list_managed_cluster_user_credentials" "credentials" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  resource_name       = azure-native_containerservice_managed_cluster.cluster.name
}

# Export the Kubeconfig of the cluster
output "cluster_name" {
  value = azure-native_containerservice_managed_cluster.cluster.name
}

output "kubeconfig" {
  value     = base64decode(data.azure-native_containerservice_list_managed_cluster_user_credentials.credentials.kubeconfigs[0].value)
  sensitive = true
}
