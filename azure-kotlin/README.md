# Minimal Azure Kotlin Pulumi Template

This template provides a minimal Pulumi program written in Kotlin that provisions an Azure Storage Account using the Pulumi Azure Native provider. It's a great starting point for building Azure infrastructure with Pulumi and Kotlin.

## Providers
- Azure Native (pulumi/azure-native)

## Resources
- `azure-native.resources.ResourceGroup`: An Azure Resource Group.
- `azure-native.storage.StorageAccount`: An Azure Storage Account resource with Standard LRS redundancy.

## Outputs
- `storageAccountName`: The name of the created Azure Storage Account.

## When to use this template
Use this template if you:
- Want a quick start with Pulumi in Kotlin
- Need an example of provisioning basic Azure resources
- Are familiar with Gradle and Kotlin development

## Prerequisites
- Java Development Kit (JDK) 21 or later
- Azure credentials configured (via Azure CLI `az login`, service principal, or managed identity)

## Getting Started
1. Create a new Pulumi project from this template:
```bash
   pulumi new azure-kotlin
```
2. Follow the interactive prompts to set your project name, description, and Azure location (default: `WestUS`).
3. Change into your project directory:
```bash
   cd <project-name>
```
4. Deploy your stack:
```bash
   pulumi up
```

## Configuration
This template supports the following configuration values:
- `azure-native:location` (string) — Azure location to deploy into. Default: `WestUS`.

View or set configuration values:
```bash
pulumi config
pulumi config set azure-native:location EastUS
```

## Getting Help
- Pulumi documentation: https://www.pulumi.com/docs/
- Azure Native provider reference: https://www.pulumi.com/registry/packages/azure-native/
- Pulumi Community Slack: https://slack.pulumi.com/
- Stack Overflow (`pulumi` tag)
- Report issues: https://github.com/pulumi/pulumi/issues