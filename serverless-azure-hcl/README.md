# Serverless Application on Azure (Pulumi HCL)

A Pulumi HCL program that deploys a serverless application on Azure: a Python Function App on a Consumption plan with a static front-end hosted on a storage account.

## Overview

A Python Azure Function returns the current time. A static website in `./www` is hosted from the storage account's static website feature and calls the function (the function endpoint is injected into a `config.json` the page reads at load time). The function code in `./app` is zipped, uploaded to a blob container, and run from that package via `WEBSITE_RUN_FROM_PACKAGE` (a service SAS grants the Function App read access). The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Azure Native (`pulumi/azure-native`)
- Synced Folder (`pulumi/synced-folder`) — uploads the website folder to the `$web` container

## Resources Created

- `azure-native_resources_resource_group` (`resource-group`): The resource group.
- `azure-native_storage_storage_account` (`account`) + `azure-native_storage_storage_account_static_website` (`website`): The account hosting both the website and the Function App package.
- `synced-folder_azure_blob_folder` (`synced-folder`): Syncs the website content to the `$web` container.
- `azure-native_storage_blob_container` (`app-container`) + `azure-native_storage_blob` (`app-blob`): A private container holding the zipped function package.
- `data azure-native_storage_list_storage_account_service_s_a_s` (`signature`): A service SAS granting read access to the package container.
- `azure-native_web_app_service_plan` (`plan`): A Linux Consumption (`Y1`) plan.
- `azure-native_web_web_app` (`app`): The Python Function App, run from the uploaded package.
- `azure-native_storage_blob` (`config`): The website's `config.json` pointing at the function endpoint.

## Outputs

- **origin_url**: The URL of the static website.
- **api_url**: The URL of the function endpoint (`/api`).

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`).
- An Azure subscription with permissions for storage, App Service, and Functions. Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.

## Usage

```bash
pulumi new serverless-azure-hcl
pulumi up
```

Open the `origin_url` output and click the button. (The Function App can take a few minutes to finish its first deployment and build.)

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── app/                  # The Function App source (data/, host.json, requirements.txt)
├── www/                  # Static front-end
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., location)
```

## Configuration

- **location**: The Azure region to deploy into. Default: `WestUS`.
- **site_path** / **app_path**: The website and function source folders.
- **index_document** / **error_document**: The website's page documents.

## Next Steps

- Add more functions and routes to the Function App.
- Put Azure Front Door in front of the static website.
- Add Application Insights for monitoring.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
