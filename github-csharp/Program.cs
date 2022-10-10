using Pulumi;
using Pulumi.Github;
using System.Collections.Generic;

return await Deployment.RunAsync(() =>
{
    // Create a GitHub Repository
    var repository = new Repository("demo-repo", new RepositoryArgs
    {
        Description = "Demo Repository for GitHub",
    });

    // Export the name of the repository
    return new Dictionary<string, object?>
    {
        ["repositoryName"] = repository.Name
    };
});
