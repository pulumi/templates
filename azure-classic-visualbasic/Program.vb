Imports Pulumi
Imports Pulumi.Azure.Core
Imports Pulumi.Azure.Storage

Module Program

    Sub Main()
        Deployment.RunAsync(AddressOf Infra).Wait()
    End Sub

    Private Function Infra() As IDictionary(Of String,Object)
        ' Create an Azure Resource Group
        Dim resourceGroup = New ResourceGroup("resourceGroup")
        
        Dim storageAccountArgs = New AccountArgs With {
            .ResourceGroupName = resourceGroup.Name,
            .AccountReplicationType = "LRS",
            .AccountTier = "Standard"
        }
        
        ' Create an Azure Storage Account
        Dim storageAccount = New Account("storage", storageAccountArgs)

        ' Export the connection string for the storage account
        Dim outputs = New Dictionary(Of String, Object)
        outputs.Add("connectionString", storageAccount.PrimaryConnectionString)
        Return outputs
    End Function
End Module
