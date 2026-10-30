import re

import pulumi
import pulumi_cloudflare as cloudflare

# Import the program's configuration settings.
config = pulumi.Config()
account_id = config.require("accountId")

name = f"{pulumi.get_project()}-{pulumi.get_stack()}".lower()
name = re.sub(r"[^a-z0-9-]+", "-", name).strip("-")

# Create a Cloudflare resource (Workers KV namespace).
namespace = cloudflare.WorkersKvNamespace(
    "my-namespace",
    account_id=account_id,
    title=name,
)

# Export the ID of the namespace.
pulumi.export("namespaceId", namespace.id)
