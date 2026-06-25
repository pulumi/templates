terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
    docker-build = {
      source = "pulumi/docker-build"
    }
    random = {
      source = "pulumi/random"
    }
  }
}

variable "app_path" {
  description = "The path to the container application to deploy"
  type        = string
  default     = "./app"
}

variable "image_name" {
  description = "The name to give the container image"
  type        = string
  default     = "my-app"
}

variable "image_tag" {
  description = "The tag to give the container image"
  type        = string
  default     = "latest"
}

variable "container_port" {
  description = "The port to expose on the container"
  type        = number
  default     = 80
}

variable "cpu" {
  description = "The number of CPU cores to allocate for the container"
  type        = number
  default     = 1
}

variable "memory" {
  description = "The amount of memory, in GB, to allocate for the container"
  type        = number
  default     = 2
}

# A random suffix to give the service a unique DNS name.
resource "random_random_string" "dns-name" {
  length  = 8
  upper   = false
  special = false
}

# Create a resource group for the container registry.
resource "azure-native_resources_resource_group" "resource-group" {
}

# Create a container registry with the admin user enabled.
resource "azure-native_containerregistry_registry" "registry" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  admin_user_enabled  = true
  sku = {
    name = "Basic"
  }
}

# Fetch login credentials for the registry.
data "azure-native_containerregistry_list_registry_credentials" "credentials" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  registry_name       = azure-native_containerregistry_registry.registry.name
}

# Build the container image and push it to the registry.
resource "docker-build_image" "image" {
  tags      = ["${azure-native_containerregistry_registry.registry.login_server}/${var.image_name}:${var.image_tag}"]
  platforms = ["linux/amd64"]
  push      = true

  context = {
    location = var.app_path
  }

  registries {
    address  = azure-native_containerregistry_registry.registry.login_server
    username = data.azure-native_containerregistry_list_registry_credentials.credentials.username
    password = data.azure-native_containerregistry_list_registry_credentials.credentials.passwords[0].value
  }
}

# Deploy the image as a publicly accessible container group.
resource "azure-native_containerinstance_container_group" "container-group" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  os_type             = "Linux"
  restart_policy      = "Always"

  image_registry_credentials {
    server   = azure-native_containerregistry_registry.registry.login_server
    username = data.azure-native_containerregistry_list_registry_credentials.credentials.username
    password = data.azure-native_containerregistry_list_registry_credentials.credentials.passwords[0].value
  }

  containers {
    name  = var.image_name
    image = docker-build_image.image.ref

    ports {
      port     = var.container_port
      protocol = "TCP"
    }

    environment_variables {
      name  = "PORT"
      value = tostring(var.container_port)
    }

    resources = {
      requests = {
        cpu           = var.cpu
        memory_in_g_b = var.memory
      }
    }
  }

  ip_address = {
    type           = "Public"
    dns_name_label = "${var.image_name}-${random_random_string.dns-name.result}"
    ports = [{
      port     = var.container_port
      protocol = "TCP"
    }]
  }
}

# Export the service's hostname, IP address, and URL.
output "hostname" {
  value = azure-native_containerinstance_container_group.container-group.ip_address.fqdn
}

output "ip" {
  value = azure-native_containerinstance_container_group.container-group.ip_address.ip
}

output "url" {
  value = "http://${azure-native_containerinstance_container_group.container-group.ip_address.fqdn}:${var.container_port}"
}
