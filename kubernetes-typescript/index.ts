import * as pulumi from "@pulumi/pulumi";
import * as k8s from "@pulumi/kubernetes";

const name = pulumi.getProject();

const config = new pulumi.Config("kubeapp");
const image = config.require("image");
const replicas = config.getNumber("replicas") || 1;

const appLabels = { app: name };
const deployment = new k8s.apps.v1beta1.Deployment(name, {
    spec: {
        selector: { matchLabels: appLabels },
        replicas: replicas,
        template: {
            metadata: { labels: appLabels },
            spec: { containers: [{ name: name, image: image }] }
        }
    }
});
