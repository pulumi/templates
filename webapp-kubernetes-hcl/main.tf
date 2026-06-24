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

# Create a namespace for the resources.
resource "kubernetes_namespace_v1" "webserver" {
  metadata {
    name = var.k8s_namespace
  }
}

# Create a ConfigMap to store the Nginx configuration.
resource "kubernetes_config_map_v1" "config" {
  metadata {
    namespace     = kubernetes_namespace_v1.webserver.metadata[0].name
    generate_name = "webserver-config-"
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

# Create a Deployment running Nginx with the ConfigMap mounted.
resource "kubernetes_deployment_v1" "webserver" {
  metadata {
    namespace     = kubernetes_namespace_v1.webserver.metadata[0].name
    generate_name = "webserver-"
  }

  spec {
    replicas = var.num_replicas

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

          volume_mount {
            name       = "nginx-conf-volume"
            mount_path = "/etc/nginx/nginx.conf"
            sub_path   = "nginx.conf"
            read_only  = true
          }
        }

        volume {
          name = "nginx-conf-volume"

          config_map {
            name = kubernetes_config_map_v1.config.metadata[0].name
            items {
              key  = "nginx.conf"
              path = "nginx.conf"
            }
          }
        }
      }
    }
  }
}

# Expose the Deployment as a Service.
resource "kubernetes_service_v1" "webserver" {
  metadata {
    namespace     = kubernetes_namespace_v1.webserver.metadata[0].name
    generate_name = "webserver-"
  }

  spec {
    selector = local.app_labels

    port {
      port        = 80
      target_port = 80
      protocol    = "TCP"
    }
  }
}

# Export the deployment and service names.
output "deployment_name" {
  value = kubernetes_deployment_v1.webserver.metadata[0].name
}

output "service_name" {
  value = kubernetes_service_v1.webserver.metadata[0].name
}
