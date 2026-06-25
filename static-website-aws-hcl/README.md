# Static Website on AWS (Pulumi HCL)

A Pulumi HCL program that deploys a static website to AWS using a private S3 bucket fronted by a CloudFront CDN.

## Overview

The website content in `./www` is uploaded to a private S3 bucket. A CloudFront distribution serves the content over HTTPS and reads from the bucket through an Origin Access Control (OAC), so the bucket itself stays private. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`pulumi/aws`)
- Synced Folder (`pulumi/synced-folder`) — uploads the website folder to the bucket

## Resources Created

- `aws_s3_bucket` (`bucket`): A private bucket holding the website content.
- `aws_s3_bucket_public_access_block` (`public-access-block`): Blocks all public access to the bucket.
- `synced-folder_s3_bucket_folder` (`bucket-folder`): Syncs the contents of `path` to the bucket as private objects.
- `aws_cloudfront_origin_access_control` (`origin-access-control`): Lets CloudFront read from the private bucket.
- `aws_cloudfront_distribution` (`cdn`): The CDN that serves and caches the site.
- `aws_s3_bucket_policy` (`bucket-policy`): Grants the distribution read access to the bucket.

## Outputs

- **cdn_url**: The HTTPS URL of the CloudFront distribution.
- **cdn_hostname**: The CloudFront distribution hostname.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- AWS credentials configured (environment variables, `~/.aws/credentials`, or `AWS_PROFILE`).
- An AWS account with permissions to create S3 and CloudFront resources.

## Usage

Initialize a new project from this template by running:

```bash
pulumi new static-website-aws-hcl
```

You will be prompted for a project name, description, region, and the website folder and document settings.

After initialization, deploy your stack:

```bash
pulumi up
```

Then open the `cdnURL` output in your browser. (A new CloudFront distribution can take several minutes to deploy.)

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── www/                  # Website content (index.html, error.html)
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., aws:region)
```

## Configuration

- **aws:region**: The AWS region to deploy into. Default: `us-west-2`.
- **path**: The folder containing the website content. Default: `./www`.
- **index_document**: The top-level page. Default: `index.html`.
- **error_document**: The error page. Default: `error.html`.

```bash
pulumi config set aws:region us-east-1
```

## Next Steps

- Add a custom domain and ACM certificate to the distribution.
- Add cache behaviors or a logging configuration.
- Wire up a build step to generate the contents of `www`.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
