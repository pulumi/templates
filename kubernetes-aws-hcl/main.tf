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
  description = "Minimum size (number of nodes) of cluster"
  type        = number
  default     = 3
}

variable "max_cluster_size" {
  description = "Maximum size (number of nodes) of cluster"
  type        = number
  default     = 6
}

variable "desired_cluster_size" {
  description = "Desired number of nodes in the cluster"
  type        = number
  default     = 3
}

variable "node_instance_type" {
  description = "Instance type to use for worker nodes"
  type        = string
  default     = "t3.medium"
}

variable "vpc_network_cidr" {
  description = "Network CIDR to use for new VPC"
  type        = string
  default     = "10.0.0.0/16"
}

# Create a VPC for the EKS cluster
resource "awsx_ec2_vpc" "eks-vpc" {
  cidr_block           = var.vpc_network_cidr
  enable_dns_hostnames = true
}

# Create the EKS cluster
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

# Output the Kubeconfig for the cluster
output "kubeconfig" {
  value     = eks_cluster.eks-cluster.kubeconfig
  sensitive = true
}

output "vpc_id" {
  value = awsx_ec2_vpc.eks-vpc.vpc_id
}
