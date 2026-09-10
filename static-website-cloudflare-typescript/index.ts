import * as pulumi from "@pulumi/pulumi";
import * as cloudflare from "@pulumi/cloudflare";

// Import the program's configuration settings.
const config = new pulumi.Config();
const accountId = config.require("accountId");
const path = config.get("path") || "./www";

// A Cloudflare Worker, exposed on its workers.dev subdomain.
const site = new cloudflare.Worker("site", {
    accountId: accountId,
    name: "my-site",
    subdomain: {
        enabled: true,
    },
});

// A version of the Worker that serves the website's files as static assets.
const version = new cloudflare.WorkerVersion("version", {
    accountId: accountId,
    workerId: site.id,
    compatibilityDate: "2025-01-01",
    assets: {
        directory: path,
        config: {
            // Serve /404.html when a request doesn't match a file.
            notFoundHandling: "404-page",
            htmlHandling: "auto-trailing-slash",
        },
    },
});

// Deploy the version so it serves all of the site's traffic.
const deployment = new cloudflare.WorkersDeployment("deployment", {
    accountId: accountId,
    scriptName: site.name,
    strategy: "percentage",
    versions: [{
        versionId: version.id,
        percentage: 100,
    }],
});

// Export the website's URL.
export const url = site.subdomain.url;
