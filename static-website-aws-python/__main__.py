import pulumi
import pulumi_aws as aws
import pulumi_synced_folder as synced_folder

# Import the program's configuration settings.
config = pulumi.Config()
path = config.get("path") or "./site"
index_document = config.get("indexDocument") or "index.html"
error_document = config.get("errorDocument") or "error.html"

# Create an S3 bucket and configure it as a website.
bucket = aws.s3.Bucket("bucket",
    acl="public-read",
    website=aws.s3.BucketWebsiteArgs(
        index_document=index_document,
        error_document=error_document,
    ))

# Create an S3BucketFolder to manage the files of the website.
bucket_folder = synced_folder.S3BucketFolder("bucket-folder",
    path=path,
    bucket_name=bucket.bucket,
    acl="public-read")

# Create a CloudFront CDN to distribute and cache the website.
cdn = aws.cloudfront.Distribution("cdn",
    enabled=True,
    origins=[aws.cloudfront.DistributionOriginArgs(
        origin_id=bucket.arn,
        domain_name=bucket.website_endpoint,
        custom_origin_config=aws.cloudfront.DistributionOriginCustomOriginConfigArgs(
            origin_protocol_policy="http-only",
            http_port=80,
            https_port=443,
            origin_ssl_protocols=["TLSv1.2"],
        ),
    )],
    default_cache_behavior=aws.cloudfront.DistributionDefaultCacheBehaviorArgs(
        target_origin_id=bucket.arn,
        viewer_protocol_policy="redirect-to-https",
        allowed_methods=[
            "GET",
            "HEAD",
            "OPTIONS",
        ],
        cached_methods=[
            "GET",
            "HEAD",
            "OPTIONS",
        ],
        default_ttl=600,
        max_ttl=600,
        min_ttl=0,
        forwarded_values=aws.cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs(
            query_string=True,
            cookies=aws.cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs(
                forward="all",
            ),
        ),
    ),
    price_class="PriceClass_100",
    custom_error_responses=[aws.cloudfront.DistributionCustomErrorResponseArgs(
        error_code=404,
        response_code=404,
        response_page_path=f"/{error_document}",
    )],
    restrictions=aws.cloudfront.DistributionRestrictionsArgs(
        geo_restriction=aws.cloudfront.DistributionRestrictionsGeoRestrictionArgs(
            restriction_type="none",
        ),
    ),
    viewer_certificate=aws.cloudfront.DistributionViewerCertificateArgs(
        cloudfront_default_certificate=True,
    ))

# Export the URLs and hostnames of the bucket and distribution.
pulumi.export("originURL", pulumi.Output.concat("http://", bucket.website_endpoint))
pulumi.export("originHostname", bucket.website_endpoint)
pulumi.export("cdnURL", pulumi.Output.concat("https://", cdn.domain_name))
pulumi.export("cdnHostname", cdn.domain_name)
