import * as pulumi from "@pulumi/pulumi";
import * as gcp from "@pulumi/gcp";
import * as docker from "@pulumi/docker";

const config = new pulumi.Config();
const imageName = config.get("imageName") || "my-app";
const appPath = config.get("appPath") || "./app";
const containerPort = config.getNumber("containerPort") || 8080;
const cpu = config.getNumber("cpu") || 1;
const memory = config.get("memory") || "1Gi";
const concurrency = config.getNumber("concurrency") || 80;

const gcpConfig = new pulumi.Config("gcp");
const location = gcpConfig.require("region");
const project = gcpConfig.require("project");

const image = new docker.Image("image", {
    imageName: `gcr.io/${project}/${imageName}`,
    build: {
        context: appPath,
    },
});

const service = new gcp.cloudrun.Service("service", {
    location,
    template: {
        spec: {
            containers: [
                {
                    image: image.imageName,
                    resources: {
                        limits: {
                            memory,
                            cpu: cpu.toString(),
                        },
                    },
                    ports: [
                        {
                            containerPort,
                        },
                    ],
                }
            ],
            containerConcurrency: concurrency,
        },
    },
});

const invoker = new gcp.cloudrun.IamMember("invoker", {
    location,
    service: service.name,
    role: "roles/run.invoker",
    member: "allUsers",
});

export const url = service.statuses.apply(statuses => statuses[0]?.url);
