import * as pulumi from "@pulumi/pulumi";
import * as aiven from "@pulumi/aiven";

const kafka1 = new aiven.Kafka("kafka1", {
    project: "<YOUR_AIVEN_PROJECT_NAME>",
    cloudName: "google-europe-west1",
    plan: "business-4",
    serviceName: "kafka-gcp-eu",
    maintenanceWindowDow: "monday",
    maintenanceWindowTime: "10:00:00",
    kafkaUserConfig: {
        kafkaRest: "true",
        kafkaConnect: "true",
        schemaRegistry: "true",
        kafkaVersion: "3.2",
        kafka: {
            groupMaxSessionTimeoutMs: "70000",
            logRetentionBytes: "1000000000",
        },
        publicAccess: {
            kafkaRest: "true",
            kafkaConnect: "true",
        },
    },
});
