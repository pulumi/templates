import pulumi
from pulumi_digitalocean import do

# Create a DigitalOcean resource (Domain)
domain = do.Domain('my-test-domain')

# Export the name of the domain
pulumi.export('domain_name',  domain.name)
