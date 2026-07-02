terraform {
  required_providers {
    kubernetes = {
      source = "pulumi/kubernetes"
    }
  }
}

locals {
  app_labels = {
    app = "nginx"
  }
}

resource "kubernetes_apps_v1_deployment" "deployment" {
  spec = {
    selector = {
      match_labels = local.app_labels
    }
    replicas = 1
    template = {
      metadata = {
        labels = local.app_labels
      }
      spec = {
        containers = [{
          name  = "nginx"
          image = "nginx"
        }]
      }
    }
  }
}

output "name" {
  value = kubernetes_apps_v1_deployment.deployment.metadata.name
}
