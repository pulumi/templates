 # Azure Native TypeScript Pulumi Template

 This template provides a minimal, ready-to-go Pulumi program for deploying Azure resources using the Azure Native provider in TypeScript. It establishes a basic infrastructure stack that you can use as a foundation for more complex deployments.

 ## When to Use This Template

 - You need a quick boilerplate for Azure Native deployments with Pulumi and TypeScript
 - You want to create a Resource Group and Storage Account as a starting point
 - You’re exploring Pulumi’s Azure Native SDK and TypeScript support

 ## Prerequisites

 - An active Azure subscription
 - Node.js (LTS) installed
 - A Pulumi account and CLI already installed and configured
 - Azure credentials available (e.g., via `az login` or environment variables)

 ## Usage

 Scaffold a new project from the Pulumi registry template:
 ```bash
 pulumi new azure-typescript
 ```

 Follow the prompts to:
 1. Name your project and stack
 2. (Optionally) override the default Azure location

 Once the project is created:
 ```bash
 cd <your-project-name>
 pulumi config set azure-native:location <your-region>
 pulumi up
 ```

 ## Project Layout

 ```
 .
 ├── Pulumi.yaml       # Project metadata & template configuration
 ├── index.ts          # Main Pulumi program defining resources
 ├── package.json      # Node.js dependencies and project metadata
 └── tsconfig.json     # TypeScript compiler options
 ```

 ## Configuration

 Pulumi configuration lets you customize deployment parameters.

 - **azure-native:location** (string)
   - Description: Azure region to provision resources in
   - Default: `WestUS2`

 Set a custom location before deployment:
 ```bash
 pulumi config set azure-native:location eastus
 ```

 ## Resources Created

 1. **Resource Group**: A container for all other resources
 2. **Storage Account**: A StorageV2 account with Standard_LRS SKU

 ## Outputs

 After `pulumi up`, the following output is exported:
 - **storageAccountName**: The name of the created Storage Account

 Retrieve it with:
 ```bash
 pulumi stack output storageAccountName
 ```

 ## Next Steps

 - Extend this template by adding more Azure Native resources (e.g., Networking, App Services)
 - Modularize your stack with Pulumi Components for reusable architectures
 - Integrate with CI/CD pipelines (GitHub Actions, Azure DevOps, etc.)

 ## Getting Help

 If you have questions or run into issues:
 - Explore the Pulumi docs: https://www.pulumi.com/docs/
 - Join the Pulumi Community on Slack: https://pulumi-community.slack.com/
 - File an issue on the Pulumi Azure Native SDK GitHub: https://github.com/pulumi/pulumi-azure-native/issues