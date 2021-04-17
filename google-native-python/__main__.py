"""A Google Cloud Python Pulumi program"""

import pulumi
from pulumi_google_native.storage import v1 as storage

config = pulumi.Config('google-native')
project = config.require('project')

# Create a Google Cloud resource (Storage Bucket)
bucket_name = "google-native-bucket-py-01"
bucket = storage.Bucket('my-bucket', name=bucket_name, bucket=bucket_name, project=project)

# Export the bucket self-link
pulumi.export('bucketSelfLink', bucket.self_link)

