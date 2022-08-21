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
backend_bucket = gcp.compute.BackendBucket("backend-bucket",
    bucket_name=bucket.name,
    enable_cdn=True)
ip = gcp.compute.GlobalAddress("ip")
url_map = gcp.compute.URLMap("url-map", default_service=backend_bucket.self_link)
http_proxy = gcp.compute.TargetHttpProxy("http-proxy", url_map=url_map.self_link)
http_forwarding_rule = gcp.compute.GlobalForwardingRule("http-forwarding-rule",
    ip_address="ip.address",
    ip_protocol="TCP",
    port_range="80",
    target=http_proxy.self_link)
pulumi.export("originURL", bucket.name.apply(lambda name: f"https://storage.googleapis.com/{name}/index.html"))
pulumi.export("cdnURL", ip.address.apply(lambda address: f"http://{address}"))
