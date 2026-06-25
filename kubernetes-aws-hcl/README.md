# Kubernetes Cluster on AWS (Pulumi HCL)

A Pulumi HCL program that provisions an Amazon EKS cluster with a managed node group.

## Overview

The `awsx` VPC component builds a VPC with public and private subnets, and the `eks` component provisions the EKS control plane and a managed node group. It exports a kubeconfig for the cluster. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`pulumi/aws`)
- AWSx (`pulumi/awsx`) — the VPC component
- EKS (`pulumi/eks`) — the EKS cluster component
- Kubernetes (`pulumi/kubernetes`)

## Resources Created

- `awsx_ec2_vpc` (`eks-vpc`): A VPC with public and private subnets for the cluster.
- `eks_cluster` (`eks-cluster`): The EKS control plane and managed node group.

## Outputs

- **kubeconfig**: A kubeconfig for the cluster (sensitive).
- **vpc_id**: The ID of the VPC.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- AWS credentials configured, and the AWS CLI installed (used by the kubeconfig to obtain tokens).
- An AWS account with permissions for EKS, EC2, and IAM.

## Usage

```bash
pulumi new kubernetes-aws-hcl
pulumi up
pulumi stack output kubeconfig --show-secrets > kubeconfig.yaml
KUBECONFIG=kubeconfig.yaml kubectl get nodes
```

A new EKS cluster typically takes 10–15 minutes to provision.

## Configuration

- **aws:region**: The AWS region to deploy into. Default: `us-west-2`.
- **min_cluster_size** / **max_cluster_size** / **desired_cluster_size**: Node group sizing. Defaults: `3` / `6` / `3`.
- **node_instance_type**: Worker node instance type. Default: `t3.medium`.
- **vpc_network_cidr**: The VPC CIDR. Default: `10.0.0.0/16`.

## Next Steps

- Add the EKS VPC CNI, CoreDNS, and kube-proxy add-ons explicitly.
- Deploy a workload with the Kubernetes provider, using this cluster's kubeconfig.
- Restrict the cluster endpoint to known networks.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
