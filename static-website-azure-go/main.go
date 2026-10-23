package main

import (
	"net/url"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v3"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v3"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		path := "www"
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

		// Create a resource group for the website.
		resourceGroup, err := resources.NewResourceGroup(ctx, "resource-group", nil)
		if err != nil {
			return err
		}

		// Create a blob storage account.
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

		// Configure the storage account as a website.
		website, err := storage.NewStorageAccountStaticWebsite(ctx, "website", &storage.StorageAccountStaticWebsiteArgs{
			ResourceGroupName: resourceGroup.Name,
			AccountName:       account.Name,
			IndexDocument:     pulumi.String(indexDocument),
			Error404Document:  pulumi.String(errorDocument),
		})
		if err != nil {
			return err
		}

		// Use a synced folder to manage the files of the website.
		_, err = synced.NewAzureBlobFolder(ctx, "synced-folder", &synced.AzureBlobFolderArgs{
			Path:               pulumi.String(path),
			ResourceGroupName:  resourceGroup.Name,
			StorageAccountName: account.Name,
			ContainerName:      website.ContainerName,
		})
		if err != nil {
			return err
		}

		// Pull the hostname out of the storage-account endpoint.
		originHostname := account.PrimaryEndpoints.ApplyT(func(endpoints storage.EndpointsResponse) (string, error) {
			parsed, err := url.Parse(endpoints.Web)
			if err != nil {
				return "", err
			}
			return parsed.Hostname(), nil
		}).(pulumi.StringOutput)

		// Export the URL and hostname of the storage account's website.
		ctx.Export("originURL", account.PrimaryEndpoints.ApplyT(func(endpoints storage.EndpointsResponse) (string, error) {
			return endpoints.Web, nil
		}).(pulumi.StringOutput))
		ctx.Export("originHostname", originHostname)

		return nil
	})
}
