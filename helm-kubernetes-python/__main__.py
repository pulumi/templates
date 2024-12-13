import pulumi
import pulumi_kubernetes as kubernetes

config = pulumi.Config()
k8s_namespace = config.get("k8sNamespace", "default")
app_labels = {
    "app": "ingress-nginx",
}

# Create a namespace (user supplies the name of the namespace)
ingress_ns = kubernetes.core.v1.Namespace(
    "ingressns",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        labels=app_labels,
        name=k8s_namespace,
    )
)

# Use Helm to install the Nginx ingress controller
ingresscontroller = kubernetes.helm.v3.Release(
    "ingresscontroller",
    chart="ingress-nginx",
    namespace=ingress_ns.metadata.name,
    repository_opts={
        "repo": "https://kubernetes.github.io/ingress-nginx",
    },
    skip_crds=True,
    values={
        "serviceAccount": {
            "automountServiceAccountToken": True,
        },
        "controller": {
            "publishService": {
                "enabled": True,
            },
        },
    },
    version="4.11.3"
)

# Export some values for use elsewhere
pulumi.export("name", ingresscontroller.name)
