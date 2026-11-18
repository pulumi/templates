import pulumi
import pulumi_delivery as delivery

# A release pipeline: test in dev, then hold the release at an approval gate
# before it updates prod.
config = pulumi.Config()
approvers = config.require_object("approvers")
approval_count = config.get_int("approvalCount") or 1

repo = delivery.GitHubSource(
    "repo",
    repository="acme/shop",
    branch="main",
    trigger={"paths": ["app/**", "delivery/**"]},
)

pipeline = delivery.Pipeline("shop")

# Stage 1: the test job must pass before the release can move on.
dev = delivery.Stage("dev", pipeline=pipeline)
delivery.Stacks(
    "dev",
    stage=dev,
    stacks={"app": {"stack": "shop-app/dev", "source": repo, "directory": "app"}},
)
delivery.Job(
    "test",
    stage=dev,
    source=repo,
    directory="app",
    command="npm ci && npm test",
)

# Stage 2: the release waits here until enough approvers allow it.
prod = delivery.Stage(
    "prod",
    pipeline=pipeline,
    from_=[dev],
    when=delivery.RuleArgs(kind="approval", approvers=approvers, count=approval_count),
)
delivery.Stacks(
    "prod",
    stage=prod,
    stacks={"app": {"stack": "shop-app/prod", "source": repo, "directory": "app"}},
)
