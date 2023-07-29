import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

// Get some values from the stack configuration, or use defaults
const config = new pulumi.Config();
const k8sNamespace = config.get("namespace") || "default";
const numReplicas = config.getNumber("replicas") || 1;
const appLabels = {
    app: "nginx",
};

// Create a new namespace
const webServerNs = new kubernetes.core.v1.Namespace("webserver", {metadata: {
    name: k8sNamespace,
}});

// Create a new ConfigMap for the Nginx configuration
const webServerConfig = new kubernetes.core.v1.ConfigMap("webserverconfig", {
    metadata: {
        namespace: webServerNs.metadata.name,
    },
    data: {
        "nginx.conf": `events { }
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
`,
    },
});

// Create a new Deployment with a user-specified number of replicas
const webServerDeployment = new kubernetes.apps.v1.Deployment("webserverdeployment", {
    metadata: {
        namespace: webServerNs.metadata.name,
    },
    spec: {
        selector: {
            matchLabels: appLabels,
        },
        replicas: numReplicas,
        template: {
            metadata: {
                labels: appLabels,
            },
            spec: {
                containers: [{
                    image: "nginx",
                    name: "nginx",
                    volumeMounts: [{
                        mountPath: "/etc/nginx/nginx.conf",
                        name: "nginx-conf-volume",
                        readOnly: true,
                        subPath: "nginx.conf",
                    }],
                }],
                volumes: [{
                    configMap: {
                        items: [{
                            key: "nginx.conf",
                            path: "nginx.conf",
                        }],
                        name: webServerConfig.metadata.name,
                    },
                    name: "nginx-conf-volume",
                }],
            },
        },
    },
});

// Expose the Deployment as a Kubernetes Service
const webServerService = new kubernetes.core.v1.Service("webserverservice", {
    metadata: {
        namespace: webServerNs.metadata.name,
    },
    spec: {
        ports: [{
            port: 80,
            targetPort: 80,
            protocol: "TCP",
        }],
        selector: appLabels,
    },
});

// Export some values for use elsewhere
export const deploymentName = webServerDeployment.metadata.name;
export const serviceName = webServerService.metadata.name;
