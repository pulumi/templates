using Pulumi;
using Pulumi.Auth0;

class MyStack : Stack
{
    public MyStack()
    {
        // Create an Auth0 Client
        var myClient = new Auth0.Client("client", new Auth0.ClientArgs
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
            JwtConfiguration = new Auth0.Inputs.ClientJwtConfigurationArgs
            {
                Alg = "RS256",
            },
        });

        // Export Client ID and Secret
        this.ClientID = myClient.ClientId;
        this.ClientSecret = myClient.ClientSecret;
    }

    [Output]
    public Output<string> ClientID { get; set; }
    public Output<string> ClientSecret { get; set; }
}
