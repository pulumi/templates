# Containerized Service on Google Cloud (Pulumi HCL)

A Pulumi HCL program that builds a container image and runs it on Google Cloud Run.

## Overview

The application in `./app` is built into a container image and pushed to Artifact Registry, then deployed as a publicly accessible Cloud Run service. The image is built and pushed with the Docker provider, so a running Docker daemon is required at deploy time. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google (`hashicorp/google`)
- Docker (`kreuzwerker/docker`) — builds and pushes the image
- Random (`hashicorp/random`)

## Resources Created

- `google_artifact_registry_repository` (`repo`): Stores the application image.
- `docker_image` / `docker_registry_image` (`app`): Builds and pushes the image.
- `google_cloud_run_v2_service` (`service`): Runs the container.
- `google_cloud_run_v2_service_iam_member` (`invoker`): Allows public access.

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
