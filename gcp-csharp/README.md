# Minimal Google Cloud C# Pulumi Template

This template provides a minimal Pulumi program in C# for provisioning a Google Cloud Storage bucket. It uses the Pulumi GCP provider to create resources in your GCP project and exports the bucket URL as an output.

## Overview

This template:
- Uses the Pulumi .NET runtime (C#)
- Leverages the `Pulumi.Gcp` provider
- Creates a Google Cloud Storage bucket in the `US` location
- Exports the bucket URL as `bucketName`

Ideal for kicking off GCP infrastructure projects with minimal setup.

## Prerequisites

Ensure you have:
- .NET 8.0 SDK installed
- A Google Cloud project with billing enabled
- Authenticated GCP credentials, e.g. via:
  - `gcloud auth application-default login`
  - or setting `GOOGLE_APPLICATION_CREDENTIALS` to a service account key file
- Pulumi CLI configured and logged in to your preferred backend

## Usage

From your preferred directory, run:
```bash
pulumi new gcp-csharp
```

Follow the prompts:
- **Project name**: a unique name for your deployment
- **Stack name**: an environment name (e.g., `dev`, `prod`)
- **gcp:project**: your Google Cloud project ID

After scaffolding:
```bash
cd <your-project-name>
pulumi up
```

Confirm the preview and watch Pulumi provision your GCS bucket.

## Project Layout

```
.
├── Pulumi.yaml          # Project and template settings
├── Program.cs           # Entry point defining resources
├── <YourProject>.csproj # C# project file with dependencies
├── bin/                 # Build output (ignored by version control)
└── obj/                 # Build artifacts (ignored by version control)
```

## Configuration

| Key           | Description                        | Default |
|---------------|------------------------------------|---------|
| `gcp:project` | Google Cloud project ID to deploy | *none*  |

Set configuration values with:
```bash
pulumi config set gcp:project YOUR_PROJECT_ID
```

## Outputs

After deployment, the following output is available:
- `bucketName` – The HTTPS URL of the created Google Cloud Storage bucket

Retrieve outputs with:
```bash
pulumi stack output bucketName
```

## When to Use This Template

- You need a quick starting point for a GCP project in C#
- You want to learn Pulumi with a minimal example
- You plan to extend the template with additional resources

## Next Steps

- Rename `my-bucket` in `Program.cs` to follow your naming conventions
- Add more GCP resources (e.g., Pub/Sub, Cloud Functions)
- Integrate Pulumi stacks for multi-environment workflows
- Explore `Pulumi.Gcp` provider docs: https://www.pulumi.com/docs/reference/pkg/gcp/

## Getting Help

- Pulumi docs: https://www.pulumi.com/docs/
- GCP provider reference: https://www.pulumi.com/docs/reference/pkg/gcp/
- Community Slack: https://slack.pulumi.com/
- GitHub Issues: https://github.com/pulumi/pulumi/issues