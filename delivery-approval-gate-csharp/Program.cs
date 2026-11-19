using Pulumi;
using Pulumi.Delivery.Delivery;
using Pulumi.Delivery.Delivery.Inputs;

// A release pipeline: test in dev, then hold the release at an approval gate
// before it updates prod.
return await Deployment.RunAsync(() =>
{
    var config = new Config();
    var approverNames = config.RequireObject<string[]>("approvers");
    var approvalCount = config.GetInt32("approvalCount") ?? 1;

    var approvers = new InputList<string>();
    foreach (var name in approverNames)
    {
        approvers.Add(name);
    }

    var repo = new GitHubSource("repo", new GitHubSourceArgs
    {
        Repository = "acme/shop",
        Branch = "main",
        Trigger = new SourceTriggerArgs { Paths = { "app/**", "delivery/**" } },
    });

    var pipeline = new Pipeline("shop");

    // Stage 1: the test job must pass before the release can move on.
    var dev = new Stage("dev", new StageArgs { Pipeline = pipeline });
    new Stacks("dev", new StacksArgs
    {
        Stage = dev,
        Stacks = { ["app"] = new StacksEntryArgs { Stack = "shop-app/dev", Source = repo, Directory = "app" } },
    });
    new Job("test", new JobArgs
    {
        Stage = dev,
        Source = repo,
        Directory = "app",
        Command = "npm ci && npm test",
    });

    // Stage 2: the release waits here until enough approvers allow it.
    var prod = new Stage("prod", new StageArgs
    {
        Pipeline = pipeline,
        From = { dev },
        When = new RuleArgs
        {
            Kind = "approval",
            Approvers = approvers,
            Count = approvalCount,
        },
    });
    new Stacks("prod", new StacksArgs
    {
        Stage = prod,
        Stacks = { ["app"] = new StacksEntryArgs { Stack = "shop-app/prod", Source = repo, Directory = "app" } },
    });
});
