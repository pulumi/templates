# Helm Chart on Kubernetes (Pulumi HCL)

A Pulumi HCL program that installs the NGINX ingress controller Helm chart onto a Kubernetes cluster.

## Overview

The program creates a namespace and installs the `nginx-ingress` Helm chart into it. It targets the cluster in your current kubeconfig context. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Kubernetes (`hashicorp/kubernetes`)
- Helm (`hashicorp/helm`)

## Resources Created

- `kubernetes_namespace` (`ingress`): The namespace for the controller.
- `helm_release` (`ingress`): The NGINX ingress controller chart.

## Outputs

- **name**: The name of the Helm release.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- A Kubernetes cluster and a kubeconfig file (the providers read `~/.kube/config` and use the current context).
- `kubectl` and Helm configured to talk to your cluster.

## Usage

```bash
pulumi new helm-kubernetes-hcl
pulumi up
```

## Configuration

- **k8s_namespace**: The namespace to deploy into. Default: `nginx-ingress`.

## Next Steps

- Set additional chart values to customize the controller.
- Create Ingress resources that use this controller.
- Pin and upgrade the chart version over time.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
