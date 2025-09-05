import * as pulumi from "@pulumi/pulumi";
import * as vultr from "@ediri/vultr";

const config = new pulumi.Config();
const instanceName = config.get("instanceName") || "my-instance";
const plan = config.get("plan") || "vc2-1c-1gb";
const region = config.get("region") || "ewr";
const hostname = config.get("hostName") || "my-instance-hostname";
const tags = config.get("tags") ? config.get("tags")!.split(",") : ["web", "dev"];
const label = config.get("label") || "my-label";

// Create a new instance:
const myInstance = new vultr.Instance(instanceName, {
    osId: 1743,
    plan,
    region,
});

// Create a new instance with options:
const vultrInstance = new vultr.Instance(instanceName, {
    activationEmail: false,
    backups: "enabled",
    backupsSchedule: {
        type: "daily",
    },
    ddosProtection: true,
    disablePublicIpv4: true,
    enableIpv6: true,
    hostname,
    label,
    osId: 1743, // Ubuntu 22.04 LTS x64
    plan,
    region,
    tags,
});

export const instanceId = vultrInstance.id;
export const instanceStatus = vultrInstance.status;
export const mainIp = vultrInstance.mainIp;
