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
  default = "WestUS"
}

# The path to the folder containing the website
variable "path" {
  type    = string
  default = "./www"
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
  # Map file extensions to the content types used when uploading blobs.
  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".json" = "application/json"
    ".svg"  = "image/svg+xml"
    ".png"  = "image/png"
    ".jpg"  = "image/jpeg"
    ".jpeg" = "image/jpeg"
    ".gif"  = "image/gif"
    ".ico"  = "image/x-icon"
    ".txt"  = "text/plain"
  }
}

# A random suffix to make the storage account name globally unique.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a resource group for the website.
resource "azurerm_resource_group" "resource_group" {
  name     = "rg-static-website-${random_string.suffix.result}"
  location = var.location
}

# Create a blob storage account.
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
  for_each = fileset(var.path, "**")

  name                   = each.value
  storage_account_name   = azurerm_storage_account.account.name
  storage_container_name = "$web"
  type                   = "Block"
  source                 = "${var.path}/${each.value}"
  content_type           = lookup(local.mime_types, regex("\\.[^.]+$", each.value), "application/octet-stream")

  depends_on = [azurerm_storage_account_static_website.website]
}

# Create an Azure Front Door profile to distribute and cache the website.
# (Front Door is the modern replacement for the now-retired classic Azure CDN.)
resource "azurerm_cdn_frontdoor_profile" "profile" {
  name                = "fd-static-website-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  sku_name            = "Standard_AzureFrontDoor"
}

# Create a Front Door endpoint that serves the website over HTTPS.
resource "azurerm_cdn_frontdoor_endpoint" "endpoint" {
  name                     = "fd-endpoint-${random_string.suffix.result}"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.profile.id
}

# Group the storage account's static website as an origin.
resource "azurerm_cdn_frontdoor_origin_group" "origin_group" {
  name                     = "static-website"
  cdn_frontdoor_profile_id = azurerm_cdn_frontdoor_profile.profile.id

  load_balancing {}
}

# Point the origin at the storage account's static website host.
resource "azurerm_cdn_frontdoor_origin" "origin" {
  name                           = "storage-account"
  cdn_frontdoor_origin_group_id  = azurerm_cdn_frontdoor_origin_group.origin_group.id
  enabled                        = true
  certificate_name_check_enabled = true
  host_name                      = azurerm_storage_account.account.primary_web_host
  origin_host_header             = azurerm_storage_account.account.primary_web_host
  http_port                      = 80
  https_port                     = 443
  priority                       = 1
  weight                         = 1000
}

# Route all requests through the endpoint to the storage origin.
resource "azurerm_cdn_frontdoor_route" "route" {
  name                          = "default"
  cdn_frontdoor_endpoint_id     = azurerm_cdn_frontdoor_endpoint.endpoint.id
  cdn_frontdoor_origin_group_id = azurerm_cdn_frontdoor_origin_group.origin_group.id
  cdn_frontdoor_origin_ids      = [azurerm_cdn_frontdoor_origin.origin.id]

  supported_protocols    = ["Http", "Https"]
  patterns_to_match      = ["/*"]
  forwarding_protocol    = "HttpsOnly"
  https_redirect_enabled = true
  link_to_default_domain = true
}

# Export the URLs and hostnames of the storage account and distribution.
output "originURL" {
  value = azurerm_storage_account.account.primary_web_endpoint
}

output "originHostname" {
  value = azurerm_storage_account.account.primary_web_host
}

output "cdnURL" {
  value = "https://${azurerm_cdn_frontdoor_endpoint.endpoint.host_name}"
}

output "cdnHostname" {
  value = azurerm_cdn_frontdoor_endpoint.endpoint.host_name
}
