using System.Collections.Generic;
using System.Text.RegularExpressions;
using Pulumi;
using Cloudflare = Pulumi.Cloudflare;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var accountId = config.Require("accountId");

    // Worker names may only contain lowercase letters, numbers and dashes.
    var name = Regex.Replace(
        $"{Deployment.Instance.ProjectName}-{Deployment.Instance.StackName}".ToLowerInvariant(),
        "[^a-z0-9-]+", "-").Trim('-');

    // A Workers KV namespace to persist state between requests.
    var kv = new Cloudflare.WorkersKvNamespace("namespace", new()
    {
        AccountId = accountId,
        Title = $"{name}-visit-counter",
    });

    // A Cloudflare Worker, exposed on its workers.dev subdomain.
    var worker = new Cloudflare.Worker("worker", new()
    {
        AccountId = accountId,
        Name = name,
        Subdomain = new Cloudflare.Inputs.WorkerSubdomainArgs
        {
            Enabled = true,
        },
    });

    // A version of the Worker containing the code and its KV binding.
    var version = new Cloudflare.WorkerVersion("version", new()
    {
        AccountId = accountId,
        WorkerId = worker.Id,
        CompatibilityDate = "2025-01-01",
        MainModule = "worker.js",
        Bindings =
        {
            new Cloudflare.Inputs.WorkerVersionBindingArgs
            {
                Name = "COUNTER",
                Type = "kv_namespace",
                NamespaceId = kv.Id,
            },
        },
        Modules =
        {
            new Cloudflare.Inputs.WorkerVersionModuleArgs
            {
                Name = "worker.js",
                ContentType = "application/javascript+module",
                ContentFile = "worker.js",
            },
        },
    });

    // Deploy the version so it serves all of the application's traffic.
    var deployment = new Cloudflare.WorkersDeployment("deployment", new()
    {
        AccountId = accountId,
        ScriptName = worker.Name,
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

    // Export the application's URL.
    return new Dictionary<string, object?>
    {
        ["url"] = worker.Subdomain.Apply(s => s?.Url),
    };
});
