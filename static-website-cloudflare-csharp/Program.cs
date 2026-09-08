using System.Collections.Generic;
using Pulumi;
using Cloudflare = Pulumi.Cloudflare;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var accountId = config.Require("accountId");
    var path = config.Get("path") ?? "./www";

    // A Cloudflare Worker, exposed on its workers.dev subdomain.
    var site = new Cloudflare.Worker("site", new()
    {
        AccountId = accountId,
        Name = "my-site",
        Subdomain = new Cloudflare.Inputs.WorkerSubdomainArgs
        {
            Enabled = true,
        },
    });

    // A version of the Worker that serves the website's files as static assets.
    var version = new Cloudflare.WorkerVersion("version", new()
    {
        AccountId = accountId,
        WorkerId = site.Id,
        CompatibilityDate = "2025-01-01",
        Assets = new Cloudflare.Inputs.WorkerVersionAssetsArgs
        {
            Directory = path,
            Config = new Cloudflare.Inputs.WorkerVersionAssetsConfigArgs
            {
                // Serve /404.html when a request doesn't match a file.
                NotFoundHandling = "404-page",
                HtmlHandling = "auto-trailing-slash",
            },
        },
    });

    // Deploy the version so it serves all of the site's traffic.
    var deployment = new Cloudflare.WorkersDeployment("deployment", new()
    {
        AccountId = accountId,
        ScriptName = site.Name,
        Strategy = "percentage",
        Versions =
        {
            new Cloudflare.Inputs.WorkersDeploymentVersionArgs
            {
                VersionId = version.Id,
                Percentage = 100,
            },
        },
    });

    // Export the website's URL.
    return new Dictionary<string, object?>
    {
        ["url"] = site.Subdomain.Apply(s => s?.Url),
    };
});
