terraform {
  required_providers {
    cloudflare = {
      source = "pulumi/cloudflare"
    }
  }
}

variable "account_id" {
  description = "The Cloudflare account ID to deploy into"
  type        = string
}

# Create a Cloudflare resource (Workers KV namespace)
resource "cloudflare_workers_kv_namespace" "namespace" {
  account_id = var.account_id
  title      = "my-namespace"
}

# Export the ID of the namespace
output "namespace_id" {
  value = cloudflare_workers_kv_namespace.namespace.id
}
