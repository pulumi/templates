# Kubernetes Cluster on Google Cloud (Pulumi HCL)

A Pulumi HCL program that provisions a Google Kubernetes Engine (GKE) cluster with a custom node pool.

## Overview

The program creates a VPC-native GKE cluster with private nodes, removes the default node pool, and adds a managed node pool backed by a dedicated service account. It exports a kubeconfig that uses the `gke-gcloud-auth-plugin` to authenticate. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google (`hashicorp/google`)
- Random (`hashicorp/random`)

## Resources Created

- `google_compute_network` / `google_compute_subnetwork`: The cluster network (with Private Google Access).
- `google_container_cluster` (`cluster`): The GKE control plane (default node pool removed).
- `google_service_account` (`nodepool`): The node pool's service account.
- `google_container_node_pool` (`nodepool`): The managed node pool.

## Outputs

- **network_name**: The name of the VPC network.
- **cluster_name**: The name of the GKE cluster.
- **kubeconfig**: A kubeconfig for the cluster (sensitive). Requires `gke-gcloud-auth-plugin`.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Google Cloud CLI authenticated (`gcloud auth application-default login`) and the `gke-gcloud-auth-plugin` installed.
- A Google Cloud project with the Kubernetes Engine and Compute APIs enabled.

## Usage

```bash
pulumi new kubernetes-gcp-hcl
pulumi up
pulumi stack output kubeconfig --show-secrets > kubeconfig.yaml
KUBECONFIG=kubeconfig.yaml kubectl get nodes
```

A new GKE cluster typically takes several minutes to provision.

## Configuration

- **google:project**: The Google Cloud project to deploy into.
- **google:region**: The region to deploy into. Default: `us-central1`.
- **nodes_per_zone**: The number of nodes per zone. Default: `1`.

## Next Steps

- Enable Workload Identity bindings for your workloads.
- Add Cloud NAT if your private nodes need outbound internet access.
- Tighten the master authorized networks for the control-plane endpoint.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
