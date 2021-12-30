using Pulumi;
using Pulumi.Github;

class MyStack : Stack
{
    public MyStack()
    {
        // Create a GitHub Repository
        var repository = new Repository("demo-repo", , new Github.RepositoryArgs
        {
            Description = "Demo Repository for GitHub",
        });

        // Export the name of the bucket
        this.RepositoryName = repository.Name;
    }

    [Output]
    public Output<string> RepositoryName { get; set; }
}
