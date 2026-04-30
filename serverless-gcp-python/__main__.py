import pulumi
import pulumi_gcp as gcp
import pulumi_synced_folder as synced

# Import the program's configuration settings.
config = pulumi.Config()
site_path = config.get("sitePath", "./www")
app_path = config.get("appPath", "./app")
index_document = config.get("indexDocument", "index.html")
error_document = config.get("errorDocument", "error.html")
region = gcp.config.region or "us-central1"

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

# Create a Cloud Function (Gen 2) that returns some data.
data_function = gcp.cloudfunctionsv2.Function(
    "data-function",
    location=region,
    build_config={
        "runtime": "python312",
        "entry_point": "data",
        "source": {
            "storage_source": {
                "bucket": app_bucket.name,
                "object": app_archive.name,
            },
        },
    },
    service_config={
        "available_memory": "256M",
        "timeout_seconds": 60,
    },
)

# Allow public, unauthenticated invocations of the underlying Cloud Run service.
invoker = gcp.cloudrun.IamMember(
    "data-function-invoker",
    location=data_function.location,
    service=data_function.name,
    role="roles/run.invoker",
    member="allUsers",
)

# Create a JSON configuration file for the website.
site_config = gcp.storage.BucketObject(
    "site-config",
    name="config.json",
    bucket=site_bucket.name,
    content_type="application/json",
    source=data_function.url.apply(
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
pulumi.export("apiURL", data_function.url)
