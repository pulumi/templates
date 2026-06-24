terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 4.0.0"
    }
    tls = {
      source  = "hashicorp/tls"
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

variable "location" {
  description = "The Azure location to deploy into"
  type        = string
  default     = "WestUS2"
}

variable "admin_username" {
  description = "The user account to create on the VM"
  type        = string
  default     = "pulumiuser"
}

variable "vm_name" {
  description = "The DNS hostname prefix and computer name to use for the VM"
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

  init_script = <<-EOF
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
}

# Create an SSH key for the VM.
resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

# A random suffix to give the VM a unique DNS name.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a resource group.
resource "azurerm_resource_group" "resource_group" {
  name     = "rg-${var.vm_name}-${random_string.suffix.result}"
  location = var.location
}

# Create a virtual network and subnet.
resource "azurerm_virtual_network" "network" {
  name                = "${var.vm_name}-net"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "subnet" {
  name                 = "default"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.network.name
  address_prefixes     = ["10.0.1.0/24"]
}

# Create a public IP address with a DNS name.
resource "azurerm_public_ip" "public_ip" {
  name                = "${var.vm_name}-ip"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "${var.vm_name}-${random_string.suffix.result}"
}

# Allow inbound access over the service port (HTTP) and port 22 (SSH).
resource "azurerm_network_security_group" "security_group" {
  name                = "${var.vm_name}-nsg"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location

  security_rule {
    name                       = "allow-http-ssh"
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

# Create a network interface and associate the security group with it.
resource "azurerm_network_interface" "nic" {
  name                = "${var.vm_name}-nic"
  resource_group_name = azurerm_resource_group.resource_group.name
  location            = azurerm_resource_group.resource_group.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.public_ip.id
  }
}

resource "azurerm_network_interface_security_group_association" "nic_nsg" {
  network_interface_id      = azurerm_network_interface.nic.id
  network_security_group_id = azurerm_network_security_group.security_group.id
}

# Create the virtual machine.
resource "azurerm_linux_virtual_machine" "vm" {
  name                  = var.vm_name
  resource_group_name   = azurerm_resource_group.resource_group.name
  location              = azurerm_resource_group.resource_group.location
  size                  = var.vm_size
  admin_username        = var.admin_username
  computer_name         = var.vm_name
  network_interface_ids = [azurerm_network_interface.nic.id]
  custom_data           = base64encode(local.init_script)

  admin_ssh_key {
    username   = var.admin_username
    public_key = tls_private_key.ssh.public_key_openssh
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = local.os_image_parts[0]
    offer     = local.os_image_parts[1]
    sku       = local.os_image_parts[2]
    version   = local.os_image_parts[3]
  }
}

# Export the VM's public IP address, hostname, URL, and SSH private key.
output "ip" {
  value = azurerm_public_ip.public_ip.ip_address
}

output "hostname" {
  value = azurerm_public_ip.public_ip.fqdn
}

output "url" {
  value = "http://${azurerm_public_ip.public_ip.fqdn}:${var.service_port}"
}

output "private_key" {
  value     = tls_private_key.ssh.private_key_openssh
  sensitive = true
}
