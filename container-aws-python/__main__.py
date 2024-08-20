from pulumi import Config, Output, export
import pulumi_aws as aws
import pulumi_awsx as awsx

config = Config()
container_port = config.get_int("containerPort", 80)
cpu = config.get_int("cpu", 512)
memory = config.get_int("memory", 128)

# An ECS cluster to deploy into
cluster = aws.ecs.Cluster("cluster")

# An ALB to serve the container endpoint to the internet
loadbalancer = awsx.lb.ApplicationLoadBalancer("loadbalancer")

# An ECR repository to store our application's container image
repo = awsx.ecr.Repository("repo", force_delete=True)

# Build and publish our application's container image from ./app to the ECR repository
image = awsx.ecr.Image(
    "image", repository_url=repo.url, context="./app", platform="linux/amd64"
)

# Deploy an ECS Service on Fargate to host the application container
service = awsx.ecs.FargateService(
    "service",
    cluster=cluster.arn,
    assign_public_ip=True,
    task_definition_args={
        "container": {
            "name": "app",
            "image": image.image_uri,
            "cpu": cpu,
            "memory": memory,
            "essential": True,
            "port_mappings": [
                {
                    "container_port": container_port,
                    "target_group": loadbalancer.default_target_group,
                }
            ],
        },
    },
)

# The URL at which the container's HTTP endpoint will be available
export("url", Output.concat("http://", loadbalancer.load_balancer.dns_name))
