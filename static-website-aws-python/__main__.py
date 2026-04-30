import json

import pulumi
import pulumi_aws as aws
import pulumi_synced_folder as synced_folder

# Import the program's configuration settings.
config = pulumi.Config()
path = config.get("path") or "./www"
index_document = config.get("indexDocument") or "index.html"
error_document = config.get("errorDocument") or "error.html"

# Create a private S3 bucket to hold the website content.
bucket = aws.s3.Bucket("bucket")

# Block all public access to the bucket; CloudFront will reach it via OAC.
public_access_block = aws.s3.BucketPublicAccessBlock(
    "public-access-block",
    bucket=bucket.bucket,
    block_public_acls=True,
    block_public_policy=True,
    ignore_public_acls=True,
    restrict_public_buckets=True,
)

# Sync the website files to the bucket as private objects.
bucket_folder = synced_folder.S3BucketFolder(
    "bucket-folder",
    acl="private",
    bucket_name=bucket.bucket,
    path=path,
    opts=pulumi.ResourceOptions(depends_on=[public_access_block]),
)

# Create an Origin Access Control so CloudFront can read from the private bucket.
origin_access_control = aws.cloudfront.OriginAccessControl(
    "origin-access-control",
    origin_access_control_origin_type="s3",
    signing_behavior="always",
    signing_protocol="sigv4",
)

# Create a CloudFront CDN to distribute and cache the website.
cdn = aws.cloudfront.Distribution(
    "cdn",
    enabled=True,
    default_root_object=index_document,
    origins=[
        {
            "origin_id": bucket.arn,
            "domain_name": bucket.bucket_regional_domain_name,
            "origin_access_control_id": origin_access_control.id,
        }
    ],
    default_cache_behavior={
        "target_origin_id": bucket.arn,
        "viewer_protocol_policy": "redirect-to-https",
        "allowed_methods": ["GET", "HEAD", "OPTIONS"],
        "cached_methods": ["GET", "HEAD", "OPTIONS"],
        "compress": True,
        # Managed-CachingOptimized
        "cache_policy_id": "658327ea-f89d-4fab-a63d-7e88639e58f6",
        # Managed-SecurityHeadersPolicy
        "response_headers_policy_id": "67f7725c-6f97-4210-82d7-5512b31e9d03",
    },
    price_class="PriceClass_100",
    custom_error_responses=[
        {
            "error_code": 404,
            "response_code": 404,
            "response_page_path": f"/{error_document}",
        }
    ],
    restrictions={
        "geo_restriction": {
            "restriction_type": "none",
        },
    },
    viewer_certificate={
        "cloudfront_default_certificate": True,
    },
)

# Grant the CloudFront distribution permission to read objects from the bucket.
bucket_policy = aws.s3.BucketPolicy(
    "bucket-policy",
    bucket=bucket.bucket,
    policy=pulumi.Output.all(bucket_arn=bucket.arn, cdn_arn=cdn.arn).apply(
        lambda args: json.dumps(
            {
                "Version": "2012-10-17",
                "Statement": [
                    {
                        "Sid": "AllowCloudFrontServicePrincipalReadOnly",
                        "Effect": "Allow",
                        "Principal": {"Service": "cloudfront.amazonaws.com"},
                        "Action": "s3:GetObject",
                        "Resource": f"{args['bucket_arn']}/*",
                        "Condition": {
                            "StringEquals": {"AWS:SourceArn": args["cdn_arn"]},
                        },
                    }
                ],
            }
        )
    ),
)

# Export the URLs and hostnames of the bucket and distribution.
pulumi.export("originHostname", bucket.bucket_regional_domain_name)
pulumi.export("cdnURL", pulumi.Output.concat("https://", cdn.domain_name))
pulumi.export("cdnHostname", cdn.domain_name)
