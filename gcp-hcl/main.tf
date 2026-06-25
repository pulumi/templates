terraform {
  required_providers {
    gcp = {
      source = "pulumi/gcp"
    }
  }
}

# Create a GCP resource (Storage Bucket)
resource "gcp_storage_bucket" "my-bucket" {
  location = "US"
}

# Export the URL of the bucket
output "bucket_name" {
  value = gcp_storage_bucket.my-bucket.url
}
