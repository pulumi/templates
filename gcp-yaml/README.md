# Pulumi YAML Google Cloud Storage Bucket Template

A minimal Pulumi YAML template for provisioning a Google Cloud Storage Bucket in a Google Cloud project, using **Pulumi YAML** as the runtime.

## Overview

This template deploys a single Google Cloud Storage bucket with default settings, and exports its URL as an output.

Key details:
- **Provider**: `gcp` (Google Cloud Platform)
- **Resource**: `gcp:storage:Bucket` (named `my-bucket`) in the `US` region
- **Output**: `bucketName` – the URL of the created bucket (`${my-bucket.url}`)

## When to Use This Template

- You want to learn Pulumi using declarative YAML
- You need a quick, reproducible Storage bucket
- You prefer minimal, infrastructure-as-code samples without writing TypeScript, Python, or Go

## Prerequisites

- A Google Cloud project (you must know the project ID)
- Credentials for GCP set up locally (`gcloud auth login` or `GOOGLE_APPLICATION_CREDENTIALS` environment variable)
- Pulumi CLI installed and authenticated (see https://www.pulumi.com/docs/get-started)

## Usage

1. Create a new Pulumi project from this template:
   ```bash
   pulumi new gcp-yaml
   ```
2. When prompted:
   - **Project Name**: choose a name for your Pulumi project
   - **Description**: add a short description
   - **gcp:project**: enter your GCP project ID
3. Deploy the stack:
   ```bash
   pulumi up
   ```
4. After deployment, view the bucket URL in the output.

## Project Layout

```text
.
├── Pulumi.yaml         # Defines the project, resources, and outputs in YAML
└── Pulumi.<stack>.yaml # Stack-specific configuration (auto-generated)
```

## Configuration

This template supports the following Pulumi configuration values:

- `gcp:project` (required): Your Google Cloud project ID

Set or override config values with:
```bash
pulumi config set gcp:project YOUR_PROJECT_ID
```

## Next Steps

- Customize the bucket (e.g., change `location`, add `storageClass`, enable versioning)
- Add IAM bindings, lifecycle rules, or notifications
- Explore additional GCP resources in Pulumi: https://www.pulumi.com/registry/packages/gcp
- Integrate deployments into CI/CD pipelines

## Getting Help

- Pulumi Documentation: https://www.pulumi.com/docs
- GCP Provider Reference: https://www.pulumi.com/registry/packages/gcp
- Community Slack: https://www.pulumi.com/community
- GitHub Issues: https://github.com/pulumi/pulumi/issues