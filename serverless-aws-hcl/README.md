# Serverless Application on AWS (Pulumi HCL)

A Pulumi HCL program that deploys a serverless application on AWS: a Lambda function behind an API Gateway HTTP API, with a static front-end hosted on S3.

## Overview

A Python Lambda function returns the current time. An API Gateway HTTP API routes `GET /date` to the function. A static website in `./www` is hosted from an S3 bucket and calls the API (the API endpoint is injected into a `config.json` the page reads at load time). The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`hashicorp/aws`)
- Archive (`hashicorp/archive`) — packages the function source
- Random (`hashicorp/random`)

## Resources Created

- `aws_iam_role` / `aws_iam_role_policy_attachment`: The Lambda execution role.
- `aws_lambda_function` (`fn`): The function, packaged from `./function`.
- `aws_apigatewayv2_api` / `_integration` / `_route` / `_stage`: The HTTP API and its `GET /date` route.
- `aws_lambda_permission` (`apigw`): Lets the API invoke the function.
- `aws_s3_bucket` + website/public-access/policy/objects: The static site and its `config.json`.

## Outputs

- **site_url**: The URL of the static website.
- **api_url**: The URL of the `GET /date` endpoint.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- AWS credentials configured (environment variables, `~/.aws/credentials`, or `AWS_PROFILE`).
- An AWS account with permissions for Lambda, API Gateway, IAM, and S3.

## Usage

```bash
pulumi new serverless-aws-hcl
pulumi up
```

Open the `siteURL` output and click the button. (The API route can take a few seconds to become live after the first deploy.)

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
├── function/handler.py   # The Lambda function source
├── www/                  # Static front-end
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., aws:region)
```

## Configuration

- **aws:region**: The AWS region to deploy into. Default: `us-west-2`.
- **site_path**: The website folder. Default: `./www`.
- **app_path**: The function source folder. Default: `./function`.

## Next Steps

- Add more routes and functions to the HTTP API.
- Put a CloudFront distribution in front of the S3 site.
- Add a custom domain to the API.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
