package main

import (
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an Azure Resource Group
		resourceGroup, err := core.NewResourceGroup(ctx, "resourceGroup", &core.ResourceGroupArgs{
			Location: pulumi.String("WestUS2"),
		})
		if err != nil {
			return err
		}

		// Create an Azure resource (Storage Account)
		account, err := storage.NewAccount(ctx, "storage", &storage.AccountArgs{
			ResourceGroupName:      resourceGroup.Name,
			AccountTier:            pulumi.String("Standard"),
			AccountReplicationType: pulumi.String("LRS"),
		})
		if err != nil {
			return err
		}

		// Export the connection string for the storage account
		ctx.Export("connectionString", account.PrimaryConnectionString)
		return nil
	})
}
