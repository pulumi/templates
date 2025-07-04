using System.Collections.Generic;
using Pulumi;
using Kubernetes = Pulumi.Kubernetes;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var k8sNamespace = config.Get("k8sNamespace") ?? "ingress-nginx";
    var appLabels = new InputMap<string>
    {
        { "app", "ingress-nginx" },
    };

    var ingressns = new Kubernetes.Core.V1.Namespace("ingressns", new()
    {
        Metadata = new Kubernetes.Types.Inputs.Meta.V1.ObjectMetaArgs
        {
            Labels = appLabels,
            Name = k8sNamespace,
        },
    });

    var ingresscontroller = new Kubernetes.Helm.V3.Release("ingresscontroller", new()
    {
        Chart = "ingress-nginx",
        Namespace = ingressns.Metadata.Apply(m => m.Name),
        RepositoryOpts = new Kubernetes.Types.Inputs.Helm.V3.RepositoryOptsArgs
        {
            Repo = "https://kubernetes.github.io/ingress-nginx",
        },
        SkipCrds = true,
        Values = new Dictionary<string, object>
        {
            ["serviceAccount"] = new Dictionary<string, object>
            {
                ["automountServiceAccountToken"] = "true"
            },
            ["controller"] = new Dictionary<string, object>
            {
                ["publishService"] = new Dictionary<string, object>
                {
                    ["enabled"] = "true"
                },
            },
        },
        Version = "4.11.3",
    });

    return new Dictionary<string, object?>
    {
        ["name"] = ingresscontroller.Name,
    };
});
