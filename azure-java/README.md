 # Pulumi Template: Azure Native Java Storage

 A minimal Pulumi template for creating a Resource Group and Storage Account on Azure using Java and the Azure Native provider. This template provisions:

 - An Azure Resource Group.
 - An Azure Storage Account (Standard_LRS, StorageV2).
 - Exports the primary storage account key as a secret output.

 ## Prerequisites

 - An Azure account with credentials configured (for example, via `az login` or setting environment variables `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_TENANT_ID`, and `ARM_SUBSCRIPTION_ID`).
 - Java 11 or higher installed.
 - Maven installed.
 - Pulumi CLI installed and logged in.

 ## Getting Started

 To create a new project from this template, run:

 ```bash
 pulumi new azure-java
 ```

 Follow the interactive prompts:

 - Project name
 - Project description
 - `azure-native:location`: The Azure location to use (default: WestUS2)

 Then, change into your project directory and preview or deploy your stack:

 ```bash
 cd <project-directory>
 pulumi up
 ```

 ## Project Layout

 ```plaintext
 .
 ├── Pulumi.yaml         # Project and template metadata
 ├── pom.xml             # Maven project file with dependencies
 └── src
     └── main
         └── java
             └── myproject
                 └── App.java  # Main program defining Azure resources
 ```

 ## Configuration

 | Key                       | Description                     | Default  |
 |---------------------------|---------------------------------|----------|
 | `azure-native:location`   | Azure region for resources      | WestUS2  |

 To override the default location, run:

 ```bash
 pulumi config set azure-native:location <your-region>
 ```

 ## Outputs

 - `primaryStorageKey` (Secret): The primary key for the Storage Account.

 ## Next Steps

 - Extend `App.java` to add more Azure resources (for example, Cosmos DB, Functions, or Networking).
 - Use multiple Pulumi stacks for different environments (development, staging, production).
 - Integrate Pulumi into your CI/CD pipeline.
 - Explore the Pulumi Azure Native SDK in the [Pulumi Registry](https://www.pulumi.com/registry/packages/azure-native/).

 ## Getting Help

 If you have questions or encounter any issues:

 - Check out the [Pulumi Documentation](https://www.pulumi.com/docs/).
 - Join the [Pulumi Community Slack](https://slack.pulumi.com/) for support.
 - File an issue in this repository.