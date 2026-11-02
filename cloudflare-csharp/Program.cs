using System.Collections.Generic;
using System.Text.RegularExpressions;
using Pulumi;
using Cloudflare = Pulumi.Cloudflare;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var accountId = config.Require("accountId");

    var name = Regex.Replace(
        $"{Deployment.Instance.ProjectName}-{Deployment.Instance.StackName}".ToLowerInvariant(),
        "[^a-z0-9-]+", "-").Trim('-');

    // Create a Cloudflare resource (Workers KV namespace).
    var ns = new Cloudflare.WorkersKvNamespace("my-namespace", new()
    {
        AccountId = accountId,
        Title = name,
    });

    // Export the ID of the namespace.
    return new Dictionary<string, object?>
    {
        ["namespaceId"] = ns.Id,
    };
});
