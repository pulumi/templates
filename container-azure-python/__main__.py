import pulumi
import pulumi_docker as docker
import pulumi_random as random
from pulumi_azure_native import resources, containerregistry, containerinstance

# Import the program's configuration settings.
config = pulumi.Config()
app_path = config.get("appPath", "./app")
image_name = config.get("imageName", "my-app")
image_tag = config.get("imageTag", "latest")
container_port = config.get_int("containerPort", 80)
cpu = config.get_int("cpu", 1)
memory = config.get_int("memory", 2)

# Create a resource group for the container registry.
resource_group = resources.ResourceGroup("resource-group")

# Create a container registry.
registry = containerregistry.Registry(
    "registry",
    containerregistry.RegistryArgs(
        resource_group_name=resource_group.name,
        admin_user_enabled=True,
        sku=containerregistry.SkuArgs(
            name=containerregistry.SkuName.BASIC,
        ),
    ),
)

# Fetch login credentials for the registry.
credentials = containerregistry.list_registry_credentials_output(
    resource_group_name=resource_group.name,
    registry_name=registry.name,
)

registry_username = credentials.apply(lambda creds: creds.username)
registry_password = credentials.apply(lambda creds: creds.passwords[0].value)

# Create a container image for the service.
image = docker.Image(
    "image",
    image_name=pulumi.Output.concat(registry.login_server, f"/{image_name}:{image_tag}"),
    build=docker.DockerBuildArgs(
        context=app_path,
        platform="linux/amd64",
    ),
    registry=docker.RegistryArgs(
        server=registry.login_server,
        username=registry_username,
        password=registry_password,
    ),
)

# Use a random string to give the service a unique DNS name.
dns_name = random.RandomString(
    "dns-name",
    random.RandomStringArgs(
        length=8,
        special=False,
    ),
).result.apply(lambda result: f"{image_name}-{result.lower()}")

# Create a container group for the service that makes it publicly accessible.
container_group = containerinstance.ContainerGroup(
    "container-group",
    containerinstance.ContainerGroupArgs(
        resource_group_name=resource_group.name,
        os_type="linux",
        restart_policy="always",
        image_registry_credentials=[
            containerinstance.ImageRegistryCredentialArgs(
                server=registry.login_server,
                username=registry_username,
                password=registry_password,
            ),
        ],
        containers=[
            containerinstance.ContainerArgs(
                name=image_name,
                image=image.image_name,
                ports=[
                    containerinstance.ContainerPortArgs(
                        port=container_port,
                        protocol="tcp",
                    ),
                ],
                environment_variables=[
                    containerinstance.EnvironmentVariableArgs(
                        name="FLASK_RUN_PORT",
                        value=str(container_port),
                    ),
                    containerinstance.EnvironmentVariableArgs(
                        name="FLASK_RUN_HOST",
                        value="0.0.0.0",
                    ),
                ],
                resources=containerinstance.ResourceRequirementsArgs(
                    requests=containerinstance.ResourceRequestsArgs(
                        cpu=cpu,
                        memory_in_gb=memory,
                    ),
                ),
            ),
        ],
        ip_address=containerinstance.IpAddressArgs(
            type=containerinstance.ContainerGroupIpAddressType.PUBLIC,
            dns_name_label=dns_name,
            ports=[
                containerinstance.PortArgs(
                    port=container_port,
                    protocol="tcp",
                ),
            ],
        ),
    ),
)

# Export the service's IP address, hostname, and fully-qualified URL.
pulumi.export("hostname", container_group.ip_address.apply(lambda addr: addr.fqdn))
pulumi.export("ip", container_group.ip_address.apply(lambda addr: addr.ip))
pulumi.export(
    "url",
    container_group.ip_address.apply(
        lambda addr: f"http://{addr.fqdn}:{container_port}"
    ),
)
