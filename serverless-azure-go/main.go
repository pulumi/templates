package main

import (
	"fmt"

	resources "github.com/pulumi/pulumi-azure-native-sdk/resources/v2"
	storage "github.com/pulumi/pulumi-azure-native-sdk/storage/v2"
	web "github.com/pulumi/pulumi-azure-native-sdk/web/v2"
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
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
		appPath := "./app"
		if param := cfg.Get("appPath"); param != "" {
			appPath = param
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

		// Create a storage container for the serverless app.
		appContainer, err := storage.NewBlobContainer(ctx, "app-container", &storage.BlobContainerArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
			PublicAccess:      storage.PublicAccessNone,
		})
		if err != nil {
			return err
		}

		// Compile the the app for the Azure Linux environment.
		buildResult, err := local.Run(ctx, &local.RunArgs{
			Command:      fmt.Sprintf("GOOS=linux GOARCH=amd64 go build cmd/app.go"),
			Dir:          pulumi.StringRef(appPath),
			ArchivePaths: []string{"**"},
		})
		if err != nil {
			return err
		}

		// Upload the serverless app to the storage container.
		appBlob, err := storage.NewBlob(ctx, "app-blob", &storage.BlobArgs{
			AccountName:       account.Name,
			ResourceGroupName: resourceGroup.Name,
			ContainerName:     appContainer.Name,
			Source:            buildResult.Archive,
		})
		if err != nil {
			return err
		}

		// Create a shared access signature to give the Function App access to the code.
		signature := storage.ListStorageAccountServiceSASOutput(ctx, storage.ListStorageAccountServiceSASOutputArgs{
			ResourceGroupName:      resourceGroup.Name,
			AccountName:            account.Name,
			Protocols:              storage.HttpProtocolHttps,
			SharedAccessStartTime:  pulumi.String("2022-01-01"),
			SharedAccessExpiryTime: pulumi.String("2030-01-01"),
			Resource:               pulumi.String("c"),
			Permissions:            pulumi.String("r"),
			ContentType:            pulumi.String("application/json"),
			CacheControl:           pulumi.String("max-age=5"),
			ContentDisposition:     pulumi.String("inline"),
			ContentEncoding:        pulumi.String("deflate"),
			CanonicalizedResource:  pulumi.Sprintf("/blob/%s/%s", account.Name, appContainer.Name),
		})

		// Create an App Service plan for the Function App.
		plan, err := web.NewAppServicePlan(ctx, "plan", &web.AppServicePlanArgs{
			ResourceGroupName: resourceGroup.Name,
			Kind:              pulumi.String("Linux"),
			Reserved:          pulumi.Bool(true),
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
						Value: pulumi.String("custom"),
					},
					&web.NameValuePairArgs{
						Name:  pulumi.String("FUNCTIONS_EXTENSION_VERSION"),
						Value: pulumi.String("~3"),
					},
					&web.NameValuePairArgs{
						Name: pulumi.String("WEBSITE_RUN_FROM_PACKAGE"),
						Value: pulumi.All(account.Name, appContainer.Name, appBlob.Name, signature.ServiceSasToken()).ApplyT(
							func(args []interface{}) string {
								accountName := args[0].(string)
								containerName := args[1].(string)
								blobName := args[2].(string)
								token := args[3].(string)
								return fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", accountName, containerName, blobName, token)
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
		ctx.Export("siteURL", account.PrimaryEndpoints.Web())
		ctx.Export("apiURL", pulumi.Sprintf("https://%v/api", app.DefaultHostName))

		return nil
	})
}
