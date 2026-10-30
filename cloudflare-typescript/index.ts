import * as pulumi from "@pulumi/pulumi";
import * as cloudflare from "@pulumi/cloudflare";

// Import the program's configuration settings.
const config = new pulumi.Config();
const accountId = config.require("accountId");

const name = `${pulumi.getProject()}-${pulumi.getStack()}`
    .toLowerCase()
    .replace(/[^a-z0-9-]+/g, "-")
    .replace(/^-+|-+$/g, "");

// Create a Cloudflare resource (Workers KV namespace).
const namespace = new cloudflare.WorkersKvNamespace("my-namespace", {
    accountId: accountId,
    title: name,
});

// Export the ID of the namespace.
export const namespaceId = namespace.id;
