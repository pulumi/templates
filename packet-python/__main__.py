import pulumi
from pulumi_packet import packet

# Create a Packet resource (project)
project = packet.Project("my-test-project",
    name='my-test-project')

# Export the name of the project
pulumi.export('project_name',  project.name)
