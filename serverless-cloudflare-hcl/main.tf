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

# A Workers KV namespace to persist state between requests.
resource "cloudflare_workers_kv_namespace" "namespace" {
  account_id = var.account_id
  title      = "visit-counter"
}

# A Cloudflare Worker, exposed on its workers.dev subdomain.
resource "cloudflare_worker" "worker" {
  account_id = var.account_id
  name       = "my-app"
  subdomain = {
    enabled = true
  }
}

# A version of the Worker containing the code and its KV binding.
resource "cloudflare_worker_version" "version" {
  account_id         = var.account_id
  worker_id          = cloudflare_worker.worker.id
  compatibility_date = "2025-01-01"
  main_module        = "worker.js"

  bindings {
    name         = "COUNTER"
    type         = "kv_namespace"
    namespace_id = cloudflare_workers_kv_namespace.namespace.id
  }

  modules {
    name         = "worker.js"
    content_type = "application/javascript+module"
    content_file = "worker.js"
  }
}

# Deploy the version so it serves all of the application's traffic.
resource "cloudflare_workers_deployment" "deployment" {
  account_id  = var.account_id
  script_name = cloudflare_worker.worker.name
  strategy    = "percentage"

  versions {
    version_id = cloudflare_worker_version.version.id
    percentage = 100
  }
}

# Export the application's URL.
output "url" {
  value = cloudflare_worker.worker.subdomain.url
}
