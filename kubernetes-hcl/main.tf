terraform {
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
  }
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

locals {
  app_labels = {
    app = "nginx"
  }
}

# Create an nginx Deployment
resource "kubernetes_deployment_v1" "nginx" {
  metadata {
    name = "nginx"
  }

  spec {
    replicas = 1

    selector {
      match_labels = local.app_labels
    }

    template {
      metadata {
        labels = local.app_labels
      }

      spec {
        container {
          name  = "nginx"
          image = "nginx"
        }
      }
    }
  }
}

# Export the name of the Deployment
output "name" {
  value = kubernetes_deployment_v1.nginx.metadata[0].name
}
