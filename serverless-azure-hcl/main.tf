terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0.0"
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
  default = "WestUS"
}

# The path to the folder containing the website
variable "site_path" {
  type    = string
  default = "./www"
}

# The path to the folder containing the function to deploy
variable "app_path" {
  type    = string
  default = "./app"
}

# The file to use for top-level pages
variable "index_document" {
  type    = string
  default = "index.html"
}

# The file to use for error pages
variable "error_document" {
  type    = string
  default = "error.html"
}

locals {
  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".json" = "application/json"
    ".svg"  = "image/svg+xml"
    ".ico"  = "image/x-icon"
    ".txt"  = "text/plain"
  }
}

# A random suffix to make resource names globally unique.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a resource group for the application.
resource "azurerm_resource_group" "resource_group" {
  name     = "rg-serverless-${random_string.suffix.result}"
  location = var.location
}

# Create a storage account that backs both the website and the Function App.
resource "azurerm_storage_account" "account" {
  name                     = "sa${random_string.suffix.result}"
  resource_group_name      = azurerm_resource_group.resource_group.name
  location                 = azurerm_resource_group.resource_group.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
}

# Configure the storage account as a website.
resource "azurerm_storage_account_static_website" "website" {
  storage_account_id = azurerm_storage_account.account.id
  index_document     = var.index_document
  error_404_document = var.error_document
}

# Sync the website files to the account's $web container.
resource "azurerm_storage_blob" "files" {
  for_each = fileset(var.site_path, "**")

  name                   = each.value
  storage_account_name   = azurerm_storage_account.account.name
  storage_container_name = "$web"
  type                   = "Block"
  source                 = "${var.site_path}/${each.value}"
  content_type           = lookup(local.mime_types, regex("\\.[^.]+$", each.value), "application/octet-stream")

  depends_on = [azurerm_storage_account_static_website.website]
}

# Write a config file the website uses to find the API endpoint.
resource "azurerm_storage_blob" "config" {
  name                   = "config.json"
  storage_account_name   = azurerm_storage_account.account.name
  storage_container_name = "$web"
  type                   = "Block"
  content_type           = "application/json"
  source_content         = jsonencode({ api = "https://${azurerm_linux_function_app.app.default_hostname}/api" })

  depends_on = [azurerm_storage_account_static_website.website]
}

# Package the function source into a deployment archive.
data "archive_file" "app" {
  type        = "zip"
  source_dir  = var.app_path
  output_path = "${path.module}/app.zip"
}

# Create a Consumption (serverless) plan for the Function App.
resource "azurerm_service_plan" "plan" {
  name                = "plan-serverless-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

# Create the Function App and deploy the function code.
resource "azurerm_linux_function_app" "app" {
  name                = "fn-serverless-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  service_plan_id     = azurerm_service_plan.plan.id

  storage_account_name       = azurerm_storage_account.account.name
  storage_account_access_key = azurerm_storage_account.account.primary_access_key

  zip_deploy_file = data.archive_file.app.output_path

  app_settings = {
    SCM_DO_BUILD_DURING_DEPLOYMENT = "true"
    ENABLE_ORYX_BUILD              = "true"
  }

  site_config {
    application_stack {
      python_version = "3.11"
    }

    cors {
      allowed_origins = ["*"]
    }
  }
}

# Export the URLs of the website and serverless endpoint.
output "siteURL" {
  value = azurerm_storage_account.account.primary_web_endpoint
}

output "apiURL" {
  value = "https://${azurerm_linux_function_app.app.default_hostname}/api/data"
}
