import * as k8s from "@pulumi/kubernetes";
import * as kx from "@pulumi/kubernetesx";

// Define and create a Kubernetes Deployment.
const appLabels = { app: "nginx" };
const deployment = new k8s.apps.v1.Deployment("nginx", {
    spec: {
        selector: { matchLabels: appLabels },
        replicas: 1,
        template: {
            metadata: { labels: appLabels },
            spec: { containers: [{ name: "nginx", image: "nginx" }] }
        }
    }
});
export const name = deployment.metadata.name;

//
// This is equivalent code using the kx package.
//
const pb = new kx.PodBuilder({
    containers: [{ name: "nginx", image: "nginx" }]
});
const deploymentKx = new kx.Deployment("nginx-kx", {
    spec: pb.asDeploymentSpec({ replicas: 1 }),
});
export const nameKx = deploymentKx.metadata.name;
