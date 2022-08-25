package main

import (
	"fmt"
	"net/url"

	cdn "github.com/pulumi/pulumi-azure-native/sdk/go/azure/cdn"
	resources "github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	storage "github.com/pulumi/pulumi-azure-native/sdk/go/azure/storage"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
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
		_, err = synced.NewAzureBlobFolder(ctx, "synced-folder", &synced.AzureBlobFolderArgs{
			Path:               pulumi.String(path),
			ResourceGroupName:  resourceGroup.Name,
			StorageAccountName: account.Name,
			ContainerName:      website.ContainerName,
		})
		if err != nil {
			return err
		}
		profile, err := cdn.NewProfile(ctx, "profile", &cdn.ProfileArgs{
			ResourceGroupName: resourceGroup.Name,
			Sku: &cdn.SkuArgs{
				Name: pulumi.String("Standard_Microsoft"),
			},
		})
		if err != nil {
			return err
		}

		originHostname := account.PrimaryEndpoints.ApplyT(func(endpoints storage.EndpointsResponse) (string, error) {
			parsed, err := url.Parse(endpoints.Web)
			if err != nil {
				return "", err
			}
			return parsed.Hostname(), nil
		}).(pulumi.StringOutput)

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
					Name: account.Name,
					HostName: originHostname,
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("originURL", account.PrimaryEndpoints.ApplyT(func(primaryEndpoints storage.EndpointsResponse) (string, error) {
			return primaryEndpoints.Web, nil
		}).(pulumi.StringOutput))
		ctx.Export("originHostname", originHostname)
		ctx.Export("cdnURL", endpoint.HostName.ApplyT(func(hostName string) (string, error) {
			return fmt.Sprintf("https://%v", hostName), nil
		}).(pulumi.StringOutput))
		ctx.Export("cdnHostname", endpoint.HostName)
		return nil
	})
}
