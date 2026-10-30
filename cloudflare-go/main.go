package main

import (
	"regexp"
	"strings"

	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Import the program's configuration settings.
		cfg := config.New(ctx, "")
		accountId := cfg.Require("accountId")

		name := strings.ToLower(ctx.Project() + "-" + ctx.Stack())
		name = strings.Trim(regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(name, "-"), "-")

		// Create a Cloudflare resource (Workers KV namespace).
		namespace, err := cloudflare.NewWorkersKvNamespace(ctx, "my-namespace", &cloudflare.WorkersKvNamespaceArgs{
			AccountId: pulumi.String(accountId),
			Title:     pulumi.String(name),
		})
		if err != nil {
			return err
		}

		// Export the ID of the namespace.
		ctx.Export("namespaceId", namespace.ID())
		return nil
	})
}
