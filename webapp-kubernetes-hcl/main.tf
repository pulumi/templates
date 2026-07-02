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
  default     = "webapp"
}

variable "num_replicas" {
  description = "The number of replicas to deploy"
  type        = number
  default     = 1
}

locals {
  app_labels = {
    app = "nginx"
  }
}

# Create a new namespace for the resources
resource "kubernetes_core_v1_namespace" "webserverns" {
  metadata = {
    name = var.k8s_namespace
  }
}

# Create a ConfigMap to store Nginx configuration
resource "kubernetes_core_v1_config_map" "webserverconfig" {
  metadata = {
    namespace = kubernetes_core_v1_namespace.webserverns.metadata.name
  }

  data = {
    "nginx.conf" = <<-EOF
      events { }
      http {
        server {
          listen 80;
          root /usr/share/nginx/html;
          index index.html index.htm index.nginx-debian.html;
          server_name _;
          location / {
            try_files $uri $uri/ =404;
          }
        }
      }
    EOF
  }
}

# Create a new Deployment
resource "kubernetes_apps_v1_deployment" "webserverdeployment" {
  metadata = {
    namespace = kubernetes_core_v1_namespace.webserverns.metadata.name
  }

  spec = {
    replicas = var.num_replicas
    selector = {
      match_labels = local.app_labels
    }
    template = {
      metadata = {
        labels = local.app_labels
      }
      spec = {
        containers = [{
          name  = "nginx"
          image = "nginx"
          volume_mounts = [{
            name       = "nginx-conf-volume"
            mount_path = "/etc/nginx/nginx.conf"
            sub_path   = "nginx.conf"
            read_only  = true
          }]
        }]
        volumes = [{
          name = "nginx-conf-volume"
          config_map = {
            name = kubernetes_core_v1_config_map.webserverconfig.metadata.name
            items = [{
              key  = "nginx.conf"
              path = "nginx.conf"
            }]
          }
        }]
      }
    }
  }
}

# Expose the Deployment as a Kubernetes Service
resource "kubernetes_core_v1_service" "webserverservice" {
  metadata = {
    namespace = kubernetes_core_v1_namespace.webserverns.metadata.name
  }

  spec = {
    selector = local.app_labels
    ports = [{
      port        = 80
      target_port = 80
      protocol    = "TCP"
    }]
  }
}

# Export some values for use elsewhere
output "deployment_name" {
  value = kubernetes_apps_v1_deployment.webserverdeployment.metadata.name
}

output "service_name" {
  value = kubernetes_core_v1_service.webserverservice.metadata.name
}
