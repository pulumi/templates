import * as pulumi from "@pulumi/pulumi";
import * as aiven from "@pulumi/aiven";

// Import the program's configuration settings.
const config = new pulumi.Config();
const project = config.require("project");
const cloudName = config.require("cloudName");
const plan = config.require("plan");
const serviceName = config.require("serviceName");

// Create a Kafka service.
const kafka = new aiven.Kafka("kafka", {
    project: project,
    cloudName: cloudName,
    plan: plan,
    serviceName: serviceName,
});

// Export the service host and port.
export const serviceHost = kafka.serviceHost;
export const servicePort = kafka.servicePort;
