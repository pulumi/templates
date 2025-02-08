package myproject;

import java.util.List;

import com.pulumi.Pulumi;
import com.pulumi.aws.ecs.Cluster;
import com.pulumi.awsx.ecr.Image;
import com.pulumi.awsx.ecr.ImageArgs;
import com.pulumi.awsx.ecr.Repository;
import com.pulumi.awsx.ecr.RepositoryArgs;
import com.pulumi.awsx.ecs.FargateService;
import com.pulumi.awsx.ecs.FargateServiceArgs;
import com.pulumi.awsx.ecs.inputs.FargateServiceTaskDefinitionArgs;
import com.pulumi.awsx.ecs.inputs.TaskDefinitionContainerDefinitionArgs;
import com.pulumi.awsx.ecs.inputs.TaskDefinitionPortMappingArgs;
import com.pulumi.awsx.lb.ApplicationLoadBalancer;

public class App {
    public static void main(String[] args) {
        Pulumi.run(ctx -> {
            var config = ctx.config();
            var containerPort = config.getInteger("containerPort").orElse(80);
            var cpu = config.getInteger("cpu").orElse(512);
            var memory = config.getInteger("memory").orElse(128);

            // An ECS cluster to deploy into
            var cluster = new Cluster("cluster");

            // An ALB to serve the container endpoint to the internet
            var loadbalancer = new ApplicationLoadBalancer("loadbalancer");

            // An ECR repository to store our application's container image
            var repo = new Repository("repo", RepositoryArgs.builder()
                    .forceDelete(true)
                    .build());

            // Build and publish our application's container image from ./app to the ECR
            // repository
            var image = new Image("image", ImageArgs.builder()
                    .repositoryUrl(repo.url())
                    .context("./app")
                    .platform("linux/amd64")
                    .build());

            // Deploy an ECS Service on Fargate to host the application container
            var service = new FargateService("service", FargateServiceArgs.builder()
                    .cluster(cluster.arn())
                    .assignPublicIp(true)
                    .taskDefinitionArgs(FargateServiceTaskDefinitionArgs.builder()
                            .container(TaskDefinitionContainerDefinitionArgs.builder()
                                    .name("app")
                                    .image(image.imageUri())
                                    .cpu(cpu)
                                    .memory(memory)
                                    .essential(true)
                                    .portMappings(List.of(
                                            TaskDefinitionPortMappingArgs.builder()
                                                    .containerPort(containerPort)
                                                    .targetGroup(loadbalancer.defaultTargetGroup())
                                                    .build()))
                                    .build())
                            .build())
                    .build());

            // Export the URL at which the container's HTTP endpoint will be available
            var dnsName = loadbalancer.loadBalancer().apply(lb -> lb.dnsName());
            ctx.export("url", dnsName.applyValue(d -> "http://" + d));
        });
    }
}
