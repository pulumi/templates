import pulumi
import pulumi_aiven as aiven

config = pulumi.Config();
projectName = config.require("projectName");
cloudName = config.require("cloudName");
planName = config.require("planName");
serviceName = config.require("serviceName");

# Create a Kafka service.
kafka = aiven.Kafka("kafka",
    project=projectName,
    cloud_name=cloudName,
    plan=planName,
    service_name=serviceName,
)

# Export the service host and port.
pulumi.export("serviceHost", kafka.service_host)
pulumi.export("servicePort", kafka.service_port)
