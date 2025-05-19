# Azure Native Go Pulumi Template

A minimal Pulumi template for Go using the Azure Native provider. This project provisions an Azure Resource Group and a Storage Account, then exports the primary access key. It serves as a solid starting point for Go-based Infrastructure as Code on Azure.

## Prerequisites
- Go 1.20 or later
- Pulumi CLI (configured and logged in)
- An active Azure subscription with credentials set (for example via Azure CLI or environment variables)

## Creating a New Project
Use this template to bootstrap a new Pulumi project:
```
pulumi new azure-go
```
Follow the interactive prompts to name your project, choose a stack, and confirm configuration values.

## Project Layout
```
.
├── Pulumi.yaml       # Project and template metadata
├── go.mod            # Go module declaration
├── go.sum            # Module checksums
└── main.go           # Pulumi program defining resources and exports
```

## Configuration
Pulumi configuration controls the deployment. Available settings:
- `azure-native:location` (string) – Azure region for your resources (default: `WestUS2`)

Set a custom location:
```
pulumi config set azure-native:location <region>
```

## Deploying
After creating your project, preview and deploy your infrastructure:
```
pulumi preview   # see what will be created or changed
pulumi up        # apply changes to your Azure subscription
```

## Resources
This template creates:
- An Azure Resource Group
- An Azure Storage Account (Standard_LRS, StorageV2)

## Outputs
- `primaryStorageKey` – the primary access key for the Storage Account

## When to Use This Template
Use this template when you want a quick, Go-based Pulumi project on Azure Native that sets up a resource group and storage account. It’s ideal for:
- Learning Pulumi with Go and Azure Native
- Bootstrapping backend storage for applications
- Extending with additional Azure services

## Next Steps
- Customize `main.go` to add or modify resources
- Add tags, networking, or compute resources (VMs, Functions, etc.)
- Integrate with CI/CD pipelines
- Explore the Pulumi Go SDK: https://www.pulumi.com/docs/reference/pkg/go/
- Review Azure Native provider docs: https://www.pulumi.com/docs/intro/cloud-providers/azure-native/

## Getting Help
- File an issue on the project’s GitHub repository
- Join the Pulumi Community Slack: #azure or #general
- Browse the Pulumi documentation: https://www.pulumi.com/docs/
  