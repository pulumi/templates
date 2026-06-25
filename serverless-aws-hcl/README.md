# Serverless Application on AWS (Pulumi HCL)

A Pulumi HCL program that deploys a serverless application on AWS: a Lambda function behind an API Gateway REST API that also serves a static front-end.

## Overview

A Python Lambda function returns the current time. An API Gateway REST API (built with the `aws-apigateway` component) serves the static front-end in `./www` at the root path and routes `GET /date` to the function. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- AWS (`pulumi/aws`)
- API Gateway (`pulumi/aws-apigateway`) — the REST API component
- Archive (`hashicorp/archive`) — packages the function source into a zip

## Resources Created

- `data archive_file` (`fn`): Packages the `./function` source into a deployment archive.
- `aws_iam_role` (`role`): The Lambda execution role.
- `aws_lambda_function` (`fn`): The function, packaged from `./function`.
- `aws-apigateway_rest_a_p_i` (`api`): A REST API that serves `./www` at `/` and routes `GET /date` to the function.

## Outputs

- **url**: The URL at which the REST API (and static front-end) is served.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- AWS credentials configured (environment variables, `~/.aws/credentials`, or `AWS_PROFILE`).
- An AWS account with permissions for Lambda, API Gateway, IAM, and S3.

## Usage

```bash
pulumi new serverless-aws-hcl
pulumi up
```

Open the `url` output and click the button. (The API route can take a few seconds to become live after the first deploy.)

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

The static front-end is read from `./www` and the function source from `./function`; both paths are set in `main.tf`.

## Next Steps

- Add more routes and functions to the HTTP API.
- Put a CloudFront distribution in front of the S3 site.
- Add a custom domain to the API.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
