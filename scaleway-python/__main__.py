import pulumi
import lbrlabs_pulumi_scaleway as scaleway

# Create a Scaleway resource (Object Bucket).
bucket = scaleway.ObjectBucket("my-bucket")

# Export the name of the bucket.
pulumi.export("bucket_name", bucket.id)
