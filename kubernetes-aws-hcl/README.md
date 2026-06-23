# Kubernetes Cluster on AWS (Pulumi HCL)

A Pulumi HCL program that provisions an Amazon EKS cluster with a managed node group.

## Overview

The program creates a VPC with public subnets across two availability zones, the IAM roles required by EKS, an EKS control plane, and a managed node group. It exports a kubeconfig that uses the AWS CLI to authenticate. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`hashicorp/aws`)

## Resources Created

- `aws_vpc` / `aws_subnet` / `aws_internet_gateway` / `aws_route_table` (+ associations): The cluster network.
- `aws_iam_role` (`cluster`, `node`) + policy attachments: The control-plane and node roles.
- `aws_eks_cluster` (`cluster`): The EKS control plane.
- `aws_eks_node_group` (`nodes`): The managed worker node group.

## Outputs

- **clusterName**: The name of the EKS cluster.
- **vpcId**: The ID of the VPC.
- **kubeconfig**: A kubeconfig for the cluster (sensitive).

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
