using System.Collections.Generic;
using Pulumi;
using Cloudflare = Pulumi.Cloudflare;

return await Deployment.RunAsync(() =>
{
    // Import the program's configuration settings.
    var config = new Config();
    var accountId = config.Require("accountId");

    // Create a Cloudflare resource (Workers KV namespace).
    var ns = new Cloudflare.WorkersKvNamespace("my-namespace", new()
    {
        AccountId = accountId,
        Title = $"{Deployment.Instance.ProjectName}-{Deployment.Instance.StackName}",
    });

    // Export the ID of the namespace.
    return new Dictionary<string, object?>
    {
        ["namespaceId"] = ns.Id,
    };
});
