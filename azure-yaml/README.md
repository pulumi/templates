# Azure Native (YAML) Pulumi Template

This repository provides a minimal Pulumi YAML template that provisions an Azure Resource Group and Storage Account using the Azure Native provider. It’s an ideal starting point for defining Azure infrastructure declaratively with YAML.

## Overview

**Providers used**
- Azure Native (`azure-native`)

**Resources created**
- `azure-native:resources:ResourceGroup` – a new Azure Resource Group
- `azure-native:storage:StorageAccount` – a StorageV2 account with Standard_LRS SKU

**Outputs returned**
- `primaryStorageKey` – the primary key of the created Storage Account

Use this template when you:
- Want to explore Azure Native via a YAML-first approach
- Need a minimal example to bootstrap your infrastructure
- Prefer a declarative, configuration-driven workflow

## Prerequisites

- An active Azure subscription
- Authentication set up via `az login` or appropriate `ARM_*` environment variables
- Pulumi CLI installed and authenticated

## Usage

Initialize a new project from this template:
```bash
pulumi new azure-yaml
```

Follow the interactive prompts to specify:
- **Project Name**
- **Description**
- **azure-native:location** (default: `WestUS2`)

Preview and deploy your stack:
```bash
pulumi up
```

## Project Layout

- `Pulumi.yaml` – the Pulumi YAML program defining resources, variables, and outputs
- `Pulumi.<stack>.yaml` – *(generated)* per-stack configuration values

## Configuration

| Key                     | Description                 | Default   |
|-------------------------|-----------------------------|-----------|
| `azure-native:location` | Azure region for resources  | `WestUS2` |

Override configuration values with:
```bash
pulumi config set azure-native:location eastus
```

## Outputs

| Name                | Description                               |
|---------------------|-------------------------------------------|
| `primaryStorageKey` | Primary access key of the Storage Account |

Retrieve outputs:
```bash
pulumi stack output primaryStorageKey
```

## Next Steps

- Customize or extend `Pulumi.yaml` to add more Azure resources (e.g., Databases, Functions, VMs)
- Explore advanced YAML features: loops, conditionals, and transforms
- Migrate to other Pulumi runtimes: TypeScript, Python, or Go
- Secure sensitive values with Pulumi secrets: https://www.pulumi.com/docs/intro/concepts/config/secrets/

## Getting Help

- Pulumi Documentation: https://www.pulumi.com/docs/
- Community Forum: https://pulumi.com/community
- Slack: https://slack.pulumi.com/
- GitHub Issues: https://github.com/pulumi/templates/issues