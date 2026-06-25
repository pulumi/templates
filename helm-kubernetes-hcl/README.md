# Helm Chart on Kubernetes (Pulumi HCL)

A Pulumi HCL program that installs the NGINX ingress controller Helm chart onto a Kubernetes cluster.

## Overview

The program creates a namespace and installs the `nginx-ingress` Helm chart into it using the native Kubernetes Helm Release resource (no separate Helm provider). It targets the cluster in your ambient kubeconfig and current context. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Kubernetes (`pulumi/kubernetes`)

## Resources Created

- `kubernetes_core_v1_namespace` (`ingressns`): The namespace for the controller.
- `kubernetes_helm.sh_v3_release` (`ingresscontroller`): The NGINX ingress controller Helm release (resource token `kubernetes:helm.sh/v3:Release`).

## Outputs

- **name**: The name of the Helm release.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- A Kubernetes cluster and a kubeconfig file (the provider uses your ambient `~/.kube/config` and current context).
- `kubectl` configured to talk to your cluster.

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
