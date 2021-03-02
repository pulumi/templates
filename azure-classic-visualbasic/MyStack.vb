Imports Pulumi
Imports Pulumi.Azure.Core
Imports Pulumi.Azure.Storage

Class MyStack
    Inherits Stack

    Public Sub New()
        ' Create an Azure Resource Group
        Dim resourceGroup = New ResourceGroup("resourceGroup")

        ' Create an Azure Storage Account
        Dim storageAccount = New Account("storage", New AccountArgs With {
            .ResourceGroupName = resourceGroup.Name,
            .AccountReplicationType = "LRS",
            .AccountTier = "Standard"
        })

        ' Export the connection string for the storage account
        Me.ConnectionString = storageAccount.PrimaryConnectionString
    End Sub

    <Output>
    Public Property ConnectionString As Output(Of String)
End Class
