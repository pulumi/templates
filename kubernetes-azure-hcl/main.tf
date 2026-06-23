terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
}

provider "azurerm" {
  features {}
}

# The Azure region to deploy into
variable "location" {
  type    = string
  default = "westus2"
}

# The number of worker nodes in the cluster
variable "node_count" {
  type    = number
  default = 3
}

# The DNS prefix to use for the cluster
variable "dns_prefix" {
  type    = string
  default = "pulumi"
}

# The VM size to use for worker nodes
variable "node_vm_size" {
  type    = string
  default = "Standard_DS2_v2"
}

# A random suffix to make names unique.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a resource group for the cluster.
resource "azurerm_resource_group" "resource_group" {
  name     = "rg-aks-${random_string.suffix.result}"
  location = var.location
}

# Create a managed Kubernetes (AKS) cluster.
resource "azurerm_kubernetes_cluster" "cluster" {
  name                = "aks-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  dns_prefix          = var.dns_prefix

  default_node_pool {
    name       = "systempool"
    node_count = var.node_count
    vm_size    = var.node_vm_size
  }

  identity {
    type = "SystemAssigned"
  }
}

# Export the cluster name and kubeconfig.
output "clusterName" {
  value = azurerm_kubernetes_cluster.cluster.name
}

output "kubeconfig" {
  value     = azurerm_kubernetes_cluster.cluster.kube_config_raw
  sensitive = true
}
