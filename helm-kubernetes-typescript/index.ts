import * as pulumi from "@pulumi/pulumi";
import * as kubernetes from "@pulumi/kubernetes";

const config = new pulumi.Config();
const k8sNamespace = config.get("k8sNamespace") || "default";
const appLabels = {
    app: "ingress-nginx",
};

// Create a namespace (user supplies the name of the namespace)
const ingressNs = new kubernetes.core.v1.Namespace("ingressns", {
    metadata: {
        labels: appLabels,
        name: k8sNamespace,
    }
});

// Use Helm to install the Nginx ingress controller
const ingressController = new kubernetes.helm.v3.Release("ingresscontroller", {
    chart: "ingress-nginx",
    namespace: ingressNs.metadata.name,
    repositoryOpts: {
        repo: "https://kubernetes.github.io/ingress-nginx",
    },
    skipCrds: true,
    values: {
        serviceAccount: {
            automountServiceAccountToken: true,
        },
        controller: {
            publishService: {
                enabled: true,
            },
        },
    },
    version: "4.11.3",
});

// Export some values for use elsewhere
export const name = ingressController.name;
