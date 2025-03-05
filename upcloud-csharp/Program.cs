using Pulumi;
using UpCloud.Pulumi.UpCloud;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    var config = new Pulumi.Config();

    var objectStorageName = config.Require("object_storage_name");
    var region = config.Require("region");
    var bucketName = config.Require("bucket_name");

    var objectStorage = new ManagedObjectStorage("objectStorage", new()
    {
        Name = objectStorageName,
        Region = region,
        ConfiguredStatus = "started"
    });

    var bucket = new ManagedObjectStorageBucket("storageBucket", new()
    {
        ServiceUuid = objectStorage.Id,
        Name = bucketName
    });

    return new Dictionary<string, object?>
    {
        ["object_storage_uuid"] = objectStorage.Id,
        ["bucket_name"] = bucket.Name
    };
});