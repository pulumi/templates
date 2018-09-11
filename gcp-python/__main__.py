import pulumi
from pulumi_gcp import storage

# Create a GCP resource (Storage Bucket)
bucket = storage.Bucket('my-bucket')

# Export the DNS name of the bucket
pulumi.output('bucket_name',  bucket.url)
