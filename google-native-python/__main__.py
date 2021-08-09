"""A Google Cloud Python Pulumi program"""

import pulumi
from pulumi_google_native.storage import v1 as storage

# Create a Google Cloud resource (Storage Bucket)
bucket = storage.Bucket('my-bucket')

# Export the bucket self-link
pulumi.export('bucketSelfLink', bucket.self_link)

