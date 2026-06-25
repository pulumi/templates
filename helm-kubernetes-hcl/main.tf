terraform {
  required_providers {
    kubernetes = {
      source = "pulumi/kubernetes"
    }
  }
}

variable "k8s_namespace" {
  description = "The Kubernetes namespace to deploy into"
  type        = string
  default     = "nginx-ingress"
}

locals {
  app_labels = {
    app = "nginx-ingress"
  }

  # The Helm release name. Defined as a local because the Release's resource
  # token (kubernetes:helm.sh/v3:Release) contains a dot, which can't be
  # referenced by traversal in HCL — so we set and export the name directly.
  release_name = "ingresscontroller"
}

# Create a namespace for the ingress controller.
resource "kubernetes_core_v1_namespace" "ingressns" {
  metadata = {
    name   = var.k8s_namespace
    labels = local.app_labels
  }
}

# Install the NGINX ingress controller with the native Helm Release resource.
# (The resource token is kubernetes:helm.sh/v3:Release.)
resource "kubernetes_helm.sh_v3_release" "ingresscontroller" {
  name      = local.release_name
  chart     = "nginx-ingress"
  namespace = kubernetes_core_v1_namespace.ingressns.metadata.name
  skip_crds = true

  repository_opts = {
    repo = "https://helm.nginx.com/stable"
  }

  values = {
    controller = {
      enableCustomResources = false
      appprotect            = { enable = false }
      appprotectdos         = { enable = false }
      service               = { extraLabels = local.app_labels }
    }
  }

  version = "0.14.1"
}

# Export the name of the Helm release.
output "name" {
  value = local.release_name
}
