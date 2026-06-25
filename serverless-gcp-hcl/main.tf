terraform {
  required_providers {
    gcp = {
      source = "pulumi/gcp"
    }
    synced-folder = {
      source = "pulumi/synced-folder"
    }
  }
}

variable "region" {
  description = "The Google Cloud region to deploy the function into"
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

# Create a storage bucket and configure it as a website.
resource "gcp_storage_bucket" "site-bucket" {
  location = "US"

  website = {
    main_page_suffix = var.index_document
    not_found_page   = var.error_document
  }
}

# Allow public read access to the website's objects.
resource "gcp_storage_bucket_i_a_m_binding" "site-bucket-iam-binding" {
  bucket  = gcp_storage_bucket.site-bucket.name
  role    = "roles/storage.objectViewer"
  members = ["allUsers"]
}

# Sync the contents of the website folder to the bucket.
resource "synced-folder_google_cloud_folder" "synced-folder" {
  path        = var.site_path
  bucket_name = gcp_storage_bucket.site-bucket.name
}

# Create a bucket to hold the function's source archive.
resource "gcp_storage_bucket" "app-bucket" {
  location = "US"
}

# Upload the zipped function source to the bucket.
resource "gcp_storage_bucket_object" "app-archive" {
  bucket = gcp_storage_bucket.app-bucket.name
  source = fileArchive(var.app_path)
}

# Create a Cloud Function (Gen 2) that returns the current time.
resource "gcp_cloudfunctionsv2_function" "data-function" {
  location = var.region

  build_config = {
    runtime     = "python312"
    entry_point = "data"
    source = {
      storage_source = {
        bucket = gcp_storage_bucket.app-bucket.name
        object = gcp_storage_bucket_object.app-archive.name
      }
    }
  }

  service_config = {
    available_memory = "256M"
    timeout_seconds  = 60
  }
}

# Allow public, unauthenticated invocations of the underlying Cloud Run service.
resource "gcp_cloudrun_iam_member" "invoker" {
  location = gcp_cloudfunctionsv2_function.data-function.location
  service  = gcp_cloudfunctionsv2_function.data-function.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Write a config file the website uses to find the function endpoint.
resource "gcp_storage_bucket_object" "site-config" {
  name         = "config.json"
  bucket       = gcp_storage_bucket.site-bucket.name
  content_type = "application/json"
  source       = stringAsset(jsonencode({ api = gcp_cloudfunctionsv2_function.data-function.url }))
}

# Export the URLs of the website and serverless endpoint.
output "site_url" {
  value = "https://storage.googleapis.com/${gcp_storage_bucket.site-bucket.name}/${var.index_document}"
}

output "api_url" {
  value = gcp_cloudfunctionsv2_function.data-function.url
}
