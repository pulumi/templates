terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.0"
    }
  }
}

# The path to the container application to deploy
variable "app_path" {
  type    = string
  default = "./app"
}

# The name to give the container image
variable "image_name" {
  type    = string
  default = "my-app"
}

# The port to expose on the container
variable "container_port" {
  type    = number
  default = 80
}

# The amount of CPU units to allocate for the task (must be a valid Fargate value)
variable "cpu" {
  type    = number
  default = 256
}

# The amount of memory (MiB) to allocate for the task (must be a valid Fargate value)
variable "memory" {
  type    = number
  default = 512
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Use the default VPC and its subnets to keep the template self-contained.
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Credentials for pushing to ECR.
data "aws_ecr_authorization_token" "token" {}

locals {
  registry = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.region}.amazonaws.com"
}

# An ECR repository to store the application's container image.
resource "aws_ecr_repository" "repo" {
  name         = var.image_name
  force_delete = true
}

# Authenticate the Docker provider to ECR.
provider "docker" {
  registry_auth {
    address  = local.registry
    username = data.aws_ecr_authorization_token.token.user_name
    password = data.aws_ecr_authorization_token.token.password
  }
}

# Build the container image from the application source.
resource "docker_image" "app" {
  name = "${aws_ecr_repository.repo.repository_url}:latest"

  build {
    context  = var.app_path
    platform = "linux/amd64"
  }
}

# Push the image to ECR.
resource "docker_registry_image" "app" {
  name          = docker_image.app.name
  keep_remotely = true
}

# An ECS cluster to deploy into.
resource "aws_ecs_cluster" "cluster" {
  name = "${var.image_name}-cluster"
}

# A security group for the load balancer, allowing inbound HTTP.
resource "aws_security_group" "lb" {
  vpc_id = data.aws_vpc.default.id

  ingress {
    protocol    = "tcp"
    from_port   = 80
    to_port     = 80
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# A security group for the service, allowing traffic from the load balancer.
resource "aws_security_group" "service" {
  vpc_id = data.aws_vpc.default.id

  ingress {
    protocol        = "tcp"
    from_port       = var.container_port
    to_port         = var.container_port
    security_groups = [aws_security_group.lb.id]
  }

  egress {
    protocol    = "-1"
    from_port   = 0
    to_port     = 0
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# An Application Load Balancer to serve the container endpoint.
resource "aws_lb" "lb" {
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb.id]
  subnets            = data.aws_subnets.default.ids
}

resource "aws_lb_target_group" "tg" {
  port        = var.container_port
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = data.aws_vpc.default.id
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.lb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.tg.arn
  }
}

# An execution role for the ECS task.
resource "aws_iam_role" "execution" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "execution" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# A task definition describing the container to run on Fargate.
resource "aws_ecs_task_definition" "task" {
  family                   = "${var.image_name}-task"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = var.cpu
  memory                   = var.memory
  execution_role_arn       = aws_iam_role.execution.arn

  container_definitions = jsonencode([{
    name      = var.image_name
    image     = docker_registry_image.app.name
    essential = true
    portMappings = [{
      containerPort = var.container_port
      hostPort      = var.container_port
      protocol      = "tcp"
    }]
  }])
}

# A Fargate service that runs and exposes the container behind the load balancer.
resource "aws_ecs_service" "service" {
  name            = "${var.image_name}-service"
  cluster         = aws_ecs_cluster.cluster.arn
  task_definition = aws_ecs_task_definition.task.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = data.aws_subnets.default.ids
    security_groups  = [aws_security_group.service.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.tg.arn
    container_name   = var.image_name
    container_port   = var.container_port
  }

  depends_on = [aws_lb_listener.http]
}

# The URL at which the container's HTTP endpoint is available.
output "url" {
  value = "http://${aws_lb.lb.dns_name}"
}
