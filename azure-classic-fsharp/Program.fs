module Program

open Pulumi.FSharp
open Pulumi.Azure.Core
open Pulumi.Azure.Storage

let infra () =
    // Create an Azure Resource Group
    let resourceGroup = ResourceGroup "resourceGroup"

    // Create an Azure Storage Account
    let storageAccount =
        Account("storage",
            AccountArgs
               (ResourceGroupName = io resourceGroup.Name,
                AccountReplicationType = input "LRS",
                AccountTier = input "Standard"))

    // Export the connection string for the storage account
    dict [("connectionString", storageAccount.PrimaryConnectionString :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
