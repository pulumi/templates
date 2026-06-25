terraform {
  required_providers {
    azure-native = {
      source = "pulumi/azure-native"
    }
    tls = {
      source = "pulumi/tls"
    }
    random = {
      source = "pulumi/random"
    }
  }
}

variable "admin_username" {
  description = "The user account to create on the VM"
  type        = string
  default     = "pulumiuser"
}

variable "vm_name" {
  description = "The computer name and DNS hostname prefix to use for the VM"
  type        = string
  default     = "my-server"
}

variable "vm_size" {
  description = "The machine size to use for the VM"
  type        = string
  default     = "Standard_A1_v2"
}

variable "os_image" {
  description = "The Azure image reference (publisher:offer:sku:version) to use for the VM"
  type        = string
  default     = "Debian:debian-11:11:latest"
}

variable "service_port" {
  description = "The HTTP service port to expose on the VM"
  type        = number
  default     = 80
}

locals {
  os_image_parts = split(":", var.os_image)
  dns_name       = "${var.vm_name}-${random_random_string.random-string.result}"

  init_script = base64encode(<<-EOF
    #!/bin/bash
    echo '<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Hello, world!</title>
    </head>
    <body>
        <h1>Hello, world! 👋</h1>
        <p>Deployed with 💜 by <a href="https://pulumi.com/">Pulumi</a>.</p>
    </body>
    </html>' > index.html
    sudo python3 -m http.server ${var.service_port} &
  EOF
  )
}

# Create an SSH key for the VM.
resource "tls_private_key" "ssh-key" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

# Use a random string to give the VM a unique DNS name.
resource "random_random_string" "random-string" {
  length  = 8
  upper   = false
  special = false
}

# Create a resource group.
resource "azure-native_resources_resource_group" "resource-group" {}

# Create a virtual network with a subnet.
resource "azure-native_network_virtual_network" "network" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  address_space = {
    address_prefixes = ["10.0.0.0/16"]
  }
  subnets {
    name           = "default"
    address_prefix = "10.0.1.0/24"
  }
}

# Create a public IP address with a DNS name.
resource "azure-native_network_public_i_p_address" "public-ip" {
  resource_group_name         = azure-native_resources_resource_group.resource-group.name
  public_ip_allocation_method = "Dynamic"
  dns_settings = {
    domain_name_label = local.dns_name
  }
}

# Allow inbound access over the service port (HTTP) and port 22 (SSH).
resource "azure-native_network_network_security_group" "security-group" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  security_rules {
    name                       = "${var.vm_name}-securityrule"
    priority                   = 1000
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
    destination_port_ranges    = [tostring(var.service_port), "22"]
  }
}

# Create a network interface bound to the subnet, public IP, and security group.
resource "azure-native_network_network_interface" "network-interface" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  network_security_group = {
    id = azure-native_network_network_security_group.security-group.id
  }
  ip_configurations {
    name                         = "${var.vm_name}-ipconfiguration"
    private_ip_allocation_method = "Dynamic"
    subnet = {
      id = azure-native_network_virtual_network.network.subnets[0].id
    }
    public_ip_address = {
      id = azure-native_network_public_i_p_address.public-ip.id
    }
  }
}

# Create the virtual machine.
resource "azure-native_compute_virtual_machine" "vm" {
  resource_group_name = azure-native_resources_resource_group.resource-group.name
  network_profile = {
    network_interfaces = [{
      id      = azure-native_network_network_interface.network-interface.id
      primary = true
    }]
  }
  hardware_profile = {
    vm_size = var.vm_size
  }
  os_profile = {
    computer_name  = var.vm_name
    admin_username = var.admin_username
    custom_data    = local.init_script
    linux_configuration = {
      disable_password_authentication = true
      ssh = {
        public_keys = [{
          key_data = tls_private_key.ssh-key.public_key_openssh
          path     = "/home/${var.admin_username}/.ssh/authorized_keys"
        }]
      }
    }
  }
  storage_profile = {
    os_disk = {
      name          = "${var.vm_name}-osdisk"
      create_option = "FromImage"
    }
    image_reference = {
      publisher = local.os_image_parts[0]
      offer     = local.os_image_parts[1]
      sku       = local.os_image_parts[2]
      version   = local.os_image_parts[3]
    }
  }
}

# Look up the VM's allocated public IP details once the machine is running.
data "azure-native_network_public_i_p_address" "address" {
  resource_group_name    = azure-native_resources_resource_group.resource-group.name
  public_ip_address_name = azure-native_network_public_i_p_address.public-ip.name
  expand                 = azure-native_compute_virtual_machine.vm.id
}

# Export the VM's public IP address, hostname, URL, and SSH private key.
output "ip" {
  value = data.azure-native_network_public_i_p_address.address.ip_address
}

output "hostname" {
  value = data.azure-native_network_public_i_p_address.address.dns_settings.fqdn
}

output "url" {
  value = "http://${data.azure-native_network_public_i_p_address.address.dns_settings.fqdn}:${var.service_port}"
}

output "private_key" {
  value     = tls_private_key.ssh-key.private_key_openssh
  sensitive = true
}
