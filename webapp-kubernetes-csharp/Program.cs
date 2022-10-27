using System.Collections.Generic;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var namespace = config.Get("namespace") ?? "default";
    var replicas = config.GetNumber("replicas") ?? 1;
    var appLabels = 
    {
        { "app", "nginx" },
    };

    var webserver = new Kubernetes.Core.V1.Namespace("webserver", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Name = @namespace,
        },
    });

    var webserverconfig = new Kubernetes.Core.V1.ConfigMap("webserverconfig", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Namespace = @namespace,
        },
        Data = 
        {
            { "nginx.conf", @"events { }
http {
  server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html index.htm index.nginx-debian.html
    server_name _;
    location / {
      try_files $uri $uri/ =404;
    }
  }
}
" },
        },
    });

    var webserverdeployment = new Kubernetes.Apps.V1.Deployment("webserverdeployment", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Namespace = @namespace,
        },
        Spec = new Kubernetes.Types.Inputs.Apps.V1.DeploymentSpecArgs
        {
            Selector = new Kubernetes.Types.Inputs.Meta.V1.LabelSelectorArgs
            {
                MatchLabels = appLabels,
            },
            Replicas = 1,
            Template = new Kubernetes.Types.Inputs.Core.V1.PodTemplateSpecArgs
            {
                Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
                {
                    Labels = appLabels,
                },
                Spec = new Kubernetes.Types.Inputs.Core.V1.PodSpecArgs
                {
                    Containers = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.ContainerArgs
                        {
                            Image = "nginx",
                            Name = "nginx",
                            VolumeMounts = new[]
                            {
                                new Kubernetes.Types.Inputs.Core.V1.VolumeMountArgs
                                {
                                    MountPath = "/etc/nginx/nginx.conf",
                                    Name = "nginx-conf-volume",
                                    ReadOnly = true,
                                    SubPath = "nginx.conf",
                                },
                            },
                        },
                    },
                    Volumes = new[]
                    {
                        new Kubernetes.Types.Inputs.Core.V1.VolumeArgs
                        {
                            ConfigMap = new Kubernetes.Types.Inputs.Core.V1.ConfigMapVolumeSourceArgs
                            {
                                Items = new[]
                                {
                                    new Kubernetes.Types.Inputs.Core.V1.KeyToPathArgs
                                    {
                                        Key = "nginx.conf",
                                        Path = "nginx.conf",
                                    },
                                },
                                Name = webserverconfig.Metadata.Apply(metadata => metadata?.Name),
                            },
                            Name = "nginx-conf-volume",
                        },
                    },
                },
            },
        },
    });

    var webserverservice = new Kubernetes.Core.V1.Service("webserverservice", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Namespace = @namespace,
        },
        Spec = new Kubernetes.Types.Inputs.Core.V1.ServiceSpecArgs
        {
            Ports = new[]
            {
                new Kubernetes.Types.Inputs.Core.V1.ServicePortArgs
                {
                    Port = 80,
                    TargetPort = 80,
                    Protocol = "TCP",
                },
            },
            Selector = appLabels,
        },
    });

    return new Dictionary<string, object?>
    {
        ["deploymentName"] = webserverdeployment.Metadata.Apply(metadata => metadata?.Name),
        ["serviceName"] = webserverservice.Metadata.Apply(metadata => metadata?.Name),
    };
});

