using System.Collections.Generic;
using Pulumi;
using Aws = Pulumi.Aws;

return await Deployment.RunAsync(() => 
{
    // Get some configuration values or set default values.
    var config = new Config();
    var instanceType = config.Get("instanceType") ?? "t3.micro";
    var vpcNetworkCidr = config.Get("vpcNetworkCidr") ?? "10.0.0.0/16";
    
    // Look up the latest Amazon Linux 2 AMI.
    var ami = Aws.Ec2.GetAmi.Invoke(new()
    {
        Filters = new[]
        {
            new Aws.Ec2.Inputs.GetAmiFilterInputArgs
            {
                Name = "name",
                Values = new[]
                {
                    "amzn2-ami-hvm-*",
                },
            },
        },
        Owners = new[]
        {
            "amazon",
        },
        MostRecent = true,
    }).Apply(invoke => invoke.Id);

    // User data to start a HTTP server in the EC2 instance
    var userData = @"#!/bin/bash
    echo ""Hello, World from Pulumi!"" > index.html
    nohup python -m SimpleHTTPServer 80 &
    ";

    // Create VPC.
    var vpc = new Aws.Ec2.Vpc("vpc", new()
    {
        CidrBlock = vpcNetworkCidr,
        EnableDnsHostnames = true,
        EnableDnsSupport = true,
    });

    // Create an internet gateway.
    var gateway = new Aws.Ec2.InternetGateway("gateway", new()
    {
        VpcId = vpc.Id,
    });

    // Create a subnet that automatically assigns new instances a public IP address.
    var subnet = new Aws.Ec2.Subnet("subnet", new()
    {
        VpcId = vpc.Id,
        CidrBlock = "10.0.1.0/24",
        MapPublicIpOnLaunch = true,
    });

    // Create a route table.
    var routeTable = new Aws.Ec2.RouteTable("routeTable", new()
    {
        VpcId = vpc.Id,
        Routes = new[]
        {
            new Aws.Ec2.Inputs.RouteTableRouteArgs
            {
                CidrBlock = "0.0.0.0/0",
                GatewayId = gateway.Id,
            },
        },
    });

    // Associate the route table with the public subnet.
    var routeTableAssociation = new Aws.Ec2.RouteTableAssociation("routeTableAssociation", new()
    {
        SubnetId = subnet.Id,
        RouteTableId = routeTable.Id,
    });

    // Create a security group allowing inbound access over port 80 and outbound
    // access to anywhere.
    var secGroup = new Aws.Ec2.SecurityGroup("secGroup", new()
    {
        Description = "Enable HTTP access",
        VpcId = vpc.Id,
        Ingress = new[]
        {
            new Aws.Ec2.Inputs.SecurityGroupIngressArgs
            {
                FromPort = 80,
                ToPort = 80,
                Protocol = "tcp",
                CidrBlocks = new[]
                {
                    "0.0.0.0/0",
                },
            },
        },
        Egress = new[]
        {
            new Aws.Ec2.Inputs.SecurityGroupEgressArgs
            {
                FromPort = 0,
                ToPort = 0,
                Protocol = "-1",
                CidrBlocks = new[]
                {
                    "0.0.0.0/0",
                },
            },
        },
    });

    // Create and launch an EC2 instance into the public subnet.
    var server = new Aws.Ec2.Instance("server", new()
    {
        InstanceType = instanceType,
        SubnetId = subnet.Id,
        VpcSecurityGroupIds = new[]
        {
            secGroup.Id,
        },
        UserData = userData,
        Ami = ami,
        Tags = 
        {
            { "Name", "webserver" },
        },
    });

    // Export the instance's publicly accessible IP address and hostname.
    return new Dictionary<string, object?>
    {
        ["ip"] = server.PublicIp,
        ["hostname"] = server.PublicDns,
        ["url"] = server.PublicDns.Apply(publicDns => $"http://{publicDns}"),
    };
});

