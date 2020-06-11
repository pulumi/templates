"""An AliCloud Python Pulumi program"""

import pulumi
from pulumi_alicloud import oss

# Create an AliCloud resource (OSS Bucket)
bucket = oss.Bucket('my-bucket')

# Export the name of the bucket
pulumi.export('bucket_name', bucket.id)
