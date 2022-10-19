using Pulumi;
using Gcp = Pulumi.Gcp;
using Docker = Pulumi.Docker;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var config = new Config();
    var appPath = config.Get("appPath") ?? "./app";
    var imageName = config.Get("imageName") ?? "my-app";
    var containerPort = config.GetInt32("containerPort") ?? 8080;
    var cpu = config.GetInt32("cpu") ?? 1;
    var memory = config.Get("memory") ?? "1Gi";
    var concurrency = config.GetInt32("concurrency") ?? 50;

    var gcpConfig = new Config("gcp");
    var location = gcpConfig.Require("region");
    var project = gcpConfig.Require("project");

    var image = new Docker.Image("image", new()
    {
        ImageName = $"gcr.io/{project}/{imageName}",
        Build = new Docker.DockerBuild {
            Context = appPath,
        },
    });

    var service = new Gcp.CloudRun.Service("service", new Gcp.CloudRun.ServiceArgs
    {
        Location = location!,
        Template = new Gcp.CloudRun.Inputs.ServiceTemplateArgs
        {
            Spec = new Gcp.CloudRun.Inputs.ServiceTemplateSpecArgs
            {
                Containers = new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerArgs[]
                {
                    new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerArgs
                    {
                        Image = image.ImageName,
                        Resources = new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerResourcesArgs
                        {
                            Limits = {
                                { "memory", memory },
                                { "cpu", cpu.ToString() },
                            },
                        },
                        Ports = new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerPortArgs[]
                        {
                            new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerPortArgs
                            {
                                ContainerPort = containerPort,
                            },
                        },
                        Envs = new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerEnvArgs[]
                        {
                            new Gcp.CloudRun.Inputs.ServiceTemplateSpecContainerEnvArgs
                            {
                                Name = "ASPNETCORE_URLS",
                                Value = $"http://*:{containerPort}",
                            },
                        },
                    },
                },
                ContainerConcurrency = concurrency,
            },
        },
    });

    var invoker = new Gcp.CloudRun.IamMember("invoker", new Gcp.CloudRun.IamMemberArgs
    {
        Location = location!,
        Service = service.Name,
        Role = "roles/run.invoker",
        Member = "allUsers",
    });

    return new Dictionary<string, object?>
    {
        ["url"] = service.Statuses.Apply(statuses => statuses[0]?.Url),
    };
});
