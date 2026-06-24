# Azure Storage Account (Pulumi HCL)

A minimal Pulumi HCL template that provisions an Azure resource group and storage account and exports the storage account name.

## Overview

This template uses the AzureRM provider to create a resource group and a `StorageV2` storage account. A random suffix keeps the storage account name globally unique. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AzureRM (`hashicorp/azurerm`)
- Random (`hashicorp/random`)

## Resources Created

- `azurerm_resource_group` (`resource_group`): The resource group that contains the storage account.
- `random_string` (`suffix`): A random suffix used to build a globally unique storage account name.
- `azurerm_storage_account` (`sa`): A Standard LRS `StorageV2` storage account.

## Outputs

- **storage_account_name**: The name of the created storage account.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`).
- An Azure subscription with permissions to create resource groups and storage accounts. The active subscription is used; set `ARM_SUBSCRIPTION_ID` to override it.

## Usage

Initialize a new project from this template by running:

```bash
pulumi new azure-hcl
```

You will be prompted for:
- A project name (default is set by the template).
- A project description.
- The Azure location to deploy into (default: `WestUS2`).

After initialization, deploy your stack:

```bash
pulumi up
```

## Project Layout

After `pulumi new`, your directory will look like:

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., location)
```

## Configuration

This template supports the following configuration keys:

- **location**: The Azure location to deploy resources into.
  - Default: `WestUS2`

To override the location, run:

```bash
pulumi config set location EastUS
```

## When to Use This Template

This template is ideal if you need:
- A lightweight starting point for an Azure storage account.
- To learn Pulumi with HCL programs.
- A quick bootstrap for small storage-focused projects.

## Next Steps

- Add blob containers, file shares, or queues to the storage account.
- Configure network rules, private endpoints, or lifecycle policies.
- Integrate with other Azure services (e.g., Functions, CDN).
- Explore additional Pulumi HCL examples.

## Getting Help

If you have questions or encounter issues:
- Visit the Pulumi documentation: https://www.pulumi.com/docs/
- Join the Pulumi Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
