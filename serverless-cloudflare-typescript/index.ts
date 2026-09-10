import * as pulumi from "@pulumi/pulumi";
import * as cloudflare from "@pulumi/cloudflare";

// Import the program's configuration settings.
const config = new pulumi.Config();
const accountId = config.require("accountId");

// A Workers KV namespace to persist state between requests.
const namespace = new cloudflare.WorkersKvNamespace("namespace", {
    accountId: accountId,
    title: "visit-counter",
});

// A Cloudflare Worker, exposed on its workers.dev subdomain.
const worker = new cloudflare.Worker("worker", {
    accountId: accountId,
    name: "my-app",
    subdomain: {
        enabled: true,
    },
});

// A version of the Worker containing the code and its KV binding.
const version = new cloudflare.WorkerVersion("version", {
    accountId: accountId,
    workerId: worker.id,
    compatibilityDate: "2025-01-01",
    mainModule: "worker.js",
    bindings: [{
        name: "COUNTER",
        type: "kv_namespace",
        namespaceId: namespace.id,
    }],
    modules: [{
        name: "worker.js",
        contentType: "application/javascript+module",
        contentFile: "worker.js",
    }],
});

// Deploy the version so it serves all of the application's traffic.
const deployment = new cloudflare.WorkersDeployment("deployment", {
    accountId: accountId,
    scriptName: worker.name,
    strategy: "percentage",
    versions: [{
        versionId: version.id,
        percentage: 100,
    }],
});

// Export the application's URL.
export const url = worker.subdomain.url;
