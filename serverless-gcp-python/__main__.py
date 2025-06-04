from cProfile import run
from pip import main
import pulumi
import pulumi_gcp as gcp
import pulumi_synced_folder as synced

# Import the program's configuration settings.
config = pulumi.Config()
site_path = config.get("sitePath", "./www")
app_path = config.get("appPath", "./app")
index_document = config.get("indexDocument", "index.html")
error_document = config.get("errorDocument", "error.html")

# Create a storage bucket and configure it as a website.
site_bucket = gcp.storage.Bucket(
    "site-bucket",
    location="US",
    website={
        "main_page_suffix": index_document,
        "not_found_page": error_document,
    },
)

# Create an IAM binding to allow public read access to the bucket.
site_bucket_iam_binding = gcp.storage.BucketIAMBinding(
    "site-bucket-iam-binding",
    bucket=site_bucket.name,
    role="roles/storage.objectViewer",
    members=["allUsers"],
)

# Use a synced folder to manage the files of the website.
synced_folder = synced.GoogleCloudFolder(
    "synced-folder",
    path=site_path,
    bucket_name=site_bucket.name,
)

# Create another storage bucket for the serverless app.
app_bucket = gcp.storage.Bucket(
    "app-bucket",
    location="US",
)

# Upload the serverless app to the storage bucket.
app_archive = gcp.storage.BucketObject(
    "app-archive",
    bucket=app_bucket.name,
    source=pulumi.asset.FileArchive(app_path),
)

# Create a Cloud Function that returns some data.
data_function = gcp.cloudfunctions.Function(
    "data-function",
    source_archive_bucket=app_bucket.name,
    source_archive_object=app_archive.name,
    runtime="python310",
    entry_point="data",
    trigger_http=True,
)

# Create an IAM member to invoke the function.
invoker = gcp.cloudfunctions.FunctionIamMember(
    "data-function-invoker",
    project=data_function.project,
    region=data_function.region,
    cloud_function=data_function.name,
    role="roles/cloudfunctions.invoker",
    member="allUsers",
)

# Create a JSON configuration file for the website.
site_config = gcp.storage.BucketObject(
    "site-config",
    name="config.json",
    bucket=site_bucket.name,
    content_type="application/json",
    source=data_function.https_trigger_url.apply(
        lambda url: pulumi.StringAsset('{ "api": "' + url + '" }')
    ),
)

# Export the URLs of the website and serverless endpoint.
pulumi.export(
    "siteURL",
    site_bucket.name.apply(
        lambda name: f"https://storage.googleapis.com/{name}/index.html"
    ),
)
pulumi.export("apiURL", data_function.https_trigger_url.apply(lambda url: url))
