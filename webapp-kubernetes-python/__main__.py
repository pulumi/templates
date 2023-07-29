import pulumi
import pulumi_kubernetes as kubernetes

# Get some values from the Pulumi stack configuration, or use defaults
config = pulumi.Config()
k8sNamespace = config.get("namespace", "default")
numReplicas = config.get_float("replicas", 1)
app_labels = {
    "app": "nginx",
}

# Create a namespace
webserverns = kubernetes.core.v1.Namespace(
    "webserver",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        name=k8sNamespace,
    )
)

# Create a ConfigMap for the Nginx configuration
webserverconfig = kubernetes.core.v1.ConfigMap(
    "webserverconfig",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        namespace=webserverns.metadata.name,
    ),
    data={
        "nginx.conf": """events { }
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
""",
    }
)

# Create a Deployment with a user-defined number of replicas
webserverdeployment = kubernetes.apps.v1.Deployment(
    "webserverdeployment",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        namespace=webserverns.metadata.name,
    ),
    spec=kubernetes.apps.v1.DeploymentSpecArgs(
        selector=kubernetes.meta.v1.LabelSelectorArgs(
            match_labels=app_labels,
        ),
        replicas=numReplicas,
        template=kubernetes.core.v1.PodTemplateSpecArgs(
            metadata=kubernetes.meta.v1.ObjectMetaArgs(
                labels=app_labels,
            ),
            spec=kubernetes.core.v1.PodSpecArgs(
                containers=[kubernetes.core.v1.ContainerArgs(
                    image="nginx",
                    name="nginx",
                    volume_mounts=[kubernetes.core.v1.VolumeMountArgs(
                        mount_path="/etc/nginx/nginx.conf",
                        name="nginx-conf-volume",
                        read_only=True,
                        sub_path="nginx.conf",
                    )],
                )],
                volumes=[kubernetes.core.v1.VolumeArgs(
                    config_map=kubernetes.core.v1.ConfigMapVolumeSourceArgs(
                        items=[kubernetes.core.v1.KeyToPathArgs(
                            key="nginx.conf",
                            path="nginx.conf",
                        )],
                        name=webserverconfig.metadata.name,
                    ),
                    name="nginx-conf-volume",
                )],
            ),
        ),
    )
)

# Expose the Deployment as a Kubernetes Service
webserverservice = kubernetes.core.v1.Service(
    "webserverservice",
    metadata=kubernetes.meta.v1.ObjectMetaArgs(
        namespace=webserverns.metadata.name,
    ),
    spec=kubernetes.core.v1.ServiceSpecArgs(
        ports=[kubernetes.core.v1.ServicePortArgs(
            port=80,
            target_port=80,
            protocol="TCP",
        )],
        selector=app_labels,
    )
)

# Export some values for use elsewhere
pulumi.export("deploymentName", webserverdeployment.metadata.name)
pulumi.export("serviceName", webserverservice.metadata.name)
