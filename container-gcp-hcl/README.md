# Containerized Service on Google Cloud (Pulumi HCL)

A Pulumi HCL program that builds a container image and runs it on Google Cloud Run.

## Overview

The application in `./app` is built into a container image and pushed to Artifact Registry, then deployed as a publicly accessible Cloud Run service. The image is built and pushed with the `docker-build` provider, so a running Docker daemon is required at deploy time. Configure Docker auth for Artifact Registry first (e.g. `gcloud auth configure-docker <region>-docker.pkg.dev`). The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google Cloud (`pulumi/gcp`)
- Docker Build (`pulumi/docker-build`) — builds and pushes the image
- Random (`pulumi/random`) — a unique repository-ID suffix

## Resources Created

- `data gcp_organizations_client_config` (`current`): Reads the active project from the provider's credentials.
- `random_random_string` (`unique-string`): A suffix giving the repository a unique ID.
- `gcp_artifactregistry_repository` (`repository`): Stores the application image.
- `docker-build_image` (`image`): Builds and pushes the image to Artifact Registry.
- `gcp_cloudrun_service` (`service`): Runs the container.
- `gcp_cloudrun_iam_member` (`invoker`): Allows public, unauthenticated access.

## Outputs

- **url**: The URL of the Cloud Run service.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Google Cloud CLI authenticated (`gcloud auth application-default login`).
- A Google Cloud project with the Artifact Registry and Cloud Run APIs enabled.
- A running Docker daemon (the image is built locally and pushed to Artifact Registry).

## Usage

```bash
pulumi new container-gcp-hcl
pulumi up
```

Open the `url` output once the service is deployed.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── app/                  # The container application (Dockerfile, source)
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., google:project)
```

## Configuration

- **google:project**: The Google Cloud project to deploy into.
- **region**: The region to deploy into. Default: `us-central1`.
- **app_path**: The container application folder. Default: `./app`.
- **image_name**: The image name. Default: `my-app`.
- **container_port**: The container port. Default: `8080`.
- **cpu** / **memory**: Per-instance CPU and memory. Defaults: `1` / `1Gi`.
- **concurrency**: Max concurrent requests per instance. Default: `80`.

## Next Steps

- Add a custom domain mapping to the service.
- Configure min/max instance scaling.
- Add a Cloud SQL or other backing service.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
