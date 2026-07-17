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

variable "site_path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
}

variable "app_path" {
  description = "The path to the folder containing the functions to be deployed"
  type        = string
  default     = "./app"
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

locals {
  app_archive = fileArchive(var.app_path)
  config_json = jsonencode({ api = "https://${azure-native_web_web_app.app.default_host_name}/api" })
  config_file = stringAsset(local.config_json)
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

# Create a storage container for the pages of the website.
resource "azure-native_storage_storage_account_static_website" "website" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  account_name        = azure-native_storage_storage_account.account.name
  index_document      = var.index_document
  error404_document   = var.error_document
}

# Use a synced folder to manage the files of the website.
resource "synced-folder_azure_blob_folder" "synced-folder" {
  path                 = var.site_path
  resource_group_name  = azure-native_resources_resource_group.resource-group.name
  storage_account_name = azure-native_storage_storage_account.account.name
  container_name       = azure-native_storage_storage_account_static_website.website.container_name
}

# Create a storage container for the serverless app.
resource "azure-native_storage_blob_container" "app-container" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  account_name        = azure-native_storage_storage_account.account.name
  public_access       = "None"
}

# Upload the serverless app to the storage container.
resource "azure-native_storage_blob" "app-blob" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  account_name        = azure-native_storage_storage_account.account.name
  container_name      = azure-native_storage_blob_container.app-container.name
  type                = "Block"
  source              = local.app_archive
}

# Create a shared access signature allowing access to function storage.
data "azure-native_storage_list_storage_account_service_s_a_s" "signature" {
  resource_group_name       = azure-native_resources_resource_group.resource-group.name
  account_name              = azure-native_storage_storage_account.account.name
  protocols                 = "https"
  shared_access_start_time  = "2022-01-01"
  shared_access_expiry_time = "2030-01-01"
  resource                  = "c"
  permissions               = "r"
  canonicalized_resource    = "/blob/${azure-native_storage_storage_account.account.name}/${azure-native_storage_blob_container.app-container.name}"
  content_type              = "application/json"
  cache_control             = "max-age=5"
  content_disposition       = "inline"
  content_encoding          = "deflate"
}

# Create an App Service plan for the Function App.
resource "azure-native_web_app_service_plan" "plan" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  kind                = "Linux"
  reserved            = true
  sku = {
    name = "Y1"
    tier = "Dynamic"
  }
}

# Create the Function App.
resource "azure-native_web_web_app" "app" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  server_farm_id      = azure-native_web_app_service_plan.plan.id
  kind                = "FunctionApp"

  site_config = {
    app_settings = [
      {
        name  = "FUNCTIONS_WORKER_RUNTIME"
        value = "python"
      },
      {
        name  = "FUNCTIONS_EXTENSION_VERSION"
        value = "~4"
      },
      {
        name  = "WEBSITE_RUN_FROM_PACKAGE"
        value = "https://${azure-native_storage_storage_account.account.name}.blob.core.windows.net/${azure-native_storage_blob_container.app-container.name}/${azure-native_storage_blob.app-blob.name}?${data.azure-native_storage_list_storage_account_service_s_a_s.signature.service_sas_token}"
      },
    ]
    cors = {
      allowed_origins = ["*"]
    }
  }
}

# Create a JSON configuration file for the website.
resource "azure-native_storage_blob" "config" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  account_name        = azure-native_storage_storage_account.account.name
  container_name      = azure-native_storage_storage_account_static_website.website.container_name
  content_type        = "application/json"
  type                = "Block"
  source              = local.config_file
}

# Export the URLs of the website and serverless endpoint.
output "origin_url" {
  value = azure-native_storage_storage_account.account.primary_endpoints.web
}

output "api_url" {
  value = "https://${azure-native_web_web_app.app.default_host_name}/api"
}
