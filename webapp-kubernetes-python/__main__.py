import pulumi
import pulumi_kubernetes as kubernetes

# Get some values from the Pulumi stack configuration, or use defaults
config = pulumi.Config()
k8sNamespace = config.get("namespace", "default")
numReplicas = config.get_int("replicas", 1)
app_labels = {
    "app": "nginx",
}

# Create a namespace
webserverns = kubernetes.core.v1.Namespace(
    "webserver",
    metadata={
        "name": k8sNamespace,
    },
)

# Create a ConfigMap for the Nginx configuration
webserverconfig = kubernetes.core.v1.ConfigMap(
    "webserverconfig",
    metadata={
        "namespace": webserverns.metadata.name,
    },
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
    },
)

# Create a Deployment with a user-defined number of replicas
webserverdeployment = kubernetes.apps.v1.Deployment(
    "webserverdeployment",
    metadata={
        "namespace": webserverns.metadata.name,
    },
    spec={
        "selector": {
            "match_labels": app_labels,
        },
        "replicas": numReplicas,
        "template": {
            "metadata": {
                "labels": app_labels,
            },
            "spec": {
                "containers": [
                    {
                        "image": "nginx",
                        "name": "nginx",
                        "volume_mounts": [
                            {
                                "mount_path": "/etc/nginx/nginx.conf",
                                "name": "nginx-conf-volume",
                                "read_only": True,
                                "sub_path": "nginx.conf",
                            }
                        ],
                    }
                ],
                "volumes": [
                    {
                        "config_map": {
                            "items": [
                                {
                                    "key": "nginx.conf",
                                    "path": "nginx.conf",
                                }
                            ],
                            "name": webserverconfig.metadata.name,
                        },
                        "name": "nginx-conf-volume",
                    }
                ],
            },
        },
    },
)

# Expose the Deployment as a Kubernetes Service
webserverservice = kubernetes.core.v1.Service(
    "webserverservice",
    metadata={
        "namespace": webserverns.metadata.name,
    },
    spec={
        "ports": [
            {
                "port": 80,
                "target_port": 80,
                "protocol": "TCP",
            }
        ],
        "selector": app_labels,
    },
)

# Export some values for use elsewhere
pulumi.export("deploymentName", webserverdeployment.metadata.name)
pulumi.export("serviceName", webserverservice.metadata.name)
