terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
}

variable "path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
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
  # Map file extensions to the content types used when uploading objects.
  mime_types = {
    ".html" = "text/html"
    ".css"  = "text/css"
    ".js"   = "application/javascript"
    ".json" = "application/json"
    ".svg"  = "image/svg+xml"
    ".png"  = "image/png"
    ".jpg"  = "image/jpeg"
    ".jpeg" = "image/jpeg"
    ".gif"  = "image/gif"
    ".ico"  = "image/x-icon"
    ".txt"  = "text/plain"
  }
}

# A random suffix to make the bucket name globally unique.
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a storage bucket and configure it as a website.
resource "google_storage_bucket" "bucket" {
  name     = "static-website-${random_string.suffix.result}"
  location = "US"

  website {
    main_page_suffix = var.index_document
    not_found_page   = var.error_document
  }
}

# Allow public read access to the bucket's objects.
resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# Sync the website files to the bucket.
resource "google_storage_bucket_object" "files" {
  for_each = fileset(var.path, "**")

  bucket       = google_storage_bucket.bucket.name
  name         = each.value
  source       = "${var.path}/${each.value}"
  content_type = lookup(local.mime_types, try(regex("\\.[^.]+$", each.value), ""), "application/octet-stream")
}

# Enable the storage bucket as a CDN backend.
resource "google_compute_backend_bucket" "backend" {
  name        = "static-website-${random_string.suffix.result}"
  bucket_name = google_storage_bucket.bucket.name
  enable_cdn  = true
}

# Provision a global IP address for the CDN.
resource "google_compute_global_address" "ip" {
  name = "static-website-${random_string.suffix.result}"
}

# Route requests to the storage bucket.
resource "google_compute_url_map" "url_map" {
  name            = "static-website-${random_string.suffix.result}"
  default_service = google_compute_backend_bucket.backend.self_link
}

# Create an HTTP proxy to route requests to the URL map.
resource "google_compute_target_http_proxy" "http_proxy" {
  name    = "static-website-${random_string.suffix.result}"
  url_map = google_compute_url_map.url_map.self_link
}

# Route incoming requests to the HTTP proxy.
resource "google_compute_global_forwarding_rule" "http" {
  name        = "static-website-${random_string.suffix.result}"
  ip_address  = google_compute_global_address.ip.address
  ip_protocol = "TCP"
  port_range  = "80"
  target      = google_compute_target_http_proxy.http_proxy.self_link
}

# Export the URLs and hostnames of the bucket and CDN.
output "origin_url" {
  value = "https://storage.googleapis.com/${google_storage_bucket.bucket.name}/${var.index_document}"
}

output "origin_hostname" {
  value = "storage.googleapis.com/${google_storage_bucket.bucket.name}"
}

output "cdn_url" {
  value = "http://${google_compute_global_address.ip.address}"
}

output "cdn_hostname" {
  value = google_compute_global_address.ip.address
}
