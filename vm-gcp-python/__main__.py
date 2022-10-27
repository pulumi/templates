import pulumi
import pulumi_gcp as gcp

config = pulumi.Config();
machine_type = config.get("machineType", "f1-micro");
os_image = config.get("osImage", "debian-11");
instance_tag = config.get("instanceTag", "webserver");
service_port = config.get_int("servicePort", 80);

network = gcp.compute.Network("network", gcp.compute.NetworkArgs(
    auto_create_subnetworks=False,
))

subnet = gcp.compute.Subnetwork("subnet", gcp.compute.SubnetworkArgs(
    ip_cidr_range="10.0.1.0/24",
    network=network.id,
))

firewall = gcp.compute.Firewall("firewall", gcp.compute.FirewallArgs(
    network=network.self_link,
    allows=[
        gcp.compute.FirewallAllowArgs(
            protocol="tcp",
            ports=[
                "22",
                str(service_port),
            ],
        ),
    ],
    direction="INGRESS",
    source_ranges=[
        "0.0.0.0/0",
    ],
    target_tags=[
        instance_tag,
    ],
))

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

instance = gcp.compute.Instance("instance", gcp.compute.InstanceArgs(
    machine_type=machine_type,
    boot_disk=gcp.compute.InstanceBootDiskArgs(
        initialize_params=gcp.compute.InstanceBootDiskInitializeParamsArgs(
            image=os_image,
        ),
    ),
    network_interfaces=[
        gcp.compute.InstanceNetworkInterfaceArgs(
            network=network.id,
            subnetwork=subnet.id,
            access_configs=[
                {},
            ],
        ),
    ],
    service_account=gcp.compute.InstanceServiceAccountArgs(
        scopes=[
            "https://www.googleapis.com/auth/cloud-platform",
        ],
    ),
    allow_stopping_for_update=True,
    metadata_startup_script=metadata_startup_script,
    tags=[
        instance_tag,
    ],
), pulumi.ResourceOptions(depends_on=firewall))

instance_ip = instance.network_interfaces.apply(lambda interfaces: interfaces[0].access_configs[0].nat_ip)

pulumi.export("name", instance.name)
pulumi.export("ip", instance_ip)
pulumi.export("url", instance_ip.apply(lambda ip: f"http://{ip}"))