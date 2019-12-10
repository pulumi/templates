using System.Collections.Generic;
using System.Threading.Tasks;

using Pulumi;
using Pulumi.Kubernetes.Core.V1;
using Pulumi.Kubernetes.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Core.V1;
using Pulumi.Kubernetes.Types.Inputs.Apps.V1;
using Pulumi.Kubernetes.Types.Inputs.Meta.V1;
using Pulumi.Kubernetes.Types.Inputs.ApiExtensions.V1Beta1;

class Program
{
    static Task<int> Main()
    {
        return Pulumi.Deployment.RunAsync(() => {

            var appLabels = new InputMap<string>{
                { "app", "nginx" },
            };

            var deployment = new Pulumi.Kubernetes.Apps.V1.Deployment("nginx", new DeploymentArgs
            {
                Spec = new DeploymentSpecArgs
                {
                    Selector = new LabelSelectorArgs
                    {
                        MatchLabels = appLabels,
                    },
                    Replicas = 1,
                    Template = new PodTemplateSpecArgs
                    {
                        Metadata = new ObjectMetaArgs
                        {
                            Labels = appLabels,
                        },
                        Spec = new PodSpecArgs
                        {
                            Containers =
                            {
                                new ContainerArgs
                                {
                                    Name = "nginx",
                                    Image = "nginx",
                                    Ports =
                                    {
                                        new ContainerPortArgs { ContainerPortValue = 80 }
                                    },
                                },
                            },
                        },
                    },
                },
            });
            
            return new Dictionary<string, object>
            {
                { "name", deployment.Metadata.Apply(m => m.Name) },
            };
        });
    }
}
