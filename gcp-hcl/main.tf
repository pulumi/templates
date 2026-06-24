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

# A random suffix to make the bucket name globally unique
resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create a GCP resource (Storage Bucket)
resource "google_storage_bucket" "my_bucket" {
  name     = "my-bucket-${random_string.suffix.result}"
  location = "US"
}

# Export the URL of the bucket
output "bucket_name" {
  value = google_storage_bucket.my_bucket.url
}
