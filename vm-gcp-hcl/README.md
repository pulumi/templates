# Virtual Machine on Google Cloud (Pulumi HCL)

A Pulumi HCL program that deploys a Google Compute Engine virtual machine running a simple web server.

## Overview

The program creates a VPC network and subnet, a firewall allowing HTTP and SSH, and a Compute Engine instance with an ephemeral public IP. A startup script serves a "Hello, world!" page. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Google Cloud (`pulumi/gcp`)

## Resources Created

- `gcp_compute_network` (`network`) / `gcp_compute_subnetwork` (`subnet`): The network.
- `gcp_compute_firewall` (`firewall`): Allows inbound SSH and HTTP to tagged instances.
- `gcp_compute_instance` (`instance`): The VM running the web server.

## Outputs

- **name**: The instance name.
- **ip**: The instance's public IP address.
- **url**: The HTTP URL of the web server.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- The Google Cloud CLI authenticated (`gcloud auth application-default login`).
- A Google Cloud project with the Compute Engine API enabled.

## Usage

```bash
pulumi new vm-gcp-hcl
pulumi up
```

Open the `url` output once the instance has booted.

## Project Layout

```
.
├── Pulumi.yaml           # Project metadata and HCL runtime configuration
├── main.tf               # HCL program
└── Pulumi.<stack>.yaml   # Stack configuration (e.g., google:project, google:zone)
```

## Configuration

- **google:project**: The Google Cloud project to deploy into.
- **google:region** / **google:zone**: The region and zone. Defaults: `us-central1` / `us-central1-a`.
- **machine_type**: The machine type. Default: `e2-micro`.
- **os_image**: The OS image. Default: `debian-11`.
- **instance_tag**: The network tag for the firewall. Default: `webserver`.
- **service_port**: The HTTP port to serve on. Default: `80`.

## Next Steps

- Replace the inline startup script with your own application.
- Add an SSH key via instance metadata for access.
- Put the instance behind a global HTTP load balancer.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
