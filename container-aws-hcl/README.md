# Containerized Service on AWS (Pulumi HCL)

A Pulumi HCL program that builds a container image and runs it on AWS ECS Fargate behind an Application Load Balancer.

## Overview

The application in `./app` is built into a container image and pushed to Amazon ECR, then deployed as an ECS Fargate service fronted by an Application Load Balancer. The `awsx` component builds and pushes the image, so a running Docker daemon is required at deploy time. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`pulumi/aws`)
- AWSx (`pulumi/awsx`) — load balancer, ECR repository/image, and Fargate service components

## Resources Created

- `aws_ecs_cluster` (`cluster`): The cluster to run the service in.
- `awsx_lb_application_load_balancer` (`loadbalancer`): The Application Load Balancer serving the container endpoint.
- `awsx_ecr_repository` (`repo`): Stores the application image.
- `awsx_ecr_image` (`image`): Builds and pushes the image from `./app` to ECR.
- `awsx_ecs_fargate_service` (`service`): The Fargate service hosting the application container.

## Outputs

- **url**: The HTTP URL of the load balancer.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- AWS credentials configured (environment variables, `~/.aws/credentials`, or `AWS_PROFILE`).
- A running Docker daemon (the image is built locally and pushed to ECR).

## Usage

```bash
pulumi new container-aws-hcl
pulumi up
```

Open the `url` output once the service is healthy.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── app/                  # The container application (Dockerfile)
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., aws:region)
```

## Configuration

- **aws:region**: The AWS region to deploy into. Default: `us-west-2`.
- **container_port**: The container port. Default: `80`.
- **cpu** / **memory**: Task CPU units and memory (MiB). Defaults: `512` / `128`.

The container image is built from `./app`, set in `main.tf`.

## Next Steps

- Add an HTTPS listener and an ACM certificate.
- Add autoscaling to the ECS service.
- Add CloudWatch logging to the task definition.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
