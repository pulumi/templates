package main

import (
	"fmt"

	resources "github.com/pulumi/pulumi-azure-native/sdk/go/azure/resources"
	storage "github.com/pulumi/pulumi-azure-native/sdk/go/azure/storage"
	web "github.com/pulumi/pulumi-azure-native/sdk/go/azure/web"
	synced "github.com/pulumi/pulumi-synced-folder/sdk/go/synced-folder"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		sitePath := "./www"
		if param := cfg.Get("sitePath"); param != "" {
			sitePath = param
		}
		apiPath := "./api"
		if param := cfg.Get("apiPath"); param != "" {
			apiPath = param
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

		// Create a storage container for the pages of the website.
		website, err := storage.NewStorageAccountStaticWebsite(ctx, "website", &storage.StorageAccountStaticWebsiteArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
			IndexDocument:     pulumi.String(indexDocument),
			Error404Document:  pulumi.String(errorDocument),
		})
		if err != nil {
			return err
		}

		// Use a synced folder to manage the files of the website.
		_, err = synced.NewAzureBlobFolder(ctx, "synced-folder", &synced.AzureBlobFolderArgs{
			Path:               pulumi.String(sitePath),
			ResourceGroupName:  resourceGroup.Name,
			StorageAccountName: account.Name,
			ContainerName:      website.ContainerName,
		})
		if err != nil {
			return err
		}

		// Create a storage container for serverless functions.
		container, err := storage.NewBlobContainer(ctx, "container", &storage.BlobContainerArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
			PublicAccess:      storage.PublicAccessNone,
		})
		if err != nil {
			return err
		}

		// Upload the functions to the container.
		blob, err := storage.NewBlob(ctx, "blob", &storage.BlobArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
			ContainerName:     container.Name,
			Source:            pulumi.NewFileArchive(apiPath),
		})
		if err != nil {
			return err
		}

		// Create an App Service plan for the Function App.
		plan, err := web.NewAppServicePlan(ctx, "plan", &web.AppServicePlanArgs{
			ResourceGroupName: resourceGroup.Name,
			Sku: &web.SkuDescriptionArgs{
				Name: pulumi.String("Y1"),
				Tier: pulumi.String("Dynamic"),
			},
		})
		if err != nil {
			return err
		}

		// Create the Function App.
		app, err := web.NewWebApp(ctx, "app", &web.WebAppArgs{
			ResourceGroupName: resourceGroup.Name,
			ServerFarmId:      plan.ID(),
			Kind:              pulumi.String("FunctionApp"),
			SiteConfig: &web.SiteConfigArgs{
				AppSettings: web.NameValuePairArray{
					&web.NameValuePairArgs{
						Name:  pulumi.String("FUNCTIONS_WORKER_RUNTIME"),
						Value: pulumi.String("node"),
					},
					&web.NameValuePairArgs{
						Name:  pulumi.String("WEBSITE_NODE_DEFAULT_VERSION"),
						Value: pulumi.String("~14"),
					},
					&web.NameValuePairArgs{
						Name:  pulumi.String("FUNCTIONS_EXTENSION_VERSION"),
						Value: pulumi.String("~3"),
					},
					&web.NameValuePairArgs{
						Name: pulumi.String("WEBSITE_RUN_FROM_PACKAGE"),
						Value: pulumi.All(resourceGroup.Name, account.Name, container.Name, blob.Name).ApplyT(
							func(args []interface{}) string {
								resourceGroupName := args[0].(string)
								accountName := args[1].(string)
								containerName := args[2].(string)
								blobName := args[3].(string)
								protocol := storage.HttpProtocolHttps

								// Create a shared access signature allowing access to function storage.
								result, err := storage.ListStorageAccountServiceSAS(ctx, &storage.ListStorageAccountServiceSASArgs{
									ResourceGroupName:      resourceGroupName,
									AccountName:            accountName,
									Protocols:              &protocol,
									SharedAccessStartTime:  pulumi.StringRef("2022-01-01"),
									SharedAccessExpiryTime: pulumi.StringRef("2030-01-01"),
									Resource:               pulumi.StringRef("c"),
									Permissions:            pulumi.StringRef("r"),
									ContentType:            pulumi.StringRef("application/json"),
									CacheControl:           pulumi.StringRef("max-age=5"),
									ContentDisposition:     pulumi.StringRef("inline"),
									ContentEncoding:        pulumi.StringRef("deflate"),
									CanonicalizedResource:  fmt.Sprintf("/blob/%s/%s", accountName, containerName),
								})
								if err != nil {
									return ""
								}

								return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", accountName, containerName, blobName, result.ServiceSasToken)
							}).(pulumi.StringPtrInput),
					},
				},
				Cors: &web.CorsSettingsArgs{
					AllowedOrigins: pulumi.StringArray{
						pulumi.String("*"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// Create a JSON configuration file for the website.
		_, err = storage.NewBlob(ctx, "config.json", &storage.BlobArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
			ContainerName:     website.ContainerName,
			ContentType:       pulumi.StringPtr("application/json"),

			Source: app.DefaultHostName.ApplyT(func(hostname string) pulumi.AssetOrArchiveOutput {
				config := fmt.Sprintf(`{ "api": "https://%s/api" }`, hostname)
				return pulumi.NewStringAsset(config).ToAssetOrArchiveOutput()
			}).(pulumi.AssetOrArchiveOutput),
		})
		if err != nil {
			return err
		}

		// Export the URLs of the website and serverless endpoint.
		ctx.Export("originURL", account.PrimaryEndpoints.Web())
		ctx.Export("apiURL", pulumi.Sprintf("https://%v/api", app.DefaultHostName))
		return nil
	})
}
