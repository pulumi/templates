# Static Website on Google Cloud (Pulumi HCL)

A Pulumi HCL program that deploys a static website to Google Cloud using a Cloud Storage bucket fronted by a Cloud CDN.

## Overview

The website content in `./www` is uploaded to a public Cloud Storage bucket configured for website hosting. A global HTTP load balancer with Cloud CDN enabled serves and caches the content. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google (`hashicorp/google`)
- Random (`hashicorp/random`)

## Resources Created

- `random_string` (`suffix`): A random suffix used to build globally unique resource names.
- `google_storage_bucket` (`bucket`): A website-enabled Cloud Storage bucket.
- `google_storage_bucket_iam_member` (`public_read`): Grants public read access to objects.
- `google_storage_bucket_object` (`files`): One object per file under `path`.
- `google_compute_backend_bucket` (`backend`): A CDN-enabled backend for the bucket.
- `google_compute_global_address` (`ip`): A global IP address for the CDN.
- `google_compute_url_map` / `google_compute_target_http_proxy` / `google_compute_global_forwarding_rule`: Route requests to the backend bucket.

## Outputs

- **originURL** / **originHostname**: The direct Cloud Storage URL and hostname.
- **cdnURL** / **cdnHostname**: The CDN URL and IP address.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Google Cloud CLI authenticated (`gcloud auth application-default login`).
- A Google Cloud project with the Compute Engine and Cloud Storage APIs enabled.

## Usage

```bash
pulumi new static-website-gcp-hcl
pulumi up
```

A new global forwarding rule and IP can take a few minutes to become reachable.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── www/                  # Website content (index.html, error.html)
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., google:project)
```

## Configuration

- **google:project**: The Google Cloud project to deploy into.
- **path**: The folder containing the website content. Default: `./www`.
- **index_document**: The top-level page. Default: `index.html`.
- **error_document**: The error page. Default: `error.html`.

```bash
pulumi config set google:project my-project-id
```

## Next Steps

- Add an HTTPS forwarding rule with a managed SSL certificate and custom domain.
- Tune the Cloud CDN cache policy on the backend bucket.
- Wire up a build step to generate the contents of `www`.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
