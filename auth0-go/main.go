package main

import (
	"github.com/pulumi/pulumi-auth0/sdk/v2/go/auth0"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		client, err := auth0.NewClient(ctx, "client", &auth0.ClientArgs{
			AllowedLogoutUrls: pulumi.StringArray{
				pulumi.String("https://example.com/logout"),
			},
			AllowedOrigins: pulumi.StringArray{
				pulumi.String("https://example.com"),
			},
			AppType: pulumi.String("regular_web"),
			Callbacks: pulumi.StringArray{
				pulumi.String("https://example.com/auth/callback"),
			},
			JwtConfiguration: &auth0.ClientJwtConfigurationArgs{
				Alg: pulumi.String("RS256"),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("clientId", client.ClientId)
		ctx.Export("clientSecret", client.ClientSecret)
		return nil
	})
}
