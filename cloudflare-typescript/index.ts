import * as pulumi from "@pulumi/pulumi";
import * as cloudflare from "@pulumi/cloudflare";

// Import the program's configuration settings.
const config = new pulumi.Config();
const accountId = config.require("accountId");

// Create a Cloudflare resource (Workers KV namespace).
const namespace = new cloudflare.WorkersKvNamespace("my-namespace", {
    accountId: accountId,
    title: "my-namespace",
});

// Export the ID of the namespace.
export const namespaceId = namespace.id;
