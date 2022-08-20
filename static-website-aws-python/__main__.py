import pulumi
import pulumi_aws as aws
import pulumi_synced_folder as synced_folder

config = pulumi.Config()
aws_region = config.get("awsRegion")
if aws_region is None:
    aws_region = "us-west-2"
path = config.get("path")
if path is None:
    path = "./site"
index_document = config.get("indexDocument")
if index_document is None:
    index_document = "index.html"
error_document = config.get("errorDocument")
if error_document is None:
    error_document = "error.html"
bucket = aws.s3.Bucket("bucket",
    acl="public-read",
    website=aws.s3.BucketWebsiteArgs(
        index_document=index_document,
        error_document=error_document,
    ))
bucket_folder = synced_folder.S3BucketFolder("bucket-folder",
    path=path,
    bucket_name=bucket.bucket,
    acl="public-read")
pulumi.export("url", bucket.website_endpoint.apply(lambda website_endpoint: f"http://{website_endpoint}"))
