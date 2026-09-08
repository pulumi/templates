import pulumi
import pulumi_cloudflare as cloudflare

# Import the program's configuration settings.
config = pulumi.Config()
account_id = config.require("accountId")

# A Workers KV namespace to persist state between requests.
namespace = cloudflare.WorkersKvNamespace(
    "namespace",
    account_id=account_id,
    title="visit-counter",
)

# A Cloudflare Worker, exposed on its workers.dev subdomain.
worker = cloudflare.Worker(
    "worker",
    account_id=account_id,
    name="my-app",
    subdomain=cloudflare.WorkerSubdomainArgs(
        enabled=True,
    ),
)

# A version of the Worker containing the code and its KV binding.
version = cloudflare.WorkerVersion(
    "version",
    account_id=account_id,
    worker_id=worker.id,
    compatibility_date="2025-01-01",
    main_module="worker.js",
    bindings=[
        cloudflare.WorkerVersionBindingArgs(
            name="COUNTER",
            type="kv_namespace",
            namespace_id=namespace.id,
        )
    ],
    modules=[
        cloudflare.WorkerVersionModuleArgs(
            name="worker.js",
            content_type="application/javascript+module",
            content_file="worker.js",
        )
    ],
)

# Deploy the version so it serves all of the application's traffic.
deployment = cloudflare.WorkersDeployment(
    "deployment",
    account_id=account_id,
    script_name=worker.name,
    strategy="percentage",
    versions=[
        cloudflare.WorkersDeploymentVersionArgs(
            version_id=version.id,
            percentage=100,
        )
    ],
)

# Export the application's URL.
pulumi.export("url", worker.subdomain.url)
