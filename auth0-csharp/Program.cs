using Pulumi;
using Pulumi.Auth0;
using Pulumi.Auth0.Inputs;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // create Auth0 client
    var auth0Client = new Client("client", new ClientArgs
    {
        AllowedLogoutUrls =
        {
            "https://example.com/logout",
        },
        AllowedOrigins =
        {
            "https://example.com",
        },
        AppType = "regular_web",
        Callbacks =
        {
            "https://example.com/auth/callback",
        },
        JwtConfiguration = new ClientJwtConfigurationArgs
        {
            Alg = "RS256",
        },
    });

    // Export Client ID and Secret
    return new Dictionary<string, object?>
    {
        ["clientID"] = auth0Client.ClientId,
        ["clientSecret"] = auth0Client.ClientSecret
    };
});
