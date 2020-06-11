"""A Linode Python Pulumi program"""

import pulumi
import pulumi_linode

# Create a Linode resource (Linode Instance)
instance = pulumi_linode.Instance('my-instance', type='g6-nanode-1', region='us-east',
                                  image='linode/ubuntu18.04')

# Export the Instance label of the instance
pulumi.export('instance_label', instance.label)
