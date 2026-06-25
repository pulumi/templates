# Static Website on Azure (Pulumi HCL)

A Pulumi HCL program that deploys a static website to Azure using a storage account's static website feature.

## Overview

The website content in `./www` is uploaded to the `$web` container of a storage account that has static website hosting enabled, and the site is served directly from the storage account's static website endpoint. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

To put a CDN in front of the site, add an Azure Front Door profile (Front Door is the modern replacement for the now-retired classic Azure CDN); `main.tf` includes a code comment outlining how.

## Providers

- Azure Native (`pulumi/azure-native`)
- Synced Folder (`pulumi/synced-folder`) — uploads the website folder to the `$web` container

## Resources Created

- `azure-native_resources_resource_group` (`resource-group`): The resource group for the website.
- `azure-native_storage_storage_account` (`account`): A `StorageV2` account hosting the content.
- `azure-native_storage_storage_account_static_website` (`website`): Enables static website hosting on the account.
- `synced-folder_azure_blob_folder` (`synced-folder`): Syncs the contents of `path` to the account's `$web` container.

## Outputs

- **origin_url** / **origin_hostname**: The storage account's static website endpoint and host.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`).
- An Azure subscription with permissions to create storage resources. Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.

## Usage

```bash
pulumi new static-website-azure-hcl
pulumi up
```

Open the `origin_url` output once the deployment finishes.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── www/                  # Website content (index.html, error.html)
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., location)
```

## Configuration

- **location**: The Azure region to deploy into. Default: `WestUS`.
- **path**: The folder containing the website content. Default: `./www`.
- **index_document**: The top-level page. Default: `index.html`.
- **error_document**: The error page. Default: `error.html`.

```bash
pulumi config set location EastUS
```

## Next Steps

- Add an Azure Front Door profile to serve the site over HTTPS with caching and a custom domain.
- Wire up a build step to generate the contents of `www`.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
