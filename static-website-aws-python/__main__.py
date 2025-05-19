import pulumi
import pulumi_aws as aws
import pulumi_synced_folder as synced_folder

# Import the program's configuration settings.
config = pulumi.Config()
path = config.get("path") or "./www"
index_document = config.get("indexDocument") or "index.html"
error_document = config.get("errorDocument") or "error.html"

# Create an S3 bucket and configure it as a website.
bucket = aws.s3.BucketV2(
    "bucket",
    website={
        "index_document": index_document,
        "error_document": error_document,
    },
)

bucket_website = aws.s3.BucketWebsiteConfigurationV2(
    "bucket",
    bucket=bucket.bucket,
    index_document={"suffix": index_document},
    error_document={"key": error_document},
)

# Set ownership controls for the new bucket
ownership_controls = aws.s3.BucketOwnershipControls(
    "ownership-controls",
    bucket=bucket.bucket,
    rule={
        "object_ownership": "ObjectWriter",
    },
)

# Configure public ACL block on the new bucket
public_access_block = aws.s3.BucketPublicAccessBlock(
    "public-access-block",
    bucket=bucket.bucket,
    block_public_acls=False,
)

# Use a synced folder to manage the files of the website.
bucket_folder = synced_folder.S3BucketFolder(
    "bucket-folder",
    acl="public-read",
    bucket_name=bucket.bucket,
    path=path,
    opts=pulumi.ResourceOptions(depends_on=[ownership_controls, public_access_block]),
)

# Create a CloudFront CDN to distribute and cache the website.
cdn = aws.cloudfront.Distribution(
    "cdn",
    enabled=True,
    origins=[
        {
            "origin_id": bucket.arn,
            "domain_name": bucket_website.website_endpoint,
            "custom_origin_config": {
                "origin_protocol_policy": "http-only",
                "http_port": 80,
                "https_port": 443,
                "origin_ssl_protocols": ["TLSv1.2"],
            },
        }
    ],
    default_cache_behavior={
        "target_origin_id": bucket.arn,
        "viewer_protocol_policy": "redirect-to-https",
        "allowed_methods": [
            "GET",
            "HEAD",
            "OPTIONS",
        ],
        "cached_methods": [
            "GET",
            "HEAD",
            "OPTIONS",
        ],
        "default_ttl": 600,
        "max_ttl": 600,
        "min_ttl": 600,
        "forwarded_values": {
            "query_string": True,
            "cookies": {
                "forward": "all",
            },
        },
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

# Export the URLs and hostnames of the bucket and distribution.
pulumi.export("originURL", pulumi.Output.concat("http://", bucket_website.website_endpoint))
pulumi.export("originHostname", bucket.website_endpoint)
pulumi.export("cdnURL", pulumi.Output.concat("https://", cdn.domain_name))
pulumi.export("cdnHostname", cdn.domain_name)
