from urllib.parse import urlparse
import pulumi
import pulumi_azure_native as azure_native
import pulumi_synced_folder as synced_folder

# Import the program's configuration settings.
config = pulumi.Config()
path = config.get("path") or "www"
index_document = config.get("indexDocument") or "index.html"
error_document = config.get("errorDocument") or "error.html"

# Create a resource group for the website.
resource_group = azure_native.resources.ResourceGroup("resource-group")

# Create a blob storage account.
account = azure_native.storage.StorageAccount(
    "account",
    resource_group_name=resource_group.name,
    kind="StorageV2",
    sku={
        "name": "Standard_LRS",
    },
)

# Configure the storage account as a website.
website = azure_native.storage.StorageAccountStaticWebsite(
    "website",
    resource_group_name=resource_group.name,
    account_name=account.name,
    index_document=index_document,
    error404_document=error_document,
)

# Use a synced folder to manage the files of the website.
synced_folder = synced_folder.AzureBlobFolder(
    "synced-folder",
    path=path,
    resource_group_name=resource_group.name,
    storage_account_name=account.name,
    container_name=website.container_name,
)

# Pull the hostname out of the storage-account endpoint.
origin_hostname = account.primary_endpoints.web.apply(
    lambda endpoint: urlparse(endpoint).hostname
)

# Export the URL and hostname of the storage account's website.
pulumi.export("originURL", account.primary_endpoints.web)
pulumi.export("originHostname", origin_hostname)
