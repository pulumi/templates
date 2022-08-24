import pulumi
import pulumi_azure_native as azure_native
import pulumi_synced_folder as synced_folder
import urllib

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
profile = azure_native.cdn.Profile("profile",
    resource_group_name=resource_group.name,
    sku=azure_native.cdn.SkuArgs(
        name="Standard_Microsoft",
    ))
origin_hostname = account.primary_endpoints.web.apply(lambda endpoint: urllib.parse.urlparse(endpoint).hostname)
endpoint = azure_native.cdn.Endpoint("endpoint",
    resource_group_name=resource_group.name,
    profile_name=profile.name,
    is_http_allowed=False,
    is_https_allowed=True,
    is_compression_enabled=True,
    content_types_to_compress=[
        "text/html",
        "text/css",
        "application/javascript",
        "application/json",
        "image/svg+xml",
        "font/woff",
        "font/woff2",
    ],
    origin_host_header=origin_hostname,
    origins=[azure_native.cdn.DeepCreatedOriginArgs(
        name=account.name,
        host_name=origin_hostname,
    )])
pulumi.export("originURL", account.primary_endpoints.web)
pulumi.export("cdnURL", endpoint.host_name.apply(lambda host_name: f"https://{host_name}"))
