# Google Cloud Storage Bucket (Pulumi HCL)

A minimal Pulumi HCL template that provisions a Google Cloud Storage bucket and exports its URL.

## Overview

This template uses the Pulumi Google Cloud provider to create a single Cloud Storage bucket. Pulumi auto-names the bucket to keep it globally unique. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google Cloud (`pulumi/gcp`)

## Resources Created

- `gcp_storage_bucket` (`my-bucket`): A multi-region (`US`) Cloud Storage bucket.

## Outputs

- **bucket_name**: The URL of the created bucket.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Google Cloud CLI installed and authenticated (`gcloud auth application-default login`).
- A Google Cloud project with the Cloud Storage API enabled and permissions to create buckets.

## Usage

Initialize a new project from this template by running:

```bash
pulumi new gcp-hcl
```

You will be prompted for:
- A project name (default is set by the template).
- A project description.
- The Google Cloud project to deploy into.

After initialization, deploy your stack:

```bash
pulumi up
```

## Project Layout

After `pulumi new`, your directory will look like:

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., google:project)
```

## Configuration

This template supports the following configuration keys:

- **google:project**: The Google Cloud project to deploy resources into.

To set the project, run:

```bash
pulumi config set google:project my-project-id
```

## When to Use This Template

This template is ideal if you need:
- A lightweight starting point for a Cloud Storage bucket.
- To learn Pulumi with HCL programs.
- A quick bootstrap for small storage-focused projects.

## Next Steps

- Enable bucket versioning, lifecycle rules, or uniform bucket-level access.
- Add IAM bindings for access control.
- Integrate with other Google Cloud services (e.g., Cloud Functions, Cloud CDN).
- Explore additional Pulumi HCL examples.

## Getting Help

If you have questions or encounter issues:
- Visit the Pulumi documentation: https://www.pulumi.com/docs/
- Join the Pulumi Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
