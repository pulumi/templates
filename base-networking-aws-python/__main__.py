import pulumi
import pulumi_aws as aws
import pulumi_awsx as awsx

# Get some values from the Pulumi configuration (or use defaults)
config = pulumi.Config()
vpc_network_cidr = config.get("vpcNetworkCidr", "10.0.0.0/16")

# Create a new VPC
vpc = awsx.ec2.Vpc(
    "vpc",
    enable_dns_hostnames=True,
    cidr_block=vpc_network_cidr,
    tags={
        "https": "allow",
    }
)

# Create a new security group to allow HTTPS traffic
security_group = aws.ec2.SecurityGroup(
    "securityGroup",
    description="Allow HTTPS inbound traffic",
    vpc_id=vpc.vpc_id,
    ingress=[aws.ec2.SecurityGroupIngressArgs(
        description="HTTPS from VPC",
        from_port=443,
        to_port=443,
        protocol="tcp",
        cidr_blocks=[vpc_network_cidr],
    )],
    egress=[aws.ec2.SecurityGroupEgressArgs(
        from_port=0,
        to_port=0,
        protocol="-1",
        cidr_blocks=["0.0.0.0/0"],
    )],
    tags={
        "Name": "https",
    }
)

# Export some values for use elsewhere
pulumi.export("vpcId", vpc.vpc_id)
pulumi.export("securityGroupId", security_group.id)
