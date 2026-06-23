# Containerized Service on AWS (Pulumi HCL)

A Pulumi HCL program that builds a container image and runs it on AWS ECS Fargate behind an Application Load Balancer.

## Overview

The application in `./app` is built into a container image and pushed to Amazon ECR, then deployed as an ECS Fargate service fronted by an Application Load Balancer. The image is built and pushed with the Docker provider, so a running Docker daemon is required at deploy time. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`hashicorp/aws`)
- Docker (`kreuzwerker/docker`) — builds and pushes the image

## Resources Created

- `aws_ecr_repository` (`repo`): Stores the application image.
- `docker_image` / `docker_registry_image` (`app`): Builds and pushes the image to ECR.
- `aws_ecs_cluster` (`cluster`): The cluster to run the service in.
- `aws_security_group` (`lb`, `service`): Allow HTTP to the load balancer and traffic from it to the task.
- `aws_lb` / `aws_lb_target_group` / `aws_lb_listener`: The Application Load Balancer.
- `aws_iam_role` (`execution`) + attachment: The ECS task execution role.
- `aws_ecs_task_definition` (`task`) / `aws_ecs_service` (`service`): The Fargate task and service.

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
- **app_path**: The container application folder. Default: `./app`.
- **image_name**: The image (and resource) name. Default: `my-app`.
- **container_port**: The container port. Default: `80`.
- **cpu** / **memory**: Task CPU units and memory (MiB) — must be a valid Fargate combination. Defaults: `256` / `512`.

## Next Steps

- Add an HTTPS listener and an ACM certificate.
- Add autoscaling to the ECS service.
- Add CloudWatch logging to the task definition.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
