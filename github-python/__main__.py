"""A GitHub Python Pulumi program"""

import pulumi
import pulumi_github as github

# Create a GitHub repository
repository = github.Repository('demo-repo', description="Demo Repository for GitHub")

# Export the Name of the repository
pulumi.export('name', repository.name)
