# Serverless Application on Google Cloud (Pulumi HCL)

A Pulumi HCL program that deploys a serverless application on Google Cloud: a Cloud Function (Gen 2) with a static front-end hosted on Cloud Storage.

## Overview

A Python Cloud Function returns the current time. A static website in `./www` is hosted from a public Cloud Storage bucket and calls the function (the function URL is injected into a `config.json` the page reads at load time). The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google Cloud (`pulumi/gcp`)
- Synced Folder (`pulumi/synced-folder`) — uploads the website folder to the bucket

## Resources Created

- `gcp_storage_bucket` (`site-bucket`) + `gcp_storage_bucket_i_a_m_binding` (`site-bucket-iam-binding`): The website bucket and its public read access.
- `synced-folder_google_cloud_folder` (`synced-folder`): Syncs the website content to the bucket.
- `gcp_storage_bucket` (`app-bucket`) + `gcp_storage_bucket_object` (`app-archive`): Holds the function's source archive.
- `gcp_cloudfunctionsv2_function` (`data-function`): The Gen 2 Cloud Function, built from `./app`.
- `gcp_cloudrun_iam_member` (`invoker`): Allows public invocation of the function's underlying Cloud Run service.
- `gcp_storage_bucket_object` (`site-config`): The website's `config.json` pointing at the function URL.

## Outputs

- **site_url**: The URL of the static website.
- **api_url**: The URL of the function endpoint.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Google Cloud CLI authenticated (`gcloud auth application-default login`).
- A Google Cloud project with the Cloud Functions, Cloud Run, Cloud Build, Artifact Registry, and Cloud Storage APIs enabled.

## Usage

```bash
pulumi new serverless-gcp-hcl
pulumi up
```

Open the `site_url` output and click the button.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── app/                  # The Cloud Function source (main.py, requirements.txt)
├── www/                  # Static front-end
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., google:project)
```

## Configuration

- **google:project**: The Google Cloud project to deploy into.
- **region**: The region for the function. Default: `us-central1`.
- **site_path** / **app_path**: The website and function source folders.
- **index_document** / **error_document**: The website's page documents.

## Next Steps

- Add more functions and routes.
- Put an HTTPS load balancer and Cloud CDN in front of the site bucket.
- Add a custom domain.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
