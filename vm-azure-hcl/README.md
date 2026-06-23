# Virtual Machine on Azure (Pulumi HCL)

A Pulumi HCL program that deploys an Azure Linux virtual machine running a simple web server.

## Overview

The program creates a resource group, a virtual network and subnet, a public IP with a DNS name, a network security group allowing HTTP and SSH, a network interface, and a Linux VM. An SSH key is generated for access, and a startup script serves a "Hello, world!" page. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AzureRM (`hashicorp/azurerm`)
- TLS (`hashicorp/tls`) — generates the SSH key
- Random (`hashicorp/random`)

## Resources Created

- `tls_private_key` (`ssh`): An SSH key pair for the VM.
- `azurerm_resource_group` / `azurerm_virtual_network` / `azurerm_subnet`: The resource group and network.
- `azurerm_public_ip` (`public_ip`): A static public IP with a DNS label.
- `azurerm_network_security_group` (+ association): Allows HTTP and SSH.
- `azurerm_network_interface` (`nic`): The VM's network interface.
- `azurerm_linux_virtual_machine` (`vm`): The virtual machine.

## Outputs

- **ip**: The VM's public IP address.
- **hostname**: The VM's fully-qualified domain name.
- **url**: The HTTP URL of the web server.
- **privatekey**: The generated SSH private key (sensitive).

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`). Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.

## Usage

```bash
pulumi new vm-azure-hcl
pulumi up
```

Open the `url` output once the VM has booted. To SSH in:

```bash
pulumi stack output privatekey --show-secrets > id_rsa && chmod 600 id_rsa
ssh -i id_rsa pulumiuser@$(pulumi stack output hostname)
```

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., location)
```

## Configuration

- **location**: The Azure region. Default: `WestUS2`.
- **admin_username**: The admin account name. Default: `pulumiuser`.
- **vm_name**: The VM name and DNS prefix. Default: `my-server`.
- **vm_size**: The VM size. Default: `Standard_A1_v2`.
- **os_image**: The image reference (`publisher:offer:sku:version`). Default: `Debian:debian-11:11:latest`.
- **service_port**: The HTTP port to serve on. Default: `80`.

## Next Steps

- Replace the inline startup script with your own application.
- Add a managed disk or attach Azure Files for storage.
- Put the VM behind a load balancer or Application Gateway.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
