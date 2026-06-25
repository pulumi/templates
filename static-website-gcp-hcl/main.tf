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

# Create a storage bucket and configure it as a website.
resource "gcp_storage_bucket" "bucket" {
  location = "US"

  website = {
    main_page_suffix = var.index_document
    not_found_page   = var.error_document
  }
}

# Create an IAM binding to allow public read access to the bucket.
resource "gcp_storage_bucket_i_a_m_binding" "bucket-iam-binding" {
  bucket  = gcp_storage_bucket.bucket.name
  role    = "roles/storage.objectViewer"
  members = ["allUsers"]
}

# Use a synced folder to manage the files of the website.
resource "synced-folder_google_cloud_folder" "synced-folder" {
  path        = var.path
  bucket_name = gcp_storage_bucket.bucket.name
}

# Enable the storage bucket as a CDN.
resource "gcp_compute_backend_bucket" "backend-bucket" {
  bucket_name = gcp_storage_bucket.bucket.name
  enable_cdn  = true
}

# Provision a global IP address for the CDN.
resource "gcp_compute_global_address" "ip" {}

# Create a URLMap to route requests to the storage bucket.
resource "gcp_compute_u_r_l_map" "url-map" {
  default_service = gcp_compute_backend_bucket.backend-bucket.self_link
}

# Create an HTTP proxy to route requests to the URLMap.
resource "gcp_compute_target_http_proxy" "http-proxy" {
  url_map = gcp_compute_u_r_l_map.url-map.self_link
}

# Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
resource "gcp_compute_global_forwarding_rule" "http-forwarding-rule" {
  ip_address  = gcp_compute_global_address.ip.address
  ip_protocol = "TCP"
  port_range  = "80"
  target      = gcp_compute_target_http_proxy.http-proxy.self_link
}

# Export the URLs and hostnames of the bucket and CDN.
output "origin_url" {
  value = "https://storage.googleapis.com/${gcp_storage_bucket.bucket.name}/${var.index_document}"
}

output "origin_hostname" {
  value = "storage.googleapis.com/${gcp_storage_bucket.bucket.name}"
}

output "cdn_url" {
  value = "http://${gcp_compute_global_address.ip.address}"
}

output "cdn_hostname" {
  value = gcp_compute_global_address.ip.address
}
