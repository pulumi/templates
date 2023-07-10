package main

import (
	"net/url"

	cdn "github.com/pulumi/pulumi-azure-native-sdk/cdn/v2"
	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v2"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		path := "./www"
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

		// Create a CDN profile.
		profile, err := cdn.NewProfile(ctx, "profile", &cdn.ProfileArgs{
			ResourceGroupName: resourceGroup.Name,
			Sku: &cdn.SkuArgs{
				Name: pulumi.String("Standard_Microsoft"),
			},
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

		// Create a CDN endpoint to distribute and cache the website.
		endpoint, err := cdn.NewEndpoint(ctx, "endpoint", &cdn.EndpointArgs{
			ResourceGroupName:    resourceGroup.Name,
			ProfileName:          profile.Name,
			IsHttpAllowed:        pulumi.Bool(false),
			IsHttpsAllowed:       pulumi.Bool(true),
			IsCompressionEnabled: pulumi.Bool(true),
			ContentTypesToCompress: pulumi.StringArray{
				pulumi.String("text/html"),
				pulumi.String("text/css"),
				pulumi.String("application/javascript"),
				pulumi.String("application/json"),
				pulumi.String("image/svg+xml"),
				pulumi.String("font/woff"),
				pulumi.String("font/woff2"),
			},
			OriginHostHeader: originHostname,
			Origins: cdn.DeepCreatedOriginArray{
				&cdn.DeepCreatedOriginArgs{
					Name:     account.Name,
					HostName: originHostname,
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the URLs and hostnames of the storage account and distribution.
		ctx.Export("originURL", account.PrimaryEndpoints.ApplyT(func(endpoints storage.EndpointsResponse) (string, error) {
			return endpoints.Web, nil
		}).(pulumi.StringOutput))
		ctx.Export("originHostname", originHostname)
		ctx.Export("cdnURL", pulumi.Sprintf("http://%s", endpoint.HostName))
		ctx.Export("cdnHostname", endpoint.HostName)

		return nil
	})
}
