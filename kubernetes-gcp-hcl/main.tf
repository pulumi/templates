terraform {
  required_providers {
    gcp = {
      source = "pulumi/gcp"
    }
  }
}

variable "region" {
  description = "The Google Cloud region to deploy into"
  type        = string
  default     = "us-central1"
}

variable "nodes_per_zone" {
  description = "The desired number of nodes PER ZONE in the nodepool"
  type        = number
  default     = 1
}

# Read the active project from the provider's credentials.
data "gcp_organizations_client_config" "current" {}

locals {
  project = data.gcp_organizations_client_config.current.project

  # Build a kubeconfig for the cluster. It uses the gke-gcloud-auth-plugin to
  # authenticate, so that plugin must be installed to run kubectl against it.
  kubeconfig = <<-EOF
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: ${gcp_container_cluster.gke-cluster.master_auth.cluster_ca_certificate}
        server: https://${gcp_container_cluster.gke-cluster.endpoint}
      name: ${gcp_container_cluster.gke-cluster.name}
    contexts:
    - context:
        cluster: ${gcp_container_cluster.gke-cluster.name}
        user: ${gcp_container_cluster.gke-cluster.name}
      name: ${gcp_container_cluster.gke-cluster.name}
    current-context: ${gcp_container_cluster.gke-cluster.name}
    kind: Config
    preferences: {}
    users:
    - name: ${gcp_container_cluster.gke-cluster.name}
      user:
        exec:
          apiVersion: client.authentication.k8s.io/v1beta1
          command: gke-gcloud-auth-plugin
          provideClusterInfo: true
  EOF
}

# Create a GCP network (global VPC)
resource "gcp_compute_network" "gke-network" {
  description             = "A virtual network for the GKE cluster"
  auto_create_subnetworks = false
}

# Create a subnet in the new GCP network
resource "gcp_compute_subnetwork" "gke-subnet" {
  region                   = var.region
  ip_cidr_range            = "10.128.0.0/12"
  network                  = gcp_compute_network.gke-network.id
  private_ip_google_access = true
}

# Create a new GKE cluster
resource "gcp_container_cluster" "gke-cluster" {
  description              = "A GKE cluster"
  location                 = var.region
  network                  = gcp_compute_network.gke-network.name
  subnetwork               = gcp_compute_subnetwork.gke-subnet.name
  networking_mode          = "VPC_NATIVE"
  initial_node_count       = 1
  remove_default_node_pool = true
  datapath_provider        = "ADVANCED_DATAPATH"

  addons_config = {
    dns_cache_config = {
      enabled = true
    }
  }

  binary_authorization = {
    evaluation_mode = "PROJECT_SINGLETON_POLICY_ENFORCE"
  }

  ip_allocation_policy = {
    cluster_ipv4_cidr_block  = "/14"
    services_ipv4_cidr_block = "/20"
  }

  master_authorized_networks_config = {
    cidr_blocks = [{
      cidr_block   = "0.0.0.0/0"
      display_name = "All networks"
    }]
  }

  private_cluster_config = {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "10.100.0.0/28"
  }

  release_channel = {
    channel = "STABLE"
  }

  workload_identity_config = {
    workload_pool = "${local.project}.svc.id.goog"
  }
}

# Create a new service account for the nodepool
resource "gcp_serviceaccount_account" "gke-nodepool-sa" {
  account_id   = "${gcp_container_cluster.gke-cluster.name}-np-1-sa"
  display_name = "Nodepool 1 Service Account"
}

# Create a new nodepool for the cluster
resource "gcp_container_node_pool" "gke-nodepool" {
  cluster    = gcp_container_cluster.gke-cluster.id
  node_count = var.nodes_per_zone

  node_config = {
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]
    service_account = gcp_serviceaccount_account.gke-nodepool-sa.email
  }
}

# Export some values to be used elsewhere
output "network_name" {
  value = gcp_compute_network.gke-network.name
}

output "cluster_name" {
  value = gcp_container_cluster.gke-cluster.name
}

output "kubeconfig" {
  value     = local.kubeconfig
  sensitive = true
}
