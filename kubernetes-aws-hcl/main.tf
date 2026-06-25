terraform {
  required_providers {
    aws = {
      source = "pulumi/aws"
    }
    awsx = {
      source = "pulumi/awsx"
    }
    eks = {
      source = "pulumi/eks"
    }
    kubernetes = {
      source = "pulumi/kubernetes"
    }
  }
}

variable "min_cluster_size" {
  description = "The minimum number of nodes in the cluster"
  type        = number
  default     = 3
}

variable "max_cluster_size" {
  description = "The maximum number of nodes in the cluster"
  type        = number
  default     = 6
}

variable "desired_cluster_size" {
  description = "The desired number of nodes in the cluster"
  type        = number
  default     = 3
}

variable "node_instance_type" {
  description = "The EC2 instance type to use for worker nodes"
  type        = string
  default     = "t3.medium"
}

variable "vpc_network_cidr" {
  description = "The network CIDR to use for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

# Create a VPC for the cluster (the awsx component builds public and private subnets).
resource "awsx_ec2_vpc" "eks-vpc" {
  cidr_block           = var.vpc_network_cidr
  enable_dns_hostnames = true
}

# Create the EKS cluster (the eks component provisions the control plane and node group).
resource "eks_cluster" "eks-cluster" {
  vpc_id              = awsx_ec2_vpc.eks-vpc.vpc_id
  public_subnet_ids   = awsx_ec2_vpc.eks-vpc.public_subnet_ids
  private_subnet_ids  = awsx_ec2_vpc.eks-vpc.private_subnet_ids
  authentication_mode = "API"

  instance_type    = var.node_instance_type
  desired_capacity = var.desired_cluster_size
  min_size         = var.min_cluster_size
  max_size         = var.max_cluster_size

  node_associate_public_ip_address = false
  endpoint_private_access          = false
  endpoint_public_access           = true
}

# Export the cluster's kubeconfig and the VPC ID.
output "kubeconfig" {
  value     = eks_cluster.eks-cluster.kubeconfig
  sensitive = true
}

output "vpc_id" {
  value = awsx_ec2_vpc.eks-vpc.vpc_id
}
