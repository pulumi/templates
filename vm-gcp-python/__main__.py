import pulumi
import pulumi_gcp as gcp

# Import the program's configuration settings.
config = pulumi.Config()
machine_type = config.get("machineType", "f1-micro")
os_image = config.get("osImage", "debian-11")
instance_tag = config.get("instanceTag", "webserver")
service_port = config.get("servicePort", "80")

# Create a new network for the virtual machine.
network = gcp.compute.Network(
    "network",
    auto_create_subnetworks=False,
)

# Create a subnet on the network.
subnet = gcp.compute.Subnetwork(
    "subnet",
    ip_cidr_range="10.0.1.0/24",
    network=network.id,
)

# Create a firewall allowing inbound access over ports 80 (for HTTP) and 22 (for SSH).
firewall = gcp.compute.Firewall(
    "firewall",
    network=network.self_link,
    allows=[
        {
            "protocol": "tcp",
            "ports": [
                "22",
                service_port,
            ],
        },
    ],
    direction="INGRESS",
    source_ranges=[
        "0.0.0.0/0",
    ],
    target_tags=[
        instance_tag,
    ],
)

# Define a script to be run when the VM starts up.
metadata_startup_script = f"""#!/bin/bash
    echo '<!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="utf-8">
        <title>Hello, world!</title>
    </head>
    <body>
        <h1>Hello, world! 👋</h1>
        <p>Deployed with 💜 by <a href="https://pulumi.com/">Pulumi</a>.</p>
    </body>
    </html>' > index.html
    sudo python3 -m http.server {service_port} &
    """

# Create the virtual machine.
instance = gcp.compute.Instance(
    "instance",
    machine_type=machine_type,
    boot_disk={
        "initialize_params": {
            "image": os_image,
        },
    },
    network_interfaces=[
        {
            "network": network.id,
            "subnetwork": subnet.id,
            "access_configs": [
                {},
            ],
        },
    ],
    service_account={
        "scopes": [
            "https://www.googleapis.com/auth/cloud-platform",
        ],
    },
    allow_stopping_for_update=True,
    metadata_startup_script=metadata_startup_script,
    tags=[
        instance_tag,
    ],
    opts=pulumi.ResourceOptions(depends_on=firewall),
)

instance_ip = instance.network_interfaces.apply(
    lambda interfaces: interfaces[0].access_configs
    and interfaces[0].access_configs[0].nat_ip
)

# Export the instance's name, public IP address, and HTTP URL.
pulumi.export("name", instance.name)
pulumi.export("ip", instance_ip)
pulumi.export("url", instance_ip.apply(lambda ip: f"http://{ip}:{service_port}"))
