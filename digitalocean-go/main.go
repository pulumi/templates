package main

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a DigitalOcean resource (Domain)
		domain, err := digitalocean.NewDomain(ctx, "my-domain", &digitalocean.DomainArgs{
			Name: pulumi.String("my-domain.io"),
		})
		if err != nil {
			return err
		}

		// Export the name of the domain
		ctx.Export("domainName", domain.Name)
		return nil
	})
}
