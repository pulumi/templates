module Program

open Pulumi
open Pulumi.FSharp
open Pulumi.AzureNative.Resources
open Pulumi.AzureNative.Storage
open Pulumi.AzureNative.Storage.Inputs

let infra () =
    // Create an Azure Resource Group
    let resourceGroup = ResourceGroup("resourceGroup")

    // Create an Azure Storage Account
    let storageAccount =
        StorageAccount("sa",
            StorageAccountArgs
                (ResourceGroupName = resourceGroup.Name,
                 Sku = input (SkuArgs(Name = SkuName.Standard_LRS)),
                 Kind = Kind.StorageV2))

    // Export the storage account name
    dict [("storageAccountName", storageAccount.Name :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
