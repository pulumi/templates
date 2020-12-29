module Program

open Pulumi.FSharp
open Pulumi.AzureNextGen.Resources.Latest
open Pulumi.AzureNextGen.Storage.Latest
open Pulumi.AzureNextGen.Storage.Latest.Inputs

// Helper function to retrieve the primary key of a storage account
let getStorageAccountPrimaryKey(resourceGroupName: string, accountName: string): Async<string> = async {
    let! accountKeys =
        ListStorageAccountKeysArgs(ResourceGroupName = resourceGroupName, AccountName = accountName) |>
        ListStorageAccountKeys.InvokeAsync
        |> Async.AwaitTask
    return accountKeys.Keys.[0].Value
}

let infra () =
    // Create an Azure Resource Group
    let resourceGroup =
        ResourceGroup("resourceGroup", 
            ResourceGroupArgs
                (ResourceGroupName = input "my-rg",
                 Location = input "WestUS"))

    // Create an Azure Storage Account
    let storageAccount =
        StorageAccount("sa", 
            StorageAccountArgs
                (ResourceGroupName = io resourceGroup.Name,
                 AccountName = input "mystorageaccount", // <-- change to a unique name
                 Location = io resourceGroup.Location,
                 Sku = input (SkuArgs(Name = inputUnion2Of2 SkuName.Standard_LRS)),
                 Kind = inputUnion2Of2 Kind.StorageV2))
        
    // Get the primary key
    let primaryKey =
        Outputs.pair resourceGroup.Name storageAccount.Name
        |> Outputs.applyAsync getStorageAccountPrimaryKey
    
    // Export the primary key for the storage account
    dict [("connectionString", primaryKey :> obj)]

[<EntryPoint>]
let main _ =
  Deployment.run infra
