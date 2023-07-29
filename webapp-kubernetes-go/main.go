package main

import (
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Get some configuration values from the Pulumi stack, or use defaults
		cfg := config.New(ctx, "")
		k8sNamespace, err := cfg.Try("namespace")
		if err != nil {
			k8sNamespace = "default"
		}
		numReplicas, err := cfg.TryInt("replicas")
		if err != nil {
			numReplicas = 1
		}
		appLabels := pulumi.StringMap{
			"app": pulumi.String("nginx"),
		}

		// Create a new namespace
		webServerNs, err := corev1.NewNamespace(ctx, "webserverns", &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.String(k8sNamespace),
			},
		})
		if err != nil {
			return err
		}

		// Create a new ConfigMap for the Nginx configuration
		webServerConfig, err := corev1.NewConfigMap(ctx, "webserverconfig", &corev1.ConfigMapArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: webServerNs.Metadata.Name(),
			},
			Data: pulumi.StringMap{
				"nginx.conf": pulumi.Sprintf(`events { }
http {
  server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html index.htm index.nginx-debian.html
    server_name _;
    location / {
      try_files $uri $uri/ =404;
    }
  }
}
`),
			},
		})
		if err != nil {
			return err
		}

		// Create a new Deployment with a user-specified number of replicas
		webServerDeployment, err := appsv1.NewDeployment(ctx, "webserverdeployment", &appsv1.DeploymentArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: webServerNs.Metadata.Name(),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Selector: &metav1.LabelSelectorArgs{
					MatchLabels: pulumi.StringMap(appLabels),
				},
				Replicas: pulumi.Int(numReplicas),
				Template: &corev1.PodTemplateSpecArgs{
					Metadata: &metav1.ObjectMetaArgs{
						Labels: pulumi.StringMap(appLabels),
					},
					Spec: &corev1.PodSpecArgs{
						Containers: corev1.ContainerArray{
							&corev1.ContainerArgs{
								Image: pulumi.String("nginx"),
								Name:  pulumi.String("nginx"),
								VolumeMounts: corev1.VolumeMountArray{
									&corev1.VolumeMountArgs{
										MountPath: pulumi.String("/etc/nginx/nginx.conf"),
										Name:      pulumi.String("nginx-conf-volume"),
										ReadOnly:  pulumi.Bool(true),
										SubPath:   pulumi.String("nginx.conf"),
									},
								},
							},
						},
						Volumes: corev1.VolumeArray{
							&corev1.VolumeArgs{
								ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
									Items: corev1.KeyToPathArray{
										&corev1.KeyToPathArgs{
											Key:  pulumi.String("nginx.conf"),
											Path: pulumi.String("nginx.conf"),
										},
									},
									Name: webServerConfig.Metadata.Name(),
								},
								Name: pulumi.String("nginx-conf-volume"),
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Expose the Deployment as a Kubernetes Service
		webServerService, err := corev1.NewService(ctx, "webserverservice", &corev1.ServiceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Namespace: webServerNs.Metadata.Name(),
			},
			Spec: &corev1.ServiceSpecArgs{
				Ports: corev1.ServicePortArray{
					&corev1.ServicePortArgs{
						Port:       pulumi.Int(80),
						TargetPort: pulumi.Any(80),
						Protocol:   pulumi.String("TCP"),
					},
				},
				Selector: pulumi.StringMap(appLabels),
			},
		})
		if err != nil {
			return err
		}

		// Export some values for use elsewhere
		ctx.Export("deploymentName", webServerDeployment.Metadata.Name())
		ctx.Export("serviceName", webServerService.Metadata.Name())

		return nil
	})
}
