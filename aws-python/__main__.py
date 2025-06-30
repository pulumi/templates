import pulumi

import infra

# Export the name of the bucket.
pulumi.export("bucket_name", infra.bucket.id)
