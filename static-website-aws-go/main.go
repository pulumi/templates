package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		path := "./site"
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
		bucket, err := s3.NewBucket(ctx, "bucket", &s3.BucketArgs{
			Acl: pulumi.String("public-read"),
			Website: &s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String(indexDocument),
				ErrorDocument: pulumi.String(errorDocument),
			},
		})
		if err != nil {
			return err
		}
		_, err = synced.NewS3BucketFolder(ctx, "bucket-folder", &synced.S3BucketFolderArgs{
			Path:       pulumi.String(path),
			BucketName: bucket.Bucket,
			Acl:        pulumi.String("public-read"),
		})
		if err != nil {
			return err
		}
		cdn, err := cloudfront.NewDistribution(ctx, "cdn", &cloudfront.DistributionArgs{
			Enabled: pulumi.Bool(true),
			Origins: cloudfront.DistributionOriginArray{
				&cloudfront.DistributionOriginArgs{
					OriginId:   bucket.Arn,
					DomainName: bucket.WebsiteEndpoint,
					CustomOriginConfig: &cloudfront.DistributionOriginCustomOriginConfigArgs{
						OriginProtocolPolicy: pulumi.String("http-only"),
						HttpPort:             pulumi.Int(80),
						HttpsPort:            pulumi.Int(443),
						OriginSslProtocols: pulumi.StringArray{
							pulumi.String("TLSv1.2"),
						},
					},
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
				MinTtl:     pulumi.Int(0),
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
				SslSupportMethod:             pulumi.String("sni-only"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("originURL", bucket.WebsiteEndpoint.ApplyT(func(websiteEndpoint string) (string, error) {
			return fmt.Sprintf("http://%v", websiteEndpoint), nil
		}).(pulumi.StringOutput))
		ctx.Export("originHostname", bucket.WebsiteEndpoint)
		ctx.Export("cdnURL", cdn.DomainName.ApplyT(func(domainName string) (string, error) {
			return fmt.Sprintf("https://%v", domainName), nil
		}).(pulumi.StringOutput))
		ctx.Export("cdnHostname", cdn.DomainName)
		return nil
	})
}
