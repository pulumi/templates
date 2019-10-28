module Program

open Pulumi.FSharp

let infra () =
  //
  // Add you resources here
  //

  // Export outputs here
  dict []

[<EntryPoint>]
let main _ =
  Deployment.run infra
