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

    // Get the primary key
    let primaryKey =
        ListStorageAccountKeysInvokeArgs(ResourceGroupName = resourceGroup.Name, AccountName = storageAccount.Name)
        |> ListStorageAccountKeys.Invoke
        |> Outputs.bind (fun storageKeys -> Output.CreateSecret(storageKeys.Keys[0].Value))

    // Export the primary key for the storage account
    dict [("connectionString", primaryKey :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
