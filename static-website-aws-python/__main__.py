import pulumi
import pulumi_aws as aws
import pulumi_synced_folder as synced_folder

config = pulumi.Config()
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
        ssl_support_method="sni-only",
    ))
pulumi.export("originURL", bucket.website_endpoint.apply(lambda website_endpoint: f"http://{website_endpoint}"))
pulumi.export("originHostname", bucket.website_endpoint)
pulumi.export("cdnURL", cdn.domain_name.apply(lambda domain_name: f"https://{domain_name}"))
pulumi.export("cdnHostname", cdn.domain_name)
