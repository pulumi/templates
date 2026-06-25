# Kubernetes nginx Deployment (Pulumi HCL)

A minimal Pulumi HCL template that deploys an nginx Deployment to a Kubernetes cluster and exports its name.

## Overview

This template uses the Pulumi Kubernetes provider to create a single-replica nginx Deployment in your currently configured cluster. The provider uses the ambient kubeconfig and current context. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Kubernetes (`pulumi/kubernetes`)

## Resources Created

- `kubernetes_apps_v1_deployment` (`deployment`): A single-replica Deployment running the `nginx` image.

## Outputs

- **name**: The name of the created Deployment.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- A Kubernetes cluster and a kubeconfig file (the template reads `~/.kube/config` and uses the current context).
- `kubectl` configured to talk to your cluster.

## Usage

Initialize a new project from this template by running:

```bash
pulumi new kubernetes-hcl
```

You will be prompted for:
- A project name (default is set by the template).
- A project description.

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
└── Pulumi.<stack>.yaml   # Stack configuration
```

## Configuration

The native Kubernetes provider uses your ambient kubeconfig (`~/.kube/config`) and current context. To target a different cluster, switch your `kubectl` context.

## When to Use This Template

This template is ideal if you need:
- A lightweight starting point for deploying workloads to Kubernetes.
- To learn Pulumi with HCL programs.
- A quick bootstrap for a containerized application.

## Next Steps

- Expose the Deployment with a Service or Ingress.
- Add a ConfigMap, Secret, or persistent volume.
- Scale the Deployment by increasing `replicas`.
- Explore additional Pulumi HCL examples.

## Getting Help

If you have questions or encounter issues:
- Visit the Pulumi documentation: https://www.pulumi.com/docs/
- Join the Pulumi Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
