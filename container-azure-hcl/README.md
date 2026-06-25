# Containerized Service on Azure (Pulumi HCL)

A Pulumi HCL program that builds a container image and runs it on Azure Container Instances.

## Overview

The application in `./app` is built into a container image and pushed to an Azure Container Registry, then deployed as a publicly accessible container group on Azure Container Instances. The image is built and pushed with the `docker-build` provider, so a running Docker daemon is required at deploy time. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Azure Native (`pulumi/azure-native`)
- Docker Build (`pulumi/docker-build`) — builds and pushes the image
- Random (`pulumi/random`) — a unique DNS-name suffix

## Resources Created

- `random_random_string` (`dns-name`): A suffix giving the service a unique DNS name.
- `azure-native_resources_resource_group` (`resource-group`): The resource group.
- `azure-native_containerregistry_registry` (`registry`): Stores the application image.
- `data azure-native_containerregistry_list_registry_credentials` (`credentials`): Fetches the registry's admin login credentials.
- `docker-build_image` (`image`): Builds and pushes the image to the registry.
- `azure-native_containerinstance_container_group` (`container-group`): Runs the container with a public IP and DNS name.

## Outputs

- **hostname**: The fully-qualified domain name of the container group.
- **ip**: The public IP address.
- **url**: The full HTTP URL of the service.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`). Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.
- A running Docker daemon (the image is built locally and pushed to the registry).

## Usage

```bash
pulumi new container-azure-hcl
pulumi up
```

Open the `url` output once the container group is running.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── app/                  # The container application (Dockerfile, source)
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., location)
```

## Configuration

- **location**: The Azure region to deploy into. Default: `WestUS`.
- **app_path**: The container application folder. Default: `./app`.
- **image_name**: The image name. Default: `my-app`.
- **container_port**: The container port. Default: `80`.
- **cpu** / **memory**: CPU cores and memory (GB). Defaults: `1` / `2`.

## Next Steps

- Front the container group with Application Gateway or Front Door for HTTPS.
- Move to Azure Container Apps for autoscaling and revisions.
- Add a managed identity for registry access instead of the admin user.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
