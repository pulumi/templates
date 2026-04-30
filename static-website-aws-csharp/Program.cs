using System.Collections.Generic;
using System.Text.Json;
using Pulumi;
using Aws = Pulumi.Aws;
using SyncedFolder = Pulumi.SyncedFolder;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var path = config.Get("path") ?? "./www";
    var indexDocument = config.Get("indexDocument") ?? "index.html";
    var errorDocument = config.Get("errorDocument") ?? "error.html";

    // Create a private S3 bucket to hold the website content.
    var bucket = new Aws.S3.Bucket("bucket");

    // Block all public access to the bucket; CloudFront will reach it via OAC.
    var publicAccessBlock = new Aws.S3.BucketPublicAccessBlock("public-access-block", new()
    {
        Bucket = bucket.Id,
        BlockPublicAcls = true,
        BlockPublicPolicy = true,
        IgnorePublicAcls = true,
        RestrictPublicBuckets = true,
    });

    // Sync the website files to the bucket as private objects.
    var bucketFolder = new SyncedFolder.S3BucketFolder("bucket-folder", new()
    {
        Path = path,
        BucketName = bucket.BucketName,
        Acl = "private",
    }, new ComponentResourceOptions {
        DependsOn = { publicAccessBlock },
    });

    // Create an Origin Access Control so CloudFront can read from the private bucket.
    var originAccessControl = new Aws.CloudFront.OriginAccessControl("origin-access-control", new()
    {
        OriginAccessControlOriginType = "s3",
        SigningBehavior = "always",
        SigningProtocol = "sigv4",
    });

    // Create a CloudFront CDN to distribute and cache the website.
    var cdn = new Aws.CloudFront.Distribution("cdn", new()
    {
        Enabled = true,
        DefaultRootObject = indexDocument,
        Origins = new[]
        {
            new Aws.CloudFront.Inputs.DistributionOriginArgs
            {
                OriginId = bucket.Arn,
                DomainName = bucket.BucketRegionalDomainName,
                OriginAccessControlId = originAccessControl.Id,
            },
        },
        DefaultCacheBehavior = new Aws.CloudFront.Inputs.DistributionDefaultCacheBehaviorArgs
        {
            TargetOriginId = bucket.Arn,
            ViewerProtocolPolicy = "redirect-to-https",
            AllowedMethods = new[] { "GET", "HEAD", "OPTIONS" },
            CachedMethods = new[] { "GET", "HEAD", "OPTIONS" },
            Compress = true,
            // Managed-CachingOptimized
            CachePolicyId = "658327ea-f89d-4fab-a63d-7e88639e58f6",
            // Managed-SecurityHeadersPolicy
            ResponseHeadersPolicyId = "67f7725c-6f97-4210-82d7-5512b31e9d03",
        },
        PriceClass = "PriceClass_100",
        CustomErrorResponses = new[]
        {
            new Aws.CloudFront.Inputs.DistributionCustomErrorResponseArgs
            {
                ErrorCode = 404,
                ResponseCode = 404,
                ResponsePagePath = $"/{errorDocument}",
            },
        },
        Restrictions = new Aws.CloudFront.Inputs.DistributionRestrictionsArgs
        {
            GeoRestriction = new Aws.CloudFront.Inputs.DistributionRestrictionsGeoRestrictionArgs
            {
                RestrictionType = "none",
            },
        },
        ViewerCertificate = new Aws.CloudFront.Inputs.DistributionViewerCertificateArgs
        {
            CloudfrontDefaultCertificate = true,
        },
    });

    // Grant the CloudFront distribution permission to read objects from the bucket.
    var bucketPolicy = new Aws.S3.BucketPolicy("bucket-policy", new()
    {
        Bucket = bucket.Id,
        Policy = Output.Tuple(bucket.Arn, cdn.Arn).Apply(arns =>
        {
            var (bucketArn, cdnArn) = arns;
            return JsonSerializer.Serialize(new
            {
                Version = "2012-10-17",
                Statement = new[]
                {
                    new
                    {
                        Sid = "AllowCloudFrontServicePrincipalReadOnly",
                        Effect = "Allow",
                        Principal = new { Service = "cloudfront.amazonaws.com" },
                        Action = "s3:GetObject",
                        Resource = $"{bucketArn}/*",
                        Condition = new
                        {
                            StringEquals = new Dictionary<string, string>
                            {
                                { "AWS:SourceArn", cdnArn },
                            },
                        },
                    },
                },
            });
        }),
    });

    // Export the URLs and hostnames of the bucket and distribution.
    return new Dictionary<string, object?>
    {
        ["originHostname"] = bucket.BucketRegionalDomainName,
        ["cdnURL"] = Output.Format($"https://{cdn.DomainName}"),
        ["cdnHostname"] = cdn.DomainName,
    };
});
