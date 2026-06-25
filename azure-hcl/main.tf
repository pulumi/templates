terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
  }
}

# Create an Azure resource group
resource "azure-native_resources_resource_group" "resource-group" {}

# Create an Azure storage account in the resource group
resource "azure-native_storage_storage_account" "sa" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  kind                = "StorageV2"
  sku = {
    name = "Standard_LRS"
  }
}

# Export the storage account name
output "storage_account_name" {
  value = azure-native_storage_storage_account.sa.name
}
