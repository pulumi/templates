module Program

open Pulumi.FSharp
open Pulumi.AliCloud.Oss

let infra () =

  // Create an AliCloud resource (OSS Bucket)
  let bucket = Bucket "my-bucket"

  // Export the name of the bucket
  dict [("bucketName", bucket.Id :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
