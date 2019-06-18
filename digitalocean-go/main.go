package main

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a DigitalOcean resource (Domain)
		domain, err := digitalocean.NewDomain(ctx, "my-test-domain", nil)
		if err != nil {
			return err
		}

		// Export the name of the domain
		ctx.Export("domainName", domain.Name())
		return nil
	})
}
