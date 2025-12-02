# Minimal GCP Kotlin Pulumi Template

This template provides a minimal Pulumi program written in Kotlin that provisions a Google Cloud Storage bucket using the Pulumi GCP provider. It's a great starting point for building GCP infrastructure with Pulumi and Kotlin.

## Providers
- GCP (pulumi/gcp)

## Resources
- `gcp.storage.Bucket`: A Google Cloud Storage bucket resource.

## Outputs
- `bucketName`: The name of the created GCS bucket.

## When to use this template
Use this template if you:
- Want a quick start with Pulumi in Kotlin
- Need an example of provisioning basic GCP resources
- Are familiar with Gradle and Kotlin development

## Prerequisites
- Java Development Kit (JDK) 21 or later
- GCP credentials configured (via `gcloud auth application-default login`, service account key, or workload identity)

## Getting Started
1. Create a new Pulumi project from this template:
```bash
   pulumi new gcp-kotlin
```
2. Follow the interactive prompts to set your project name, description, and GCP project ID and region (default: `us-central1`).
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
- `gcp:project` (string) — GCP project ID to deploy into. Required.
- `gcp:region` (string) — GCP region to deploy into. Default: `us-central1`.

View or set configuration values:
```bash
pulumi config
pulumi config set gcp:project my-gcp-project-id
pulumi config set gcp:region us-west1
```

## Getting Help
- Pulumi documentation: https://www.pulumi.com/docs/
- GCP provider reference: https://www.pulumi.com/registry/packages/gcp/
- Pulumi Community Slack: https://slack.pulumi.com/
- Stack Overflow (`pulumi` tag)
- Report issues: https://github.com/pulumi/pulumi/issues