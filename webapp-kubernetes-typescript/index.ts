import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const config = new pulumi.Config();
const namespace = config.get("namespace") || "default";
const numReplicas = config.getNumber("replicas") || 1;
const appLabels = {
    app: "nginx",
};
const webserver = new kubernetes.core.v1.Namespace("webserver", {metadata: {
    name: namespace,
}});
const webserverconfig = new kubernetes.core.v1.ConfigMap("webserverconfig", {
    metadata: {
        namespace: namespace,
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
const webserverdeployment = new kubernetes.apps.v1.Deployment("webserverdeployment", {
    metadata: {
        namespace: namespace,
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
                        name: webserverconfig.metadata.apply(metadata => metadata?.name),
                    },
                    name: "nginx-conf-volume",
                }],
            },
        },
    },
});
const webserverservice = new kubernetes.core.v1.Service("webserverservice", {
    metadata: {
        namespace: namespace,
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
export const deploymentName = webserverdeployment.metadata.apply(metadata => metadata?.name);
export const serviceName = webserverservice.metadata.apply(metadata => metadata?.name);
