terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.0.0"
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
    echo '<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Hello, world!</title>
    </head>
    <body>
        <h1>Hello, world! 👋</h1>
        <p>Deployed with 💜 by <a href="https://pulumi.com/">Pulumi</a>.</p>
    </body>
    </html>' > index.html
    nohup python3 -m http.server ${var.service_port} &
  EOF
}

# Look up the latest Amazon Linux 2023 AMI.
data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-2023.*-x86_64"]
  }
}

# Create a VPC.
resource "aws_vpc" "vpc" {
  cidr_block           = var.vpc_network_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true
}

# Create an internet gateway.
resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.vpc.id
}

# Create a subnet that assigns instances a public IP address.
resource "aws_subnet" "subnet" {
  vpc_id                  = aws_vpc.vpc.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
}

# Create a route table that routes outbound traffic through the gateway.
resource "aws_route_table" "route_table" {
  vpc_id = aws_vpc.vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }
}

resource "aws_route_table_association" "association" {
  subnet_id      = aws_subnet.subnet.id
  route_table_id = aws_route_table.route_table.id
}

# Create a security group allowing inbound HTTP and all outbound traffic.
resource "aws_security_group" "sec_group" {
  description = "Enable HTTP access"
  vpc_id      = aws_vpc.vpc.id

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
resource "aws_instance" "server" {
  ami                    = data.aws_ami.amazon_linux.id
  instance_type          = var.instance_type
  subnet_id              = aws_subnet.subnet.id
  vpc_security_group_ids = [aws_security_group.sec_group.id]
  user_data              = local.user_data

  tags = {
    Name = "webserver"
  }
}

# Export the instance's public IP address, hostname, and URL.
output "ip" {
  value = aws_instance.server.public_ip
}

output "hostname" {
  value = aws_instance.server.public_dns
}

output "url" {
  value = "http://${aws_instance.server.public_dns}:${var.service_port}"
}
