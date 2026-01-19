# Minimal AWS Kotlin Pulumi Template

This template provides a minimal Pulumi program written in Kotlin that provisions an Amazon S3 bucket using the Pulumi AWS provider. It’s a great starting point for building AWS infrastructure with Pulumi and kotlin.

## Providers
- AWS (pulumi/aws)

## Resources
- `aws.s3.Bucket`: An S3 bucket resource.

## Outputs
- `bucketName`: The name of the created S3 bucket.

## When to use this template
Use this template if you:
- Want a quick start with Pulumi in Kotlin
- Need an example of provisioning basic AWS resources
- Are familiar with Maven and Kotlin development

## Prerequisites
- Java Development Kit (JDK) 21 or later
- AWS credentials configured (via AWS CLI, environment variables, or shared credentials file)

## Getting Started
1. Create a new Pulumi project from this template:
   ```bash
   pulumi new aws-kotlin
   ```
2. Follow the interactive prompts to set your project name, description, and AWS region (default: `us-east-1`).
3. Change into your project directory:
   ```bash
   cd <project-name>
   ```
4. Deploy your stack:
   ```bash
   pulumi up
   ```

## Configuration
This template supports the following configuration values:
- `aws:region` (string) — AWS region to deploy into. Default: `us-east-1`.

View or set configuration values:
```bash
pulumi config
pulumi config set aws:region us-west-2
```

## Getting Help
- Pulumi documentation: https://www.pulumi.com/docs/
- AWS provider reference: https://www.pulumi.com/docs/reference/pkg/aws/
- Pulumi Community Slack: https://slack.pulumi.com/
- Stack Overflow (`pulumi` tag)
- Report issues: https://github.com/pulumi/pulumi/issues