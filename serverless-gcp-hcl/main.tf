terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0.0"
    }
    archive = {
      source  = "hashicorp/archive"
      version = ">= 2.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
}

variable "region" {
  description = "The Google Cloud region to deploy into"
  type        = string
  default     = "us-central1"
}

variable "site_path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
}

variable "app_path" {
  description = "The path to the folder containing the function to deploy"
  type        = string
  default     = "./app"
}

variable "index_document" {
  description = "The file to use for top-level pages"
  type        = string
  default     = "index.html"
}

variable "error_document" {
  description = "The file to use for error pages"
  type        = string
  default     = "error.html"
}

locals {
  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".json" = "application/json"
    ".svg"  = "image/svg+xml"
    ".ico"  = "image/x-icon"
    ".txt"  = "text/plain"
  }
}

# A random suffix to make resource names globally unique.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a storage bucket and configure it as a website.
resource "google_storage_bucket" "site" {
  name     = "serverless-site-${random_string.suffix.result}"
  location = "US"

  website {
    main_page_suffix = var.index_document
    not_found_page   = var.error_document
  }
}

# Allow public read access to the website's objects.
resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.site.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# Sync the website files to the bucket.
resource "google_storage_bucket_object" "files" {
  for_each = fileset(var.site_path, "**")

  bucket       = google_storage_bucket.site.name
  name         = each.value
  source       = "${var.site_path}/${each.value}"
  content_type = lookup(local.mime_types, try(regex("\\.[^.]+$", each.value), ""), "application/octet-stream")
}

# Create a bucket to hold the serverless app's source archive.
resource "google_storage_bucket" "app" {
  name     = "serverless-app-${random_string.suffix.result}"
  location = "US"
}

# Package the function source into a deployment archive.
data "archive_file" "app" {
  type        = "zip"
  source_dir  = var.app_path
  output_path = "${path.module}/app.zip"
}

# Upload the serverless app to the bucket.
resource "google_storage_bucket_object" "app" {
  bucket = google_storage_bucket.app.name
  name   = "app-${data.archive_file.app.output_md5}.zip"
  source = data.archive_file.app.output_path
}

# Create a Cloud Function (Gen 2) that returns the current time.
resource "google_cloudfunctions2_function" "data" {
  name     = "serverless-fn-${random_string.suffix.result}"
  location = var.region

  build_config {
    runtime     = "python312"
    entry_point = "data"
    source {
      storage_source {
        bucket = google_storage_bucket.app.name
        object = google_storage_bucket_object.app.name
      }
    }
  }

  service_config {
    available_memory = "256M"
    timeout_seconds  = 60
  }
}

# Allow public, unauthenticated invocations of the underlying Cloud Run service.
resource "google_cloud_run_v2_service_iam_member" "invoker" {
  location = google_cloudfunctions2_function.data.location
  name     = google_cloudfunctions2_function.data.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Write a config file the website uses to find the function endpoint.
resource "google_storage_bucket_object" "config" {
  bucket       = google_storage_bucket.site.name
  name         = "config.json"
  content      = jsonencode({ api = google_cloudfunctions2_function.data.url })
  content_type = "application/json"
}

# The URL of the website.
output "site_url" {
  value = "https://storage.googleapis.com/${google_storage_bucket.site.name}/${var.index_document}"
}

# The URL of the serverless endpoint.
output "api_url" {
  value = google_cloudfunctions2_function.data.url
}
