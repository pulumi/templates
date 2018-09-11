import pulumi
from pulumi_azure import core, storage

# Create an Azure Resource Group
resource_group = core.ResourceGroup("resource_group", 
    location='WestUS')

# Create an Azure resource (Storage Account)
account = storage.Account("storage", 
    resource_group_name=resource_group.name,
    location=resource_group.location,
    account_tier='Standard',
    account_replication_type='LRS')

# Export the connection string for the storage account
pulumi.output('connection_string', account.primary_connection_string)
