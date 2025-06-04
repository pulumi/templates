 # AWS S3 Bucket (Pulumi YAML)

 A minimal Pulumi YAML template that provisions an AWS S3 Bucket and exports its name.

 ## Overview

 This template uses the AWS provider to create a single S3 bucket. It is a great starting point for projects that require simple object storage with minimal setup.

 ## Providers

 - AWS

 ## Resources Created

 - aws:s3:BucketV2 (`my-bucket`): A basic S3 bucket.

 ## Outputs

 - **bucketName**: The name (ID) of the created S3 bucket.

 ## Prerequisites

 - Pulumi CLI configured and logged in to your chosen backend.
 - AWS credentials configured (environment variables, `~/.aws/credentials`, or `AWS_PROFILE`).
 - An AWS account with permissions to create S3 buckets.

 ## Usage

 Initialize a new project from this template by running:

 ```bash
 pulumi new aws-yaml
 ```

 You will be prompted for:
 - A project name (default is set by the template).
 - A project description.
 - The AWS region to deploy into (default: `us-east-1`).

 After initialization, deploy your stack:

 ```bash
 pulumi up
 ```

 ## Project Layout

 After `pulumi new`, your directory will look like:

 ```
 .
 ├── Pulumi.yaml           # Project metadata and YAML program
 └── Pulumi.<stack>.yaml   # Stack configuration (e.g., aws:region)
 ```

 ## Configuration

 This template supports the following configuration keys:

 - **aws:region**: The AWS region to deploy resources into.
   - Default: `us-east-1`

 To override the region, run:

 ```bash
 pulumi config set aws:region us-west-2
 ```

 ## When to Use This Template

 This template is ideal if you need:
 - A lightweight starting point for creating an S3 bucket.
 - To learn Pulumi with YAML programs.
 - A quick bootstrap for small storage-focused projects.

 ## Next Steps

 - Enable bucket versioning, encryption, or lifecycle rules.
 - Add IAM policies or roles for access control.
 - Integrate with other AWS services (e.g., Lambda, CloudFront).
 - Explore additional Pulumi AWS YAML examples.

 ## Getting Help

 If you have questions or encounter issues:
 - Visit the Pulumi documentation: https://www.pulumi.com/docs/
 - Join the Pulumi Community Slack: https://www.pulumi.com/slack
 - Open an issue in this GitHub repository.