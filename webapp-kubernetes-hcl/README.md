# Web Application on Kubernetes (Pulumi HCL)

A Pulumi HCL program that deploys a simple Nginx web application onto a Kubernetes cluster.

## Overview

The program creates a namespace, a ConfigMap holding an Nginx configuration, a Deployment that mounts that configuration, and a Service that exposes the Deployment. It targets the cluster in your current kubeconfig context. The program is written in HCL (`main.tf`) and run by Pulumi's native HCL runtime.

## Providers

- Kubernetes (`hashicorp/kubernetes`)

## Resources Created

- `kubernetes_namespace` (`webserver`): The namespace for the application.
- `kubernetes_config_map` (`config`): The Nginx configuration.
- `kubernetes_deployment_v1` (`webserver`): The Nginx Deployment.
- `kubernetes_service_v1` (`webserver`): A Service exposing the Deployment on port 80.

## Outputs

- **deployment_name**: The name of the Deployment.
- **service_name**: The name of the Service.

## Prerequisites

- Pulumi CLI configured and logged in to your chosen backend.
- A Kubernetes cluster and a kubeconfig file (the provider reads `~/.kube/config` and uses the current context).
- `kubectl` configured to talk to your cluster.

## Usage

```bash
pulumi new webapp-kubernetes-hcl
pulumi up
```

## Configuration

- **k8s_namespace**: The namespace to deploy into. Default: `webapp`.
- **num_replicas**: The number of replicas. Default: `1`.

## Next Steps

- Change the Service type to `LoadBalancer` to expose it externally.
- Add an Ingress resource to route traffic to the Service.
- Replace the Nginx image with your own application image.

## Getting Help

- Pulumi documentation: https://www.pulumi.com/docs/
- Community Slack: https://www.pulumi.com/slack
- Open an issue in this GitHub repository.
