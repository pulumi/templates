# Kubernetes Cluster on Azure (Pulumi HCL)

A Pulumi HCL program that provisions an Azure Kubernetes Service (AKS) cluster.

## Overview

The program creates a resource group and a managed AKS cluster with a system-assigned identity and a default node pool. It exports a kubeconfig for the cluster. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Azure Native (`pulumi/azure-native`)

## Resources Created

- `azure-native_resources_resource_group` (`resource-group`): The resource group.
- `azure-native_containerservice_managed_cluster` (`cluster`): The AKS cluster and its system node pool.
- `data azure-native_containerservice_list_managed_cluster_user_credentials` (`credentials`): Fetches the cluster's user credentials for the exported kubeconfig.

## Outputs

- **cluster_name**: The name of the AKS cluster.
- **kubeconfig**: A kubeconfig for the cluster (sensitive).

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Azure CLI installed and signed in (`az login`). Set `ARM_SUBSCRIPTION_ID` to choose a specific subscription.
- An Azure subscription with permissions to create AKS clusters.

## Usage

```bash
pulumi new kubernetes-azure-hcl
pulumi up
pulumi stack output kubeconfig --show-secrets > kubeconfig.yaml
KUBECONFIG=kubeconfig.yaml kubectl get nodes
```

## Configuration

- **location**: The Azure region to deploy into. Default: `westus2`.
- **node_count**: The number of worker nodes. Default: `3`.
- **dns_prefix**: The cluster DNS prefix. Default: `pulumi`.
- **node_vm_size**: The worker node VM size. Default: `Standard_DS2_v2`.

## Next Steps

- Add Microsoft Entra ID (Azure AD) integration and Azure RBAC.
- Add additional node pools for different workloads.
- Enable monitoring with Azure Monitor / Container Insights.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
