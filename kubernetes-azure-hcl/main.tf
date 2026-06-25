terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
  }
}

variable "node_count" {
  description = "The number of worker nodes in the cluster"
  type        = number
  default     = 3
}

variable "dns_prefix" {
  description = "The DNS prefix to use for the cluster"
  type        = string
  default     = "pulumi"
}

variable "node_vm_size" {
  description = "The VM size to use for worker nodes"
  type        = string
  default     = "Standard_DS2_v2"
}

# Create a resource group for the cluster.
resource "azure-native_resources_resource_group" "resource-group" {
}

# Create a managed Kubernetes (AKS) cluster.
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

# Fetch the cluster's user credentials so we can export a kubeconfig.
data "azure-native_containerservice_list_managed_cluster_user_credentials" "credentials" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  resource_name       = azure-native_containerservice_managed_cluster.cluster.name
}

# Export the cluster name and kubeconfig.
output "cluster_name" {
  value = azure-native_containerservice_managed_cluster.cluster.name
}

output "kubeconfig" {
  value     = base64decode(data.azure-native_containerservice_list_managed_cluster_user_credentials.credentials.kubeconfigs[0].value)
  sensitive = true
}
