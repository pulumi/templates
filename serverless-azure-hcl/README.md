# Serverless Application on Azure (Pulumi HCL)

A Pulumi HCL program that deploys a serverless application on Azure: a Python Function App on a Consumption plan with a static front-end hosted on a storage account.

## Overview

A Python Azure Function returns the current time. A static website in `./www` is hosted from the storage account's static website feature and calls the function (the function endpoint is injected into a `config.json` the page reads at load time). The function code in `./app` is packaged and deployed to the Function App via zip deploy. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AzureRM (`hashicorp/azurerm`)
- Archive (`hashicorp/archive`) — packages the function source
- Random (`hashicorp/random`)

## Resources Created

- `azurerm_resource_group` (`resource_group`): The resource group.
- `azurerm_storage_account` (`account`) + `azurerm_storage_account_static_website` (`website`): The account hosting both the website and the Function App runtime.
- `azurerm_storage_blob` (`files`, `config`): The website content and its `config.json`.
- `azurerm_service_plan` (`plan`): A Linux Consumption (`Y1`) plan.
- `azurerm_linux_function_app` (`app`): The Python Function App, deployed from `./app`.

## Outputs

- **site_url**: The URL of the static website.
- **api_url**: The URL of the function endpoint (`/api/data`).

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`).
- An Azure subscription with permissions for storage, App Service, and Functions. Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.

## Usage

```bash
pulumi new serverless-azure-hcl
pulumi up
```

Open the `siteURL` output and click the button. (The Function App can take a few minutes to finish its first deployment and build.)

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
