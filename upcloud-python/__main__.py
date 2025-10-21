import pulumi
import pulumi_upcloud as upcloud

# Load user input from Pulumi config
config = pulumi.Config()
object_storage_name = config.require("object_storage_name")
region = config.require("region")
bucket_name = config.require("bucket_name")

# Create an UpCloud Managed Object Storage
object_storage = upcloud.ManagedObjectStorage(
    "objectStorage",
    name=object_storage_name,
    region=region,
    configured_status="started",
)

# Create a Bucket inside the Object Storage
bucket = upcloud.ManagedObjectStorageBucket(
    "storageBucket",
    service_uuid=object_storage.id,
    name=bucket_name,
)

# Export outputs
pulumi.export("object_storage_uuid", object_storage.id)
pulumi.export("bucket_name", bucket.name)
