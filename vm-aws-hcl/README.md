# Virtual Machine on AWS (Pulumi HCL)

A Pulumi HCL program that deploys an EC2 virtual machine running a simple web server.

## Overview

The program creates a VPC with a public subnet, an internet gateway and routing, a security group allowing inbound HTTP, and an EC2 instance running the latest Amazon Linux 2023 AMI. A startup script serves a "Hello, world!" page. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`pulumi/aws`)

## Resources Created

- `data aws_ec2_ami` (`amazon_linux`): Looks up the latest Amazon Linux 2023 AMI.
- `aws_ec2_vpc` (`vpc`) / `aws_ec2_internet_gateway` (`gateway`) / `aws_ec2_subnet` (`subnet`) / `aws_ec2_route_table` (`route-table`) (+ `aws_ec2_route_table_association`): The network.
- `aws_ec2_security_group` (`sec-group`): Allows inbound HTTP and all outbound traffic.
- `aws_ec2_instance` (`server`): The EC2 instance running the web server.

## Outputs

- **ip**: The instance's public IP address.
- **hostname**: The instance's public DNS name.
- **url**: The HTTP URL of the web server.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- AWS credentials configured (environment variables, `~/.aws/credentials`, or `AWS_PROFILE`).

## Usage

```bash
pulumi new vm-aws-hcl
pulumi up
```

Open the `url` output once the instance has booted.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., aws:region)
```

## Configuration

- **aws:region**: The AWS region to deploy into. Default: `us-west-2`.
- **instance_type**: The EC2 instance type. Default: `t3.micro`.
- **vpc_network_cidr**: The VPC CIDR. Default: `10.0.0.0/16`.
- **service_port**: The HTTP port to serve on. Default: `80`.

## Next Steps

- Add an SSH key pair and a rule to allow SSH access.
- Replace the inline startup script with your own application.
- Put the instance behind a load balancer.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
