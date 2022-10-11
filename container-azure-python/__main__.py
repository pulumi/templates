import pulumi
import pulumi_docker as docker
import pulumi_random as random
from pulumi_azure_native import resources, containerregistry, containerinstance

config = pulumi.Config()
image_name = config.get("imageName", "my-app")
app_path = config.get("appPath", "./app")
container_port = config.get_number("containerPort", 80)

resource_group = resources.ResourceGroup('resource_group')

registry = containerregistry.Registry("registry", containerregistry.RegistryArgs(
    resource_group_name=resource_group.name,
    admin_user_enabled=True,
    sku=containerregistry.SkuArgs(
        name=containerregistry.SkuName.BASIC,
    ),
))

credentials = containerregistry.list_registry_credentials_output(
    resource_group_name=resource_group.name,
    registry_name=registry.name,
)

registry_username = credentials.apply(lambda creds: creds.username)
registry_password = credentials.apply(lambda creds: creds.passwords[0].value)

image = docker.Image("image",
    image_name=pulumi.Output.concat(registry.login_server, "/", image_name),
    build=docker.DockerBuild(
        context=app_path,
    ),
    registry=docker.ImageRegistry(
        server=registry.login_server,
        username=registry_username,
        password=registry_password,
    ),
)

hostname = random.RandomPet("hostname", random.RandomPetArgs(
    length=2,
))

group = containerinstance.ContainerGroup("group", containerinstance.ContainerGroupArgs(
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
                    cpu=1.0,
                    memory_in_gb=1.5,
                ),
            ),
        ),
    ],
    ip_address=containerinstance.IpAddressArgs(
        type=containerinstance.ContainerGroupIpAddressType.PUBLIC,
        dns_name_label=hostname,
        ports=[
            containerinstance.PortArgs(
                port=container_port,
                protocol="tcp",
            ),
        ],
    ),
))

pulumi.export("ipAddress", group.ip_address.apply(lambda addr: addr.ip))
pulumi.export("hostname", group.ip_address.apply(lambda addr: addr.fqdn))
pulumi.export("url", group.ip_address.apply(lambda addr: f"http://{addr.fqdn}:{container_port}"))
