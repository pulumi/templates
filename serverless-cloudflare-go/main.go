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

		// A Workers KV namespace to persist state between requests.
		namespace, err := cloudflare.NewWorkersKvNamespace(ctx, "namespace", &cloudflare.WorkersKvNamespaceArgs{
			AccountId: pulumi.String(accountId),
			Title:     pulumi.String("visit-counter"),
		})
		if err != nil {
			return err
		}

		// A Cloudflare Worker, exposed on its workers.dev subdomain.
		worker, err := cloudflare.NewWorker(ctx, "worker", &cloudflare.WorkerArgs{
			AccountId: pulumi.String(accountId),
			Name:      pulumi.String("my-app"),
			Subdomain: &cloudflare.WorkerSubdomainArgs{
				Enabled: pulumi.Bool(true),
			},
		})
		if err != nil {
			return err
		}

		// A version of the Worker containing the code and its KV binding.
		version, err := cloudflare.NewWorkerVersion(ctx, "version", &cloudflare.WorkerVersionArgs{
			AccountId:         pulumi.String(accountId),
			WorkerId:          worker.ID(),
			CompatibilityDate: pulumi.String("2025-01-01"),
			MainModule:        pulumi.String("worker.js"),
			Bindings: cloudflare.WorkerVersionBindingArray{
				&cloudflare.WorkerVersionBindingArgs{
					Name:        pulumi.String("COUNTER"),
					Type:        pulumi.String("kv_namespace"),
					NamespaceId: namespace.ID(),
				},
			},
			Modules: cloudflare.WorkerVersionModuleArray{
				&cloudflare.WorkerVersionModuleArgs{
					Name:        pulumi.String("worker.js"),
					ContentType: pulumi.String("application/javascript+module"),
					ContentFile: pulumi.String("worker.js"),
				},
			},
		})
		if err != nil {
			return err
		}

		// Deploy the version so it serves all of the application's traffic.
		_, err = cloudflare.NewWorkersDeployment(ctx, "deployment", &cloudflare.WorkersDeploymentArgs{
			AccountId:  pulumi.String(accountId),
			ScriptName: worker.Name,
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

		// Export the application's URL.
		ctx.Export("url", worker.Subdomain.Url())
		return nil
	})
}
