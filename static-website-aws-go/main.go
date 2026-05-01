package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		path := "www"
		if param := cfg.Get("path"); param != "" {
			path = param
		}
		indexDocument := "index.html"
		if param := cfg.Get("indexDocument"); param != "" {
			indexDocument = param
		}
		errorDocument := "error.html"
		if param := cfg.Get("errorDocument"); param != "" {
			errorDocument = param
		}

		// Create a private S3 bucket to hold the website content.
		bucket, err := s3.NewBucket(ctx, "bucket", nil)
		if err != nil {
			return err
		}

		// Block all public access to the bucket; CloudFront will reach it via OAC.
		publicAccessBlock, err := s3.NewBucketPublicAccessBlock(ctx, "public-access-block", &s3.BucketPublicAccessBlockArgs{
			Bucket:                bucket.Bucket,
			BlockPublicAcls:       pulumi.Bool(true),
			BlockPublicPolicy:     pulumi.Bool(true),
			IgnorePublicAcls:      pulumi.Bool(true),
			RestrictPublicBuckets: pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		// Sync the website files to the bucket as private objects.
		_, err = synced.NewS3BucketFolder(ctx, "bucket-folder", &synced.S3BucketFolderArgs{
			Path:       pulumi.String(path),
			BucketName: bucket.Bucket,
			Acl:        pulumi.String("private"),
		}, pulumi.DependsOn([]pulumi.Resource{publicAccessBlock}))
		if err != nil {
			return err
		}

		// Create an Origin Access Control so CloudFront can read from the private bucket.
		originAccessControl, err := cloudfront.NewOriginAccessControl(ctx, "origin-access-control", &cloudfront.OriginAccessControlArgs{
			OriginAccessControlOriginType: pulumi.String("s3"),
			SigningBehavior:               pulumi.String("always"),
			SigningProtocol:               pulumi.String("sigv4"),
		})
		if err != nil {
			return err
		}

		// Create a CloudFront CDN to distribute and cache the website.
		cdn, err := cloudfront.NewDistribution(ctx, "cdn", &cloudfront.DistributionArgs{
			Enabled:           pulumi.Bool(true),
			DefaultRootObject: pulumi.String(indexDocument),
			Origins: cloudfront.DistributionOriginArray{
				&cloudfront.DistributionOriginArgs{
					OriginId:              bucket.Arn,
					DomainName:            bucket.BucketRegionalDomainName,
					OriginAccessControlId: originAccessControl.ID(),
				},
			},
			DefaultCacheBehavior: &cloudfront.DistributionDefaultCacheBehaviorArgs{
				TargetOriginId:       bucket.Arn,
				ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
				AllowedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
					pulumi.String("OPTIONS"),
				},
				CachedMethods: pulumi.StringArray{
					pulumi.String("GET"),
					pulumi.String("HEAD"),
					pulumi.String("OPTIONS"),
				},
				DefaultTtl: pulumi.Int(600),
				MaxTtl:     pulumi.Int(600),
				MinTtl:     pulumi.Int(600),
				ForwardedValues: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
					QueryString: pulumi.Bool(true),
					Cookies: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
						Forward: pulumi.String("all"),
					},
				},
			},
			PriceClass: pulumi.String("PriceClass_100"),
			CustomErrorResponses: cloudfront.DistributionCustomErrorResponseArray{
				&cloudfront.DistributionCustomErrorResponseArgs{
					ErrorCode:        pulumi.Int(404),
					ResponseCode:     pulumi.Int(404),
					ResponsePagePath: pulumi.String(fmt.Sprintf("/%v", errorDocument)),
				},
			},
			Restrictions: &cloudfront.DistributionRestrictionsArgs{
				GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
					RestrictionType: pulumi.String("none"),
				},
			},
			ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
				CloudfrontDefaultCertificate: pulumi.Bool(true),
			},
		})
		if err != nil {
			return err
		}

		// Grant the CloudFront distribution permission to read objects from the bucket.
		_, err = s3.NewBucketPolicy(ctx, "bucket-policy", &s3.BucketPolicyArgs{
			Bucket: bucket.Bucket,
			Policy: pulumi.All(bucket.Arn, cdn.Arn).ApplyT(func(args []interface{}) (string, error) {
				bucketArn := args[0].(string)
				cdnArn := args[1].(string)
				doc := map[string]interface{}{
					"Version": "2012-10-17",
					"Statement": []interface{}{
						map[string]interface{}{
							"Sid":      "AllowCloudFrontServicePrincipalReadOnly",
							"Effect":   "Allow",
							"Principal": map[string]interface{}{"Service": "cloudfront.amazonaws.com"},
							"Action":   "s3:GetObject",
							"Resource": fmt.Sprintf("%s/*", bucketArn),
							"Condition": map[string]interface{}{
								"StringEquals": map[string]interface{}{"AWS:SourceArn": cdnArn},
							},
						},
					},
				}
				b, err := json.Marshal(doc)
				if err != nil {
					return "", err
				}
				return string(b), nil
			}).(pulumi.StringOutput),
		})
		if err != nil {
			return err
		}

		// Export the URL and hostname of the CloudFront distribution.
		ctx.Export("cdnURL", pulumi.Sprintf("https://%s", cdn.DomainName))
		ctx.Export("cdnHostname", cdn.DomainName)
		return nil
	})
}
