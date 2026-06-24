terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
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
}

# Create a namespace for the ingress controller.
resource "kubernetes_namespace_v1" "ingress" {
  metadata {
    name   = var.k8s_namespace
    labels = local.app_labels
  }
}

# Install the NGINX ingress controller via Helm.
resource "helm_release" "ingress" {
  name       = "ingresscontroller"
  namespace  = kubernetes_namespace_v1.ingress.metadata[0].name
  repository = "https://helm.nginx.com/stable"
  chart      = "nginx-ingress"
  version    = "0.14.1"
  skip_crds  = true

  set {
    name  = "controller.enableCustomResources"
    value = "false"
  }

  set {
    name  = "controller.appprotect.enable"
    value = "false"
  }

  set {
    name  = "controller.appprotectdos.enable"
    value = "false"
  }
}

# Export the name of the Helm release.
output "name" {
  value = helm_release.ingress.name
}
