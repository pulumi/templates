"""An Azure DevOps Python Pulumi program"""

import pulumi
import pulumi_random as random
import pulumi_azuredevops as ado

# Generate Azure DevOps project name
project_name = random.RandomPet("demo-project-name")

# Create Azure DevOps project
project = ado.Project("demo-project",
                           name=project_name)

# Export the name of created project
pulumi.export("project_name", project.name)
