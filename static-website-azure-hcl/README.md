# Static Website on Azure (Pulumi HCL)

A Pulumi HCL program that deploys a static website to Azure using a storage account's static website feature fronted by an Azure CDN.

## Overview

The website content in `./www` is uploaded to the `$web` container of a storage account that has static website hosting enabled. An Azure Front Door profile distributes and caches the content over HTTPS. (Front Door is the modern replacement for the now-retired classic Azure CDN.) The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AzureRM (`hashicorp/azurerm`)
- Random (`hashicorp/random`)

## Resources Created

- `random_string` (`suffix`): A random suffix used to build globally unique names.
- `azurerm_resource_group` (`resource_group`): The resource group for the website.
- `azurerm_storage_account` (`account`): A `StorageV2` account hosting the content.
- `azurerm_storage_account_static_website` (`website`): Enables static website hosting.
- `azurerm_storage_blob` (`files`): One blob per file under `path`, uploaded to `$web`.
- `azurerm_cdn_frontdoor_profile` / `azurerm_cdn_frontdoor_endpoint`: The Front Door profile and endpoint.
- `azurerm_cdn_frontdoor_origin_group` / `azurerm_cdn_frontdoor_origin` / `azurerm_cdn_frontdoor_route`: Route requests to the storage account's static website origin.

## Outputs

- **originURL** / **originHostname**: The storage account's static website endpoint and host.
- **cdnURL** / **cdnHostname**: The Front Door endpoint URL and hostname.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`).
- An Azure subscription with permissions to create storage and CDN resources. Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.

## Usage

```bash
pulumi new static-website-azure-hcl
pulumi up
```

A new CDN endpoint can take several minutes to propagate.

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

- Add a custom domain and managed certificate to the CDN endpoint.
- Configure CDN caching and compression rules.
- Wire up a build step to generate the contents of `www`.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
