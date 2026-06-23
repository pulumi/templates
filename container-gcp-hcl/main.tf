terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0.0"
    }
    docker = {
      source  = "kreuzwerker/docker"
      version = ">= 3.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
}

# The Google Cloud region to deploy into
variable "region" {
  type    = string
  default = "us-central1"
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
  default = 8080
}

# The number of vCPUs to allocate per container instance
variable "cpu" {
  type    = number
  default = 1
}

# The amount of memory to allocate per container instance
variable "memory" {
  type    = string
  default = "1Gi"
}

# The maximum number of concurrent requests per container instance
variable "concurrency" {
  type    = number
  default = 80
}

# Read the active project from the provider's credentials.
data "google_client_config" "current" {}

locals {
  project  = data.google_client_config.current.project
  repo_url = "${var.region}-docker.pkg.dev/${local.project}/${google_artifact_registry_repository.repo.repository_id}"
}

# A random suffix to give the repository a unique ID.
resource "random_string" "suffix" {
  length  = 4
  special = false
  upper   = false
}

# Create an Artifact Registry repository for the container image.
resource "google_artifact_registry_repository" "repo" {
  location      = var.region
  repository_id = "repo-${random_string.suffix.result}"
  description   = "Repository for the container image"
  format        = "DOCKER"
}

# Authenticate the Docker provider to Artifact Registry using a short-lived token.
provider "docker" {
  registry_auth {
    address  = "${var.region}-docker.pkg.dev"
    username = "oauth2accesstoken"
    password = data.google_client_config.current.access_token
  }
}

# Build the container image from the application source.
resource "docker_image" "app" {
  name = "${local.repo_url}/${var.image_name}:latest"

  build {
    context  = var.app_path
    platform = "linux/amd64"
  }
}

# Push the image to Artifact Registry.
resource "docker_registry_image" "app" {
  name          = docker_image.app.name
  keep_remotely = true
}

# Deploy the image as a Cloud Run service.
resource "google_cloud_run_v2_service" "service" {
  name                = "service-${random_string.suffix.result}"
  location            = var.region
  deletion_protection = false

  template {
    containers {
      image = docker_registry_image.app.name

      ports {
        container_port = var.container_port
      }

      resources {
        limits = {
          cpu    = tostring(var.cpu)
          memory = var.memory
        }
      }
    }

    max_instance_request_concurrency = var.concurrency
  }
}

# Allow public, unauthenticated access to the service.
resource "google_cloud_run_v2_service_iam_member" "invoker" {
  location = google_cloud_run_v2_service.service.location
  name     = google_cloud_run_v2_service.service.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Export the URL of the service.
output "url" {
  value = google_cloud_run_v2_service.service.uri
}
