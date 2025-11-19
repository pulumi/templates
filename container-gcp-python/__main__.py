import pulumi
import pulumi_docker_build as docker_build
from pulumi_gcp import artifactregistry, cloudrun
import pulumi_random as random

# Import the program's configuration settings.
config = pulumi.Config()
app_path = config.get("appPath", "./app")
image_name = config.get("imageName", "my-app")
container_port = config.get_int("containerPort", 8080)
cpu = config.get_int("cpu", 1)
memory = config.get("memory", "1Gi")
concurrency = config.get_int("concurrency", 50)

# Import the provider's configuration settings.
gcp_config = pulumi.Config("gcp")
location = gcp_config.require("region")
project = gcp_config.require("project")

# Create a unique Artifact Registry repository ID
unique_string = random.RandomString(
    "unique-string",
    length=4,
    lower=True,
    upper=False,
    numeric=True,
    special=False,
)
repo_id = pulumi.Output.concat(
    "repo-",
    unique_string.result
)

# Create an Artifact Registry repository
repository = artifactregistry.Repository(
    "repository",
    description="Repository for container image",
    format="DOCKER",
    location=location,
    repository_id=repo_id,
)

# Form the repository URL
repo_url = pulumi.Output.concat(
    location,
    "-docker.pkg.dev/",
    project,
    "/",
    repository.repository_id
)

# Create a container image for the service.
# Before running `pulumi up`, configure Docker for Artifact Registry authentication
# as described here: https://cloud.google.com/artifact-registry/docs/docker/authentication
image = docker_build.Image(
    "image",
    tags=[pulumi.Output.concat(repo_url, "/", image_name)],
    context=docker_build.BuildContextArgs(
        location=app_path,
    ),
    # Cloud Run currently requires x86_64 images
    # https://cloud.google.com/run/docs/container-contract#languages
    platforms=[docker_build.Platform.LINUX_AMD64],
)

# Create a Cloud Run service definition.
service = cloudrun.Service(
    "service",
    location=location,
    template={
        "spec": {
            "containers": [
                {
                    "image": image.repo_digest,
                    "resources": {
                        "limits": {
                            "memory": memory,
                            "cpu": str(cpu),
                        },
                    },
                    "ports": [
                        {
                            "container_port": container_port,
                        },
                    ],
                    "envs": [
                        {
                            "name": "FLASK_RUN_PORT",
                            "value": str(container_port),
                        },
                    ],
                },
            ],
            "container_concurrency": concurrency,
        },
    },
)

# Create an IAM member to make the service publicly accessible.
invoker = cloudrun.IamMember(
    "invoker",
    location=location,
    service=service.name,
    role="roles/run.invoker",
    member="allUsers",
)

# Export the URL of the service.
pulumi.export("url", service.statuses.apply(lambda statuses: statuses[0].url))
