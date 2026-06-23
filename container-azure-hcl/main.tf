terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.0"
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

# The path to the container application to deploy
variable "app_path" {
  type    = string
  default = "./app"
}

# The name to give the container image
variable "image_name" {
  type    = string
  default = "my-app"
}

# The port to expose on the container
variable "container_port" {
  type    = number
  default = 80
}

# The number of CPU cores to allocate for the container
variable "cpu" {
  type    = number
  default = 1
}

# The amount of memory, in GB, to allocate for the container
variable "memory" {
  type    = number
  default = 2
}

# A random suffix to make names globally unique.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a resource group for the service.
resource "azurerm_resource_group" "resource_group" {
  name     = "rg-container-${random_string.suffix.result}"
  location = var.location
}

# Create a container registry.
resource "azurerm_container_registry" "registry" {
  name                = "acr${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  sku                 = "Basic"
  admin_enabled       = true
}

# Authenticate the Docker provider to the registry.
provider "docker" {
  registry_auth {
    address  = azurerm_container_registry.registry.login_server
    username = azurerm_container_registry.registry.admin_username
    password = azurerm_container_registry.registry.admin_password
  }
}

# Build the container image from the application source.
resource "docker_image" "app" {
  name = "${azurerm_container_registry.registry.login_server}/${var.image_name}:latest"

  build {
    context  = var.app_path
    platform = "linux/amd64"
  }
}

# Push the image to the registry.
resource "docker_registry_image" "app" {
  name          = docker_image.app.name
  keep_remotely = true
}

# Deploy the image as a publicly accessible container group.
resource "azurerm_container_group" "container_group" {
  name                = "container-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  os_type             = "Linux"
  restart_policy      = "Always"
  ip_address_type     = "Public"
  dns_name_label      = "${var.image_name}-${random_string.suffix.result}"

  image_registry_credential {
    server   = azurerm_container_registry.registry.login_server
    username = azurerm_container_registry.registry.admin_username
    password = azurerm_container_registry.registry.admin_password
  }

  container {
    name   = var.image_name
    image  = docker_registry_image.app.name
    cpu    = var.cpu
    memory = var.memory

    ports {
      port     = var.container_port
      protocol = "TCP"
    }

    environment_variables = {
      PORT = tostring(var.container_port)
    }
  }
}

# Export the service's hostname, IP address, and URL.
output "hostname" {
  value = azurerm_container_group.container_group.fqdn
}

output "ip" {
  value = azurerm_container_group.container_group.ip_address
}

output "url" {
  value = "http://${azurerm_container_group.container_group.fqdn}:${var.container_port}"
}
