module Program

open Pulumi.FSharp
open Pulumi.Aws.S3

let infra () =

  // Create an AWS resource (S3 Bucket)
  let bucket = Bucket "my-bucket"

  // Export the name of the bucket
  dict [("bucketName", bucket.Id :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
