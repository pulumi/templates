import pulumi
import pulumi_cloudflare as cloudflare

# Import the program's configuration settings.
config = pulumi.Config()
account_id = config.require("accountId")

# Create a Cloudflare resource (Workers KV namespace).
namespace = cloudflare.WorkersKvNamespace(
    "my-namespace",
    account_id=account_id,
    title="my-namespace",
)

# Export the ID of the namespace.
pulumi.export("namespaceId", namespace.id)
