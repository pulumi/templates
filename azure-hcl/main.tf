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

# The Azure location to deploy into
variable "location" {
  type    = string
  default = "WestUS2"
}

# Create an Azure Resource Group
resource "azurerm_resource_group" "resource_group" {
  name     = "rg-${random_string.suffix.result}"
  location = var.location
}

# A random suffix to make the storage account name globally unique
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create an Azure Storage Account
resource "azurerm_storage_account" "sa" {
  name                     = "sa${random_string.suffix.result}"
  resource_group_name      = azurerm_resource_group.resource_group.name
  location                 = azurerm_resource_group.resource_group.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
}

# Export the storage account name
output "storageAccountName" {
  value = azurerm_storage_account.sa.name
}
