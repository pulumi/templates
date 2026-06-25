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
  # Define some labels that will be applied to resources
  app_labels = {
    app = "nginx-ingress"
  }

  # The Helm release name.
  release_name = "ingresscontroller"
}

# Create a namespace (name of the namespace supplied by the user)
resource "kubernetes_core_v1_namespace" "ingressns" {
  metadata = {
    name   = var.k8s_namespace
    labels = local.app_labels
  }
}

# Use Helm to install the Nginx ingress controller
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

# Export some values for use elsewhere
output "name" {
  value = local.release_name
}
