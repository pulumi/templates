"""An OpenStack Python Pulumi program"""

import pulumi
from pulumi_openstack import compute

# Create an OpenStack resource (Compute Instance)
instance = compute.Instance('test',
	flavor_name='s1-2',
	image_name='Ubuntu 22.04')

# Export the IP of the instance
pulumi.export('instance_ip', instance.access_ip_v4)
