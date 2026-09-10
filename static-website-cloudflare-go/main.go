package main

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		accountId := cfg.Require("accountId")
		path := cfg.Get("path")
		if path == "" {
			path = "./www"
		}

		// A Cloudflare Worker, exposed on its workers.dev subdomain.
		site, err := cloudflare.NewWorker(ctx, "site", &cloudflare.WorkerArgs{
			AccountId: pulumi.String(accountId),
			Name:      pulumi.String("my-site"),
			Subdomain: &cloudflare.WorkerSubdomainArgs{
				Enabled: pulumi.Bool(true),
			},
		})
		if err != nil {
			return err
		}

		// A version of the Worker that serves the website's files as static assets.
		version, err := cloudflare.NewWorkerVersion(ctx, "version", &cloudflare.WorkerVersionArgs{
			AccountId:         pulumi.String(accountId),
			WorkerId:          site.ID(),
			CompatibilityDate: pulumi.String("2025-01-01"),
			Assets: &cloudflare.WorkerVersionAssetsArgs{
				Directory: pulumi.String(path),
				Config: &cloudflare.WorkerVersionAssetsConfigArgs{
					// Serve /404.html when a request doesn't match a file.
					NotFoundHandling: pulumi.String("404-page"),
					HtmlHandling:     pulumi.String("auto-trailing-slash"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Deploy the version so it serves all of the site's traffic.
		_, err = cloudflare.NewWorkersDeployment(ctx, "deployment", &cloudflare.WorkersDeploymentArgs{
			AccountId:  pulumi.String(accountId),
			ScriptName: site.Name,
			Strategy:   pulumi.String("percentage"),
			Versions: cloudflare.WorkersDeploymentVersionArray{
				&cloudflare.WorkersDeploymentVersionArgs{
					VersionId:  version.ID(),
					Percentage: pulumi.Float64(100),
				},
			},
		})
		if err != nil {
			return err
		}

		// Export the website's URL.
		ctx.Export("url", site.Subdomain.Url())
		return nil
	})
}
