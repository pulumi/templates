import pulumi
import pulumi_azure_native as azure_native
import pulumi_synced_folder as synced_folder

config = pulumi.Config()
path = config.get("path")
if path is None:
    path = "./site"
index_document = config.get("indexDocument")
if index_document is None:
    index_document = "index.html"
error_document = config.get("errorDocument")
if error_document is None:
    error_document = "error.html"
resource_group = azure_native.resources.ResourceGroup("resource-group")
account = azure_native.storage.StorageAccount("account",
    resource_group_name=resource_group.name,
    kind="StorageV2",
    sku=azure_native.storage.SkuArgs(
        name="Standard_LRS",
    ))
website = azure_native.storage.StorageAccountStaticWebsite("website",
    resource_group_name=resource_group.name,
    account_name=account.name,
    index_document=index_document,
    error404_document=error_document)
synced_folder = synced_folder.AzureBlobFolder("synced-folder",
    path=path,
    resource_group_name=resource_group.name,
    storage_account_name=account.name,
    container_name=website.container_name)
pulumi.export("url", account.primary_endpoints.web)
