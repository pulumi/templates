import pulumi
import pulumi_gcp as gcp
import pulumi_synced_folder as synced_folder

config = pulumi.Config()
gcp_project = config.get("gcpProject")
if gcp_project is None:
    gcp_project = "pulumi-development"
path = config.get("path")
if path is None:
    path = "./site"
index_document = config.get("indexDocument")
if index_document is None:
    index_document = "index.html"
error_document = config.get("errorDocument")
if error_document is None:
    error_document = "error.html"
bucket = gcp.storage.Bucket("bucket",
    location="US",
    website=gcp.storage.BucketWebsiteArgs(
        main_page_suffix=index_document,
        not_found_page=error_document,
    ))
bucket_iam_binding = gcp.storage.BucketIAMBinding("bucket-iam-binding",
    bucket=bucket.name,
    role="roles/storage.objectViewer",
    members=["allUsers"])
synced_folder = synced_folder.GoogleCloudFolder("synced-folder",
    path=path,
    bucket_name=bucket.name)
pulumi.export("url", bucket.name.apply(lambda name: f"https://storage.googleapis.com/{name}/index.html"))
