terraform {
  required_providers {
    aws = {
      source = "pulumi/aws"
    }
  }
}

variable "instance_type" {
  description = "The Amazon EC2 instance type"
  type        = string
  default     = "t3.micro"
}

variable "vpc_network_cidr" {
  description = "The network CIDR to use for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "service_port" {
  description = "The HTTP service port to expose on the VM"
  type        = number
  default     = 80
}

locals {
  user_data = <<-EOF
    #!/bin/bash
    echo "Hello, World from Pulumi!" > index.html
    nohup python3 -m http.server ${var.service_port} &
  EOF
}

# Look up the latest Amazon Linux 2023 AMI.
data "aws_ec2_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filters {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }
}

# Create a VPC.
resource "aws_ec2_vpc" "vpc" {
  cidr_block           = var.vpc_network_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
}

# Create an internet gateway.
resource "aws_ec2_internet_gateway" "gateway" {
  vpc_id = aws_ec2_vpc.vpc.id
}

# Create a subnet that assigns instances a public IP address.
resource "aws_ec2_subnet" "subnet" {
  vpc_id                  = aws_ec2_vpc.vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

# Create a route table that routes outbound traffic through the gateway.
resource "aws_ec2_route_table" "route-table" {
  vpc_id = aws_ec2_vpc.vpc.id

  routes {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_ec2_internet_gateway.gateway.id
  }
}

resource "aws_ec2_route_table_association" "route-table-association" {
  subnet_id      = aws_ec2_subnet.subnet.id
  route_table_id = aws_ec2_route_table.route-table.id
}

# Create a security group allowing inbound HTTP and all outbound traffic.
resource "aws_ec2_security_group" "sec-group" {
  description = "Enable HTTP access"
  vpc_id      = aws_ec2_vpc.vpc.id

  ingress {
    from_port   = var.service_port
    to_port     = var.service_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Create and launch an EC2 instance into the public subnet.
resource "aws_ec2_instance" "server" {
  instance_type          = var.instance_type
  subnet_id              = aws_ec2_subnet.subnet.id
  vpc_security_group_ids = [aws_ec2_security_group.sec-group.id]
  user_data              = local.user_data
  ami                    = data.aws_ec2_ami.amazon_linux.id

  tags = {
    Name = "webserver"
  }
}

# Export the instance's public IP address, hostname, and URL.
output "ip" {
  value = aws_ec2_instance.server.public_ip
}

output "hostname" {
  value = aws_ec2_instance.server.public_dns
}

output "url" {
  value = "http://${aws_ec2_instance.server.public_dns}:${var.service_port}"
}
