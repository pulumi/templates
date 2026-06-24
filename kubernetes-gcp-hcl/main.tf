terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.0.0"
    }
    random = {
      source  = "hashicorp/random"
      version = ">= 3.0.0"
    }
  }
}

variable "region" {
  description = "The Google Cloud region to deploy into"
  type        = string
  default     = "us-central1"
}

variable "nodes_per_zone" {
  description = "The number of nodes per zone in the node pool"
  type        = number
  default     = 1
}

data "google_client_config" "current" {}

locals {
  project = data.google_client_config.current.project
}

# A random suffix to make names unique.
resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

# Create a VPC network for the GKE cluster.
resource "google_compute_network" "network" {
  name                    = "gke-network-${random_string.suffix.result}"
  description             = "A virtual network for the GKE cluster"
  auto_create_subnetworks = false
}

# Create a subnet with Private Google Access enabled.
resource "google_compute_subnetwork" "subnet" {
  name                     = "gke-subnet-${random_string.suffix.result}"
  ip_cidr_range            = "10.128.0.0/12"
  region                   = var.region
  network                  = google_compute_network.network.id
  private_ip_google_access = true
}

# Create a GKE cluster with its default node pool removed.
resource "google_container_cluster" "cluster" {
  name     = "gke-cluster-${random_string.suffix.result}"
  location = var.region

  network    = google_compute_network.network.name
  subnetwork = google_compute_subnetwork.subnet.name

  networking_mode          = "VPC_NATIVE"
  remove_default_node_pool = true
  initial_node_count       = 1
  deletion_protection      = false

  ip_allocation_policy {}

  release_channel {
    channel = "STABLE"
  }

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "10.100.0.0/28"
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "0.0.0.0/0"
      display_name = "All networks"
    }
  }

  workload_identity_config {
    workload_pool = "${local.project}.svc.id.goog"
  }
}

# Create a service account for the node pool.
resource "google_service_account" "nodepool" {
  account_id   = "gke-np-${random_string.suffix.result}"
  display_name = "GKE node pool service account"
}

# Create a node pool for the cluster.
resource "google_container_node_pool" "nodepool" {
  name       = "nodepool-${random_string.suffix.result}"
  cluster    = google_container_cluster.cluster.id
  node_count = var.nodes_per_zone

  node_config {
    machine_type    = "e2-medium"
    service_account = google_service_account.nodepool.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}

# Export network and cluster details, plus a kubeconfig for the cluster.
output "network_name" {
  value = google_compute_network.network.name
}

output "cluster_name" {
  value = google_container_cluster.cluster.name
}

output "kubeconfig" {
  sensitive = true
  value     = <<-EOF
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: ${google_container_cluster.cluster.master_auth[0].cluster_ca_certificate}
        server: https://${google_container_cluster.cluster.endpoint}
      name: ${google_container_cluster.cluster.name}
    contexts:
    - context:
        cluster: ${google_container_cluster.cluster.name}
        user: ${google_container_cluster.cluster.name}
      name: ${google_container_cluster.cluster.name}
    current-context: ${google_container_cluster.cluster.name}
    kind: Config
    preferences: {}
    users:
    - name: ${google_container_cluster.cluster.name}
      user:
        exec:
          apiVersion: client.authentication.k8s.io/v1beta1
          command: gke-gcloud-auth-plugin
          provideClusterInfo: true
  EOF
}
