import pulumi
import pulumi_azure_native as azure
import pulumi_synced_folder as synced

# Import the program's configuration settings.
config = pulumi.Config()
www_path = config.get("sitePath", "./www")
app_path = config.get("appPath", "./app")
index_document = config.get("indexDocument", "index.html")
error_document = config.get("errorDocument", "error.html")

# Create a resource group for the website.
resource_group = azure.resources.ResourceGroup("resource-group")

# Create a blob storage account.
account = azure.storage.StorageAccount(
    "account",
    resource_group_name=resource_group.name,
    kind=azure.storage.Kind.STORAGE_V2,
    sku={
        "name": azure.storage.SkuName.STANDARD_LRS,
    },
)

# Create a storage container for the pages of the website.
website = azure.storage.StorageAccountStaticWebsite(
    "website",
    account_name=account.name,
    resource_group_name=resource_group.name,
    index_document=index_document,
    error404_document=error_document,
)

# Use a synced folder to manage the files of the website.
synced_folder = synced.AzureBlobFolder(
    "synced-folder",
    path=www_path,
    resource_group_name=resource_group.name,
    storage_account_name=account.name,
    container_name=website.container_name,
)

# Create a storage container for the serverless app.
app_container = azure.storage.BlobContainer(
    "app-container",
    account_name=account.name,
    resource_group_name=resource_group.name,
    public_access=azure.storage.PublicAccess.NONE,
)

# Upload the serverless app to the storage container.
app_blob = azure.storage.Blob(
    "app-blob",
    account_name=account.name,
    resource_group_name=resource_group.name,
    container_name=app_container.name,
    source=pulumi.FileArchive(app_path),
)

# Create a shared access signature to give the Function App access to the code.
signature = (
    pulumi.Output.all(resource_group.name, account.name, app_container.name)
    .apply(
        lambda args: azure.storage.list_storage_account_service_sas_output(
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
            canonicalized_resource=f"/blob/{args[1]}/{args[2]}",
        )
    )
    .apply(lambda result: result.service_sas_token)
)

# Create an App Service plan for the Function App.
plan = azure.web.AppServicePlan(
    "plan",
    resource_group_name=resource_group.name,
    kind="Linux",
    reserved=True,
    sku={
        "name": "Y1",
        "tier": "Dynamic",
    },
)

# Create the Function App.
app = azure.web.WebApp(
    "app",
    resource_group_name=resource_group.name,
    server_farm_id=plan.id,
    kind="FunctionApp",
    site_config={
        "app_settings": [
            {
                "name": "FUNCTIONS_WORKER_RUNTIME",
                "value": "python",
            },
            {
                "name": "FUNCTIONS_EXTENSION_VERSION",
                "value": "~3",
            },
            {
                "name": "WEBSITE_RUN_FROM_PACKAGE",
                "value": pulumi.Output.all(
                    account.name, app_container.name, app_blob.name, signature
                ).apply(
                    lambda args: f"https://{args[0]}.blob.core.windows.net/{args[1]}/{args[2]}?{args[3]}"
                ),
            },
        ],
        "cors": {
            "allowed_origins": ["*"],
        },
    },
)

# Create a JSON configuration file for the website.
site_config = azure.storage.Blob(
    "config.json",
    account_name=account.name,
    resource_group_name=resource_group.name,
    container_name=website.container_name,
    content_type="application/json",
    source=app.default_host_name.apply(
        lambda hostname: pulumi.StringAsset('{ "api": "https://' + hostname + '/api" }')
    ),
)

# Export the URLs of the website and serverless endpoint.
pulumi.export("siteURL", account.primary_endpoints.web)
pulumi.export(
    "apiURL",
    app.default_host_name.apply(
        lambda default_host_name: f"https://{default_host_name}/api"
    ),
)
