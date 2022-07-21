using Pulumi;
using Pulumi.Azure.Core;
using Pulumi.Azure.Storage;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Create an Azure Resource Group
    var resourceGroup = new ResourceGroup("resourceGroup");

    // Create an Azure Storage Account
    var storageAccount = new Account("storage", new AccountArgs
    {
        ResourceGroupName = resourceGroup.Name,
        AccountReplicationType = "LRS",
        AccountTier = "Standard"
    });

    // Export the connection string for the storage account
    return new Dictionary<string, object?>
    {
        ["connectionString"] = storageAccount.PrimaryConnectionString
    };
});
