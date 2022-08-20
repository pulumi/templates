package main

import (
	resources "github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	storage "github.com/pulumi/pulumi-azure-native/sdk/go/azure/storage"
	"github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		path := "./site"
		if param := cfg.Get("path"); param != "" {
			path = param
		}
		indexDocument := "index.html"
		if param := cfg.Get("indexDocument"); param != "" {
			indexDocument = param
		}
		errorDocument := "error.html"
		if param := cfg.Get("errorDocument"); param != "" {
			errorDocument = param
		}
		resourceGroup, err := resources.NewResourceGroup(ctx, "resource-group", nil)
		if err != nil {
			return err
		}
		account, err := storage.NewStorageAccount(ctx, "account", &storage.StorageAccountArgs{
			ResourceGroupName: resourceGroup.Name,
			Kind:              pulumi.String("StorageV2"),
			Sku: &storage.SkuArgs{
				Name: pulumi.String("Standard_LRS"),
			},
		})
		if err != nil {
			return err
		}
		website, err := storage.NewStorageAccountStaticWebsite(ctx, "website", &storage.StorageAccountStaticWebsiteArgs{
			ResourceGroupName: resourceGroup.Name,
			AccountName:       account.Name,
			IndexDocument:     pulumi.String(indexDocument),
			Error404Document:  pulumi.String(errorDocument),
		})
		if err != nil {
			return err
		}
		_, err = synced - folder.NewAzureBlobFolder(ctx, "synced-folder", &synced-folder.AzureBlobFolderArgs{
			Path:               pulumi.String(path),
			ResourceGroupName:  resourceGroup.Name,
			StorageAccountName: account.Name,
			ContainerName:      website.ContainerName,
		})
		if err != nil {
			return err
		}
		ctx.Export("url", account.PrimaryEndpoints.ApplyT(func(primaryEndpoints storage.EndpointsResponse) (string, error) {
			return primaryEndpoints.Web, nil
		}).(pulumi.StringOutput))
		return nil
	})
}
