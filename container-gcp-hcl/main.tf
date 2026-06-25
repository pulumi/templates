terraform {
  required_providers {
    gcp = {
      source = "pulumi/gcp"
    }
    docker-build = {
      source = "pulumi/docker-build"
    }
    random = {
      source = "pulumi/random"
    }
  }
}

variable "region" {
  description = "The Google Cloud region to deploy into"
  type        = string
  default     = "us-central1"
}

variable "image_name" {
  description = "The name to give the container image"
  type        = string
  default     = "my-app"
}

variable "app_path" {
  description = "The path to the container application to deploy"
  type        = string
  default     = "./app"
}

variable "container_port" {
  description = "The port to expose on the container"
  type        = number
  default     = 8080
}

variable "cpu" {
  description = "The number of vCPUs to allocate per container instance"
  type        = number
  default     = 1
}

variable "memory" {
  description = "The amount of memory to allocate per container instance"
  type        = string
  default     = "1Gi"
}

variable "concurrency" {
  description = "The maximum number of concurrent requests per container instance"
  type        = number
  default     = 80
}

# Read the active project from the provider's credentials.
data "gcp_organizations_client_config" "current" {}

locals {
  repo_url = "${var.region}-docker.pkg.dev/${data.gcp_organizations_client_config.current.project}/${gcp_artifactregistry_repository.repository.repository_id}"
}

# A random suffix to give the repository a unique ID.
resource "random_random_string" "unique-string" {
  length  = 4
  lower   = true
  upper   = false
  numeric = true
  special = false
}

# Create an Artifact Registry repository for the container image.
resource "gcp_artifactregistry_repository" "repository" {
  description   = "Repository for the container image"
  format        = "DOCKER"
  location      = var.region
  repository_id = "repo-${random_random_string.unique-string.result}"
}

# Build the container image and push it to Artifact Registry.
# Before running `pulumi up`, configure Docker auth for Artifact Registry, e.g.:
#   gcloud auth configure-docker ${var.region}-docker.pkg.dev
resource "docker-build_image" "image" {
  tags      = ["${local.repo_url}/${var.image_name}"]
  platforms = ["linux/amd64"]
  push      = true

  context = {
    location = var.app_path
  }
}

# Deploy the image as a Cloud Run service.
resource "gcp_cloudrun_service" "service" {
  location = var.region

  template = {
    spec = {
      container_concurrency = var.concurrency
      containers = [{
        image = docker-build_image.image.ref
        ports = [{
          container_port = var.container_port
        }]
        resources = {
          limits = {
            cpu    = tostring(var.cpu)
            memory = var.memory
          }
        }
      }]
    }
  }
}

# Allow public, unauthenticated access to the service.
resource "gcp_cloudrun_iam_member" "invoker" {
  location = var.region
  service  = gcp_cloudrun_service.service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Export the URL of the service.
output "url" {
  value = gcp_cloudrun_service.service.statuses[0].url
}
