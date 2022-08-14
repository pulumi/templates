package myproject

import com.pulumi.kubernetes.apps_v1.Deployment
import com.pulumi.kubernetes.apps_v1.DeploymentArgs
import com.pulumi.kubernetes.apps_v1.inputs.DeploymentSpecArgs
import com.pulumi.kubernetes.core_v1.inputs.ContainerArgs
import com.pulumi.kubernetes.core_v1.inputs.ContainerPortArgs
import com.pulumi.kubernetes.core_v1.inputs.PodSpecArgs
import com.pulumi.kubernetes.core_v1.inputs.PodTemplateSpecArgs
import com.pulumi.kubernetes.meta_v1.inputs.LabelSelectorArgs
import com.pulumi.kubernetes.meta_v1.inputs.ObjectMetaArgs

fun getContainers(): PodSpecArgs.Builder = PodSpecArgs.builder().containers(
    ContainerArgs.builder()
        .name("nginx")
        .image("nginx")
        .ports(
            ContainerPortArgs.builder()
                .containerPort(80)
                .build()
        )
        .build()
)

fun getTemplate(labels: Map<String, String>): PodTemplateSpecArgs =
    PodTemplateSpecArgs.builder()
        .metadata(
            ObjectMetaArgs.builder()
                .labels(labels)
                .build()
        )
        .spec(
            getContainers()
                .build()
        ).build()


fun getDeployment(labels: Map<String, String>): Deployment =
    Deployment(
        "nginx", DeploymentArgs.builder()
            .spec(
                DeploymentSpecArgs.builder()
                    .selector(
                        LabelSelectorArgs.builder()
                            .matchLabels(labels)
                            .build()
                    )
                    .replicas(1)
                    .template(getTemplate(labels))
                    .build()
            )
            .build()
    )


