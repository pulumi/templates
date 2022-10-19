
import pulumi
import pulumi_docker as docker
from pulumi_gcp import cloudrun, config as gcp_config

config = pulumi.Config()
app_path = config.get("appPath", "./app")
image_name = config.get("imageName", "my-app")
container_port = config.get_int("containerPort", 8080)
cpu = config.get_int("cpu", 1)
memory = config.get("memory", "1Gi")
concurrency = config.get_float("concurrency", 50)

gcp_config = pulumi.Config("gcp")
location = gcp_config.require("region")
project = gcp_config.require("project")

# https://cloud.google.com/run/docs/about-concurrency
# gcloud services enable containerregistry.googleapis.com (may need this, I don't think I did)
# gcloud auth configure-docker

image = docker.Image(
    "image",
    image_name=f"gcr.io/{project}/{image_name}",
    build=docker.DockerBuild(
        context=app_path,
    ),
)

service = cloudrun.Service("service", cloudrun.ServiceArgs(
    location=location,
    template=cloudrun.ServiceTemplateArgs(
        spec=cloudrun.ServiceTemplateSpecArgs(
            containers=[
                cloudrun.ServiceTemplateSpecContainerArgs(
                    image=image.image_name,
                    resources=cloudrun.ServiceTemplateSpecContainerResourcesArgs(
                        limits=dict(
                            memory=memory,
                            cpu=cpu,
                        ),
                    ),
                    ports=[
                        cloudrun.ServiceTemplateSpecContainerPortArgs(
                            container_port=container_port,
                            protocol="TCP",
                        ),
                    ],
                    envs=[
                        cloudrun.ServiceTemplateSpecContainerEnvArgs(
                            name="FLASK_RUN_PORT",
                            value=container_port,
                        ),
                    ],
                ),
            ],
            container_concurrency=concurrency,
        ),
    ),
))

invoker = cloudrun.IamMember("invoker", cloudrun.IamMemberArgs(
    location=location,
    service=service.name,
    role="roles/run.invoker",
    member="allUsers",
))

pulumi.export("url", service.statuses.apply(lambda statuses: statuses[0].url))
