import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";
import * as synced_folder from "@pulumi/synced-folder";

// Import the program's configuration settings.
const config = new pulumi.Config();
const path = config.get("path") || "www";
const indexDocument = config.get("indexDocument") || "index.html";
const errorDocument = config.get("errorDocument") || "error.html";

// Create a private S3 bucket to hold the website content.
const bucket = new aws.s3.Bucket("bucket");

// Block all public access to the bucket; CloudFront will reach it via OAC.
const publicAccessBlock = new aws.s3.BucketPublicAccessBlock("public-access-block", {
    bucket: bucket.bucket,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true,
});

// Sync the website files to the bucket as private objects.
const bucketFolder = new synced_folder.S3BucketFolder("bucket-folder", {
    path: path,
    bucketName: bucket.bucket,
    acl: "private",
}, { dependsOn: [publicAccessBlock] });

// Create an Origin Access Control so CloudFront can read from the private bucket.
const originAccessControl = new aws.cloudfront.OriginAccessControl("origin-access-control", {
    originAccessControlOriginType: "s3",
    signingBehavior: "always",
    signingProtocol: "sigv4",
});

// Create a CloudFront CDN to distribute and cache the website.
const cdn = new aws.cloudfront.Distribution("cdn", {
    enabled: true,
    defaultRootObject: indexDocument,
    origins: [{
        originId: bucket.arn,
        domainName: bucket.bucketRegionalDomainName,
        originAccessControlId: originAccessControl.id,
    }],
    defaultCacheBehavior: {
        targetOriginId: bucket.arn,
        viewerProtocolPolicy: "redirect-to-https",
        allowedMethods: [
            "GET",
            "HEAD",
            "OPTIONS",
        ],
        cachedMethods: [
            "GET",
            "HEAD",
            "OPTIONS",
        ],
        defaultTtl: 600,
        maxTtl: 600,
        minTtl: 600,
        forwardedValues: {
            queryString: true,
            cookies: {
                forward: "all",
            },
        },
    },
    priceClass: "PriceClass_100",
    customErrorResponses: [{
        errorCode: 404,
        responseCode: 404,
        responsePagePath: `/${errorDocument}`,
    }],
    restrictions: {
        geoRestriction: {
            restrictionType: "none",
        },
    },
    viewerCertificate: {
        cloudfrontDefaultCertificate: true,
    },
});

// Grant the CloudFront distribution permission to read objects from the bucket.
const bucketPolicy = new aws.s3.BucketPolicy("bucket-policy", {
    bucket: bucket.bucket,
    policy: pulumi.jsonStringify({
        Version: "2012-10-17",
        Statement: [{
            Sid: "AllowCloudFrontServicePrincipalReadOnly",
            Effect: "Allow",
            Principal: { Service: "cloudfront.amazonaws.com" },
            Action: "s3:GetObject",
            Resource: pulumi.interpolate`${bucket.arn}/*`,
            Condition: {
                StringEquals: { "AWS:SourceArn": cdn.arn },
            },
        }],
    }),
});

// Export the URL and hostname of the CloudFront distribution.
export const cdnURL = pulumi.interpolate`https://${cdn.domainName}`;
export const cdnHostname = cdn.domainName;
