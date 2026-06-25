terraform {
  required_providers {
    aws = {
      source = "pulumi/aws"
    }
    awsx = {
      source = "pulumi/awsx"
    }
  }
}

variable "container_port" {
  description = "The port to expose on the container"
  type        = number
  default     = 80
}

variable "cpu" {
  description = "The amount of CPU to allocate for the container"
  type        = number
  default     = 512
}

variable "memory" {
  description = "The amount of memory to allocate for the container"
  type        = number
  default     = 128
}

# An ECS cluster to deploy into.
resource "aws_ecs_cluster" "cluster" {}

# An Application Load Balancer to serve the container endpoint to the internet.
resource "awsx_lb_application_load_balancer" "loadbalancer" {}

# An ECR repository to store the application's container image.
resource "awsx_ecr_repository" "repo" {
  force_delete = true
}

# Build and publish the image from ./app to the ECR repository.
resource "awsx_ecr_image" "image" {
  repository_url = awsx_ecr_repository.repo.url
  context        = "./app"
  platform       = "linux/amd64"
}

# Deploy an ECS Service on Fargate to host the application container.
resource "awsx_ecs_fargate_service" "service" {
  cluster          = aws_ecs_cluster.cluster.arn
  assign_public_ip = true

  task_definition_args = {
    container = {
      name      = "app"
      image     = awsx_ecr_image.image.image_uri
      cpu       = var.cpu
      memory    = var.memory
      essential = true
      port_mappings = [{
        container_port = var.container_port
        target_group   = awsx_lb_application_load_balancer.loadbalancer.default_target_group
      }]
    }
  }
}

# The URL at which the container's HTTP endpoint is available.
output "url" {
  value = "http://${awsx_lb_application_load_balancer.loadbalancer.load_balancer.dns_name}"
}
