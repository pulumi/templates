import pulumi
import pulumi_aws as aws

# Get some configuration values or set default values.
config = pulumi.Config()
instance_type = config.get("instanceType")
if instance_type is None:
    instance_type = "t3.micro"
vpc_network_cidr = config.get("vpcNetworkCidr")
if vpc_network_cidr is None:
    vpc_network_cidr = "10.0.0.0/16"

# Look up the latest Amazon Linux 2 AMI.
ami = aws.ec2.get_ami(filters=[aws.ec2.GetAmiFilterArgs(
        name="name",
        values=["amzn2-ami-hvm-*"],
    )],
    owners=["amazon"],
    most_recent=True).id

# User data to start a HTTP server in the EC2 instance
user_data = """#!/bin/bash
echo "Hello, World from Pulumi!" > index.html
nohup python -m SimpleHTTPServer 80 &
"""

# Create VPC.
vpc = aws.ec2.Vpc("vpc",
    cidr_block=vpc_network_cidr,
    enable_dns_hostnames=True,
    enable_dns_support=True)

# Create an internet gateway.
gateway = aws.ec2.InternetGateway("gateway", vpc_id=vpc.id)

# Create a subnet that automatically assigns new instances a public IP address.
subnet = aws.ec2.Subnet("subnet",
    vpc_id=vpc.id,
    cidr_block="10.0.1.0/24",
    map_public_ip_on_launch=True)

# Create a route table.
route_table = aws.ec2.RouteTable("routeTable",
    vpc_id=vpc.id,
    routes=[aws.ec2.RouteTableRouteArgs(
        cidr_block="0.0.0.0/0",
        gateway_id=gateway.id,
    )])

# Associate the route table with the public subnet.
route_table_association = aws.ec2.RouteTableAssociation("routeTableAssociation",
    subnet_id=subnet.id,
    route_table_id=route_table.id)

# Create a security group allowing inbound access over port 80 and outbound
# access to anywhere.
sec_group = aws.ec2.SecurityGroup("secGroup",
    description="Enable HTTP access",
    vpc_id=vpc.id,
    ingress=[aws.ec2.SecurityGroupIngressArgs(
        from_port=80,
        to_port=80,
        protocol="tcp",
        cidr_blocks=["0.0.0.0/0"],
    )],
    egress=[aws.ec2.SecurityGroupEgressArgs(
        from_port=0,
        to_port=0,
        protocol="-1",
        cidr_blocks=["0.0.0.0/0"],
    )])

# Create and launch an EC2 instance into the public subnet.
server = aws.ec2.Instance("server",
    instance_type=instance_type,
    subnet_id=subnet.id,
    vpc_security_group_ids=[sec_group.id],
    user_data=user_data,
    ami=ami,
    tags={
        "Name": "webserver",
    })

# Export the instance's publicly accessible IP address and hostname.
pulumi.export("ip", server.public_ip)
pulumi.export("hostname", server.public_dns)
pulumi.export("url", server.public_dns.apply(lambda public_dns: f"http://{public_dns}"))
