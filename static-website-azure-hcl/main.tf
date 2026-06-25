terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
    synced-folder = {
      source = "pulumi/synced-folder"
    }
  }
}

variable "path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
}

variable "index_document" {
  description = "The file to use for top-level pages"
  type        = string
  default     = "index.html"
}

variable "error_document" {
  description = "The file to use for error pages"
  type        = string
  default     = "error.html"
}

# Create a resource group for the website.
resource "azure-native_resources_resource_group" "resource-group" {}

# Create a blob storage account.
resource "azure-native_storage_storage_account" "account" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  kind                = "StorageV2"
  sku = {
    name = "Standard_LRS"
  }
}

# Configure the storage account as a website.
resource "azure-native_storage_storage_account_static_website" "website" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  account_name        = azure-native_storage_storage_account.account.name
  index_document      = var.index_document
  error404_document   = var.error_document
}

# Use a synced folder to manage the files of the website.
resource "synced-folder_azure_blob_folder" "synced-folder" {
  path                 = var.path
  resource_group_name  = azure-native_resources_resource_group.resource-group.name
  storage_account_name = azure-native_storage_storage_account.account.name
  container_name       = azure-native_storage_storage_account_static_website.website.container_name
}

# Export the URL and hostname of the storage account's website.
output "origin_url" {
  value = azure-native_storage_storage_account.account.primary_endpoints.web
}

output "origin_hostname" {
  value = split("/", azure-native_storage_storage_account.account.primary_endpoints.web)[2]
}
