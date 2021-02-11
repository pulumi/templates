"""An Equinix Metal Python Pulumi program"""

import pulumi
from pulumi_equinix_metal import metal

# Create an Equinix Metal resource (project)
project = metal.Project('my-test-project',
    name='my-test-project')

# Export the name of the project
pulumi.export('project_name', project.name)
