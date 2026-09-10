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

variable "path" {
  description = "The path to the folder containing the website"
  type        = string
  default     = "./www"
}

# A Cloudflare Worker, exposed on its workers.dev subdomain.
resource "cloudflare_worker" "site" {
  account_id = var.account_id
  name       = "my-site"
  subdomain = {
    enabled = true
  }
}

# A version of the Worker that serves the website's files as static assets.
resource "cloudflare_worker_version" "version" {
  account_id         = var.account_id
  worker_id          = cloudflare_worker.site.id
  compatibility_date = "2025-01-01"

  assets = {
    directory = var.path
    config = {
      # Serve /404.html when a request doesn't match a file.
      not_found_handling = "404-page"
      html_handling      = "auto-trailing-slash"
    }
  }
}

# Deploy the version so it serves all of the site's traffic.
resource "cloudflare_workers_deployment" "deployment" {
  account_id  = var.account_id
  script_name = cloudflare_worker.site.name
  strategy    = "percentage"

  versions {
    version_id = cloudflare_worker_version.version.id
    percentage = 100
  }
}

# Export the website's URL.
output "url" {
  value = cloudflare_worker.site.subdomain.url
}
