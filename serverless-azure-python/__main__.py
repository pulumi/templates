import pulumi
import pulumi_azure_native as azure
import pulumi_synced_folder as synced

# Import the program's configuration settings.
config = pulumi.Config()
www_path = config.get("sitePath") or "./www"
api_path = config.get("apiPath") or "./api"
index_document = config.get("indexDocument") or "index.html"
error_document = config.get("errorDocument") or "error.html"

# Create a resource group for the website.
resource_group = azure.resources.ResourceGroup("resource-group")

# Create a blob storage account.
account = azure.storage.StorageAccount("account",
    resource_group_name=resource_group.name,
    kind=azure.storage.Kind.STORAGE_V2,
    sku=azure.storage.SkuArgs(
        name=azure.storage.SkuName.STANDARD_LRS,
    ))

# Create a storage container for the pages of the website.
website = azure.storage.StorageAccountStaticWebsite("website",
    account_name=account.name,
    resource_group_name=resource_group.name,
    index_document=index_document,
    error404_document=error_document)

# Use a synced folder to manage the files of the website.
synced_folder = synced.AzureBlobFolder("synced-folder",
    path=www_path,
    resource_group_name=resource_group.name,
    storage_account_name=account.name,
    container_name=website.container_name)

# Create a storage container for serverless functions.
container = azure.storage.BlobContainer("container",
    account_name=account.name,
    resource_group_name=resource_group.name,
    public_access=azure.storage.PublicAccess.NONE)

# Upload the functions to the container.
blob = azure.storage.Blob("blob",
    account_name=account.name,
    resource_group_name=resource_group.name,
    container_name=container.name,
    source=pulumi.FileArchive(api_path))

# Create a shared access signature allowing access to function storage.
blob_sas = pulumi.Output.all(resource_group.name, account.name, container.name).apply(lambda args: azure.storage.list_storage_account_service_sas_output(
    resource_group_name=args[0],
    account_name=args[1],
    protocols=azure.storage.HttpProtocol.HTTPS,
    shared_access_start_time="2022-01-01",
    shared_access_expiry_time="2030-01-01",
    resource=azure.storage.SignedResource.C,
    permissions=azure.storage.Permissions.R,
    content_type="application/json",
    cache_control="max-age=5",
    content_disposition="inline",
    content_encoding="deflate",
    canonicalized_resource=f"/blob/{args[1]}/{args[2]}"))

# Create an App Service plan for the Function App.
plan = azure.web.AppServicePlan("plan",
    resource_group_name=resource_group.name,
    sku=azure.web.SkuDescriptionArgs(
        name="Y1",
        tier="Dynamic",
    ))

# Create the Function App.
app = azure.web.WebApp("app",
    resource_group_name=resource_group.name,
    server_farm_id=plan.id,
    kind="FunctionApp",
    site_config=azure.web.SiteConfigArgs(
        app_settings=[
            azure.web.NameValuePairArgs(
                name="FUNCTIONS_WORKER_RUNTIME",
                value="node",
            ),
            azure.web.NameValuePairArgs(
                name="WEBSITE_NODE_DEFAULT_VERSION",
                value="~14",
            ),
            azure.web.NameValuePairArgs(
                name="FUNCTIONS_EXTENSION_VERSION",
                value="~3",
            ),
            azure.web.NameValuePairArgs(
                name="WEBSITE_RUN_FROM_PACKAGE",
                value=pulumi.Output.all(account.name, container.name, blob.name, blob_sas).apply(
                    lambda args: f"https://{args[0]}.blob.core.windows.net/{args[1]}/{args[2]}?{args[3].service_sas_token}"),
            ),
        ],
        cors=azure.web.CorsSettingsArgs(
            allowed_origins=["*"],
        ),
    ))

# Create a JSON configuration file for the website.
site_config = azure.storage.Blob("config.json",
    account_name=account.name,
    resource_group_name=resource_group.name,
    container_name=website.container_name,
    content_type="application/json",
    source=app.default_host_name.apply(lambda hostname: pulumi.StringAsset("{ \"api\": \"https://" + hostname  + "/api\" }")))

# Export the URLs of the website and serverless endpoint.
pulumi.export("originURL", account.primary_endpoints.web)
pulumi.export("apiURL", app.default_host_name.apply(lambda default_host_name: f"https://{default_host_name}/api"))
