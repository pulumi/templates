"""A Google Cloud Python Pulumi program"""

import pulumi
from pulumi_google_native.storage import v1 as storage
from pulumi_random import random_string

config = pulumi.Config('google-native')
project = config.require('project')

# Generate a random bucket name
suffix = random_string.RandomString('suffix', length=5, number=False, upper=False, special=False)
bucket_name = pulumi.Output.concat("pulumi-goog-native-bucket-py-", suffix.result)

# Create a Google Cloud resource (Storage Bucket)
bucket = storage.Bucket('my-bucket', name=bucket_name, bucket=bucket_name, project=project)

# Export the bucket self-link
pulumi.export('bucketSelfLink', bucket.self_link)

