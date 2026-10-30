module Program

open Pulumi.FSharp
open Pulumi.Gcp.Storage

let infra () =
  // Create a GCP resource (Storage Bucket)
  let bucket = Bucket("my-bucket", BucketArgs(Location = "US"))

  // Export the bucket's gs:// URL
  dict [("bucketName", bucket.Url :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
