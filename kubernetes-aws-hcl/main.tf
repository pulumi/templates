terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
    }
  }
}

# The minimum number of nodes in the cluster
variable "min_cluster_size" {
  type    = number
  default = 3
}

# The maximum number of nodes in the cluster
variable "max_cluster_size" {
  type    = number
  default = 6
}

# The desired number of nodes in the cluster
variable "desired_cluster_size" {
  type    = number
  default = 3
}

# The EC2 instance type to use for worker nodes
variable "node_instance_type" {
  type    = string
  default = "t3.medium"
}

# The network CIDR to use for the VPC
variable "vpc_network_cidr" {
  type    = string
  default = "10.0.0.0/16"
}

data "aws_region" "current" {}

data "aws_availability_zones" "available" {
  state = "available"
}

# Create a VPC for the cluster.
resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_network_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.vpc.id
}

# Create two public subnets in two availability zones.
resource "aws_subnet" "public" {
  count                   = 2
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = cidrsubnet(var.vpc_network_cidr, 8, count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  map_public_ip_on_launch = true

  tags = {
    "kubernetes.io/role/elb" = "1"
  }
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }
}

resource "aws_route_table_association" "public" {
  count          = 2
  subnet_id      = aws_subnet.public[count.index].id
  route_table_id = aws_route_table.public.id
}

# An IAM role for the EKS control plane.
resource "aws_iam_role" "cluster" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "eks.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "cluster" {
  role       = aws_iam_role.cluster.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
}

# An IAM role for the worker nodes.
resource "aws_iam_role" "node" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = { Service = "ec2.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "node" {
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
    "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
  ])
  role       = aws_iam_role.node.name
  policy_arn = each.value
}

# Create the EKS cluster.
resource "aws_eks_cluster" "cluster" {
  name     = "eks-cluster"
  role_arn = aws_iam_role.cluster.arn

  vpc_config {
    subnet_ids              = aws_subnet.public[*].id
    endpoint_public_access  = true
    endpoint_private_access = true
  }

  depends_on = [aws_iam_role_policy_attachment.cluster]
}

# Create a managed node group for the cluster.
resource "aws_eks_node_group" "nodes" {
  cluster_name    = aws_eks_cluster.cluster.name
  node_group_name = "default"
  node_role_arn   = aws_iam_role.node.arn
  subnet_ids      = aws_subnet.public[*].id
  instance_types  = [var.node_instance_type]

  scaling_config {
    desired_size = var.desired_cluster_size
    min_size     = var.min_cluster_size
    max_size     = var.max_cluster_size
  }

  depends_on = [aws_iam_role_policy_attachment.node]
}

# Export the cluster name, VPC ID, and a kubeconfig for the cluster.
output "clusterName" {
  value = aws_eks_cluster.cluster.name
}

output "vpcId" {
  value = aws_vpc.vpc.id
}

output "kubeconfig" {
  sensitive = true
  value     = <<-EOF
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: ${aws_eks_cluster.cluster.certificate_authority[0].data}
        server: ${aws_eks_cluster.cluster.endpoint}
      name: ${aws_eks_cluster.cluster.name}
    contexts:
    - context:
        cluster: ${aws_eks_cluster.cluster.name}
        user: ${aws_eks_cluster.cluster.name}
      name: ${aws_eks_cluster.cluster.name}
    current-context: ${aws_eks_cluster.cluster.name}
    kind: Config
    preferences: {}
    users:
    - name: ${aws_eks_cluster.cluster.name}
      user:
        exec:
          apiVersion: client.authentication.k8s.io/v1beta1
          command: aws
          args:
          - eks
          - get-token
          - --cluster-name
          - ${aws_eks_cluster.cluster.name}
          - --region
          - ${data.aws_region.current.region}
  EOF
}
