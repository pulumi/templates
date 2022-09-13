module Program

open Pulumi.FSharp
open Pulumi.AwsNative.S3

let infra () =

  // Create an AWS resource (S3 Bucket)
  let bucket = Bucket "my-aws-native-fsharp-bucket"

  // Export the name of the bucket
  dict [("bucketName", bucket.Id :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
