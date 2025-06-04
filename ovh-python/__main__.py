import pulumi
import pulumi_ovh as ovh

# Get some configuration values
config = pulumi.Config();

ovh_service_name = config.require('ovhServiceName')
ovh_region = config.get("ovhRegion", "GRA")
plan_name = config.get("planName", "SMALL")

registry_name = config.get("registryName", "my-registry")
registry_user_name = config.get("registryUserName", "user")
registry_user_email = config.get("registryUserEmail", "myuser@ovh.com")
registry_user_login = config.get("registryUserLogin", "myuser")
	
# Initiate the configuration of the registry
regcap = ovh.cloudproject.get_capabilities_container_filter(service_name=ovh_service_name,
    plan_name=plan_name,
    region=ovh_region)

# Deploy a new Managed private registry
my_registry = ovh.cloudproject.ContainerRegistry(registry_name,
    service_name=regcap.service_name,
    plan_id=regcap.id,
    region=regcap.region)

# Create a Private Registry User
user = ovh.cloudproject.ContainerRegistryUser(registry_user_name,
    service_name=ovh_service_name,
    registry_id=my_registry.id,
    email=registry_user_email,
    login=registry_user_login)

# Add as an output registry information
pulumi.export("registryURL", my_registry.url)
pulumi.export("registryUser", user.user)
pulumi.export("registryPassword", user.password)
