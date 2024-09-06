import pulumi
import pulumi_gcp as gcp
import pulumi_synced_folder as synced_folder

# Import the program's configuration settings.
config = pulumi.Config()
path = config.get("path") or "./www"
index_document = config.get("indexDocument") or "index.html"
error_document = config.get("errorDocument") or "error.html"

# Create a storage bucket and configure it as a website.
bucket = gcp.storage.Bucket(
    "bucket",
    location="US",
    website={
        "main_page_suffix": index_document,
        "not_found_page": error_document,
    },
)

# Create an IAM binding to allow public read access to the bucket.
bucket_iam_binding = gcp.storage.BucketIAMBinding(
    "bucket-iam-binding",
    bucket=bucket.name,
    role="roles/storage.objectViewer",
    members=["allUsers"],
)

# Use a synced folder to manage the files of the website.
synced_folder = synced_folder.GoogleCloudFolder(
    "synced-folder", path=path, bucket_name=bucket.name
)

# Enable the storage bucket as a CDN.
backend_bucket = gcp.compute.BackendBucket(
    "backend-bucket", bucket_name=bucket.name, enable_cdn=True
)

# Provision a global IP address for the CDN.
ip = gcp.compute.GlobalAddress("ip")

# Create a URLMap to route requests to the storage bucket.
url_map = gcp.compute.URLMap("url-map", default_service=backend_bucket.self_link)

# Create an HTTP proxy to route requests to the URLMap.
http_proxy = gcp.compute.TargetHttpProxy("http-proxy", url_map=url_map.self_link)

# Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
http_forwarding_rule = gcp.compute.GlobalForwardingRule(
    "http-forwarding-rule",
    ip_address=ip.address,
    ip_protocol="TCP",
    port_range="80",
    target=http_proxy.self_link,
)

# Export the URLs and hostnames of the bucket and CDN.
pulumi.export(
    "originURL",
    bucket.name.apply(lambda name: f"https://storage.googleapis.com/{name}/index.html"),
)
pulumi.export(
    "originHostname", bucket.name.apply(lambda name: f"storage.googleapis.com/{name}")
)
pulumi.export("cdnURL", ip.address.apply(lambda address: f"http://{address}"))
pulumi.export("cdnHostname", ip.address)
