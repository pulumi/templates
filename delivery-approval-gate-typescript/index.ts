import * as delivery from "@pulumi/delivery";
import * as pulumi from "@pulumi/pulumi";

// A release pipeline: test in dev, then hold the release at an approval gate
// before it updates prod.
const config = new pulumi.Config();
const approvers = config.requireObject<string[]>("approvers");
const approvalCount = config.getNumber("approvalCount") ?? 1;

const pipeline = new delivery.Pipeline("shop");

const repo = new delivery.GitHubSource("app", {
    repository: "acme/shop",
    branch: "main",
    trigger: { paths: ["app/**"] },
});

// Stage 1: run the tests. The stage only passes if the job passes.
const dev = new delivery.Stage("dev", { pipeline });

new delivery.Job("test", {
    stage: dev,
    source: repo,
    directory: "app",
    command: "npm ci && npm test",
});

new delivery.Stacks("dev", {
    stage: dev,
    stacks: { app: { stack: "shop-app/dev", source: repo, directory: "app" } },
});

// Stage 2: the release waits here until enough approvers allow it.
const prod = new delivery.Stage("prod", {
    pipeline,
    from: [dev],
    when: delivery.approval({ approvers, count: approvalCount }),
});

new delivery.Stacks("prod", {
    stage: prod,
    stacks: { app: { stack: "shop-app/prod", source: repo, directory: "app" } },
});
