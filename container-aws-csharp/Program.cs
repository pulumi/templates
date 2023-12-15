using System.Collections.Generic;
using Pulumi;
using Aws = Pulumi.Aws;
using Awsx = Pulumi.Awsx;

return await Deployment.RunAsync(() =>
{
    var config = new Config();
    var containerPort = config.GetInt32("containerPort") ?? 80;
    var cpu = config.GetInt32("cpu") ?? 512;
    var memory = config.GetInt32("memory") ?? 128;
    var cluster = new Aws.Ecs.Cluster("cluster");

    var loadbalancer = new Awsx.Lb.ApplicationLoadBalancer("loadbalancer");

    var repo = new Awsx.Ecr.Repository("repo", new()
    {
        ForceDelete = true,
    });

    var image = new Awsx.Ecr.Image("image", new()
    {
        RepositoryUrl = repo.Url,
        Context = "./app",
        Platform = "linux/amd64",
    });

    var service = new Awsx.Ecs.FargateService("service", new()
    {
        Cluster = cluster.Arn,
        AssignPublicIp = true,
        TaskDefinitionArgs = new Awsx.Ecs.Inputs.FargateServiceTaskDefinitionArgs
        {
            Container = new Awsx.Ecs.Inputs.TaskDefinitionContainerDefinitionArgs
            {
                Name = "app",
                Image = image.ImageUri,
                Cpu = cpu,
                Memory = memory,
                Essential = true,
                PortMappings = new[]
                {
                    new Awsx.Ecs.Inputs.TaskDefinitionPortMappingArgs
                    {
                        ContainerPort = containerPort,
                        TargetGroup = loadbalancer.DefaultTargetGroup,
                    },
                },
            },
        },
    });

    return new Dictionary<string, object?>
    {
        ["url"] = loadbalancer.LoadBalancer.Apply(loadBalancer => Output.Format($"http://{loadBalancer.DnsName}")),
    };
});
