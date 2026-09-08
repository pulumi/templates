import pulumi
import pulumi_cloudflare as cloudflare

# Import the program's configuration settings.
config = pulumi.Config()
account_id = config.require("accountId")
path = config.get("path") or "./www"

# A Cloudflare Worker, exposed on its workers.dev subdomain.
site = cloudflare.Worker(
    "site",
    account_id=account_id,
    name="my-site",
    subdomain=cloudflare.WorkerSubdomainArgs(
        enabled=True,
    ),
)

# A version of the Worker that serves the website's files as static assets.
version = cloudflare.WorkerVersion(
    "version",
    account_id=account_id,
    worker_id=site.id,
    compatibility_date="2025-01-01",
    assets=cloudflare.WorkerVersionAssetsArgs(
        directory=path,
        config=cloudflare.WorkerVersionAssetsConfigArgs(
            # Serve /404.html when a request doesn't match a file.
            not_found_handling="404-page",
            html_handling="auto-trailing-slash",
        ),
    ),
)

# Deploy the version so it serves all of the site's traffic.
deployment = cloudflare.WorkersDeployment(
    "deployment",
    account_id=account_id,
    script_name=site.name,
    strategy="percentage",
    versions=[
        cloudflare.WorkersDeploymentVersionArgs(
            version_id=version.id,
            percentage=100,
        )
    ],
)

# Export the website's URL.
pulumi.export("url", site.subdomain.url)
