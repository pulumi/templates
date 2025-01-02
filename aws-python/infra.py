import pulumi
from pulumi_aws import s3

# Create an AWS resource (S3 Bucket) with tags.
bucket = s3.BucketV2("my-bucket",
    tags={
        "Name": "My bucket",
    },
)
