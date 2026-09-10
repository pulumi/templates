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

		// Create a Cloudflare resource (Workers KV namespace).
		namespace, err := cloudflare.NewWorkersKvNamespace(ctx, "my-namespace", &cloudflare.WorkersKvNamespaceArgs{
			AccountId: pulumi.String(accountId),
			Title:     pulumi.String("my-namespace"),
		})
		if err != nil {
			return err
		}

		// Export the ID of the namespace.
		ctx.Export("namespaceId", namespace.ID())
		return nil
	})
}
