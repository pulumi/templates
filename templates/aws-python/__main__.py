import pulumi
from pulumi_aws import s3

# Create an AWS resource (S3 Bucket)
bucket = s3.Bucket('my-bucket')

# Export the DNS name of the bucket
pulumi.output('bucket_name',  bucket.bucket_domain_name)
