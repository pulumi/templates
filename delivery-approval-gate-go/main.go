package main

import (
	"encoding/json"

	"github.com/pulumi/pulumi-delivery/sdk/go/delivery"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// A release pipeline: test in dev, then hold the release at an approval gate
// before it updates prod.
func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")

		var names []string
		if err := json.Unmarshal([]byte(cfg.Require("approvers")), &names); err != nil {
			return err
		}
		approvers := pulumi.StringArray{}
		for _, name := range names {
			approvers = append(approvers, pulumi.String(name))
		}
		approvalCount := cfg.GetInt("approvalCount")
		if approvalCount == 0 {
			approvalCount = 1
		}

		repo, err := delivery.NewGitHubSource(ctx, "repo", &delivery.GitHubSourceArgs{
			Repository: pulumi.String("acme/shop"),
			Branch:     pulumi.String("main"),
			Trigger: delivery.SourceTriggerArgs{
				Paths: pulumi.StringArray{pulumi.String("app/**"), pulumi.String("delivery/**")},
			},
		})
		if err != nil {
			return err
		}

		pipeline, err := delivery.NewPipeline(ctx, "shop", nil)
		if err != nil {
			return err
		}

		// Stage 1: the test job must pass before the release can move on.
		dev, err := delivery.NewStage(ctx, "dev", &delivery.StageArgs{Pipeline: pipeline})
		if err != nil {
			return err
		}
		if _, err := delivery.NewStacks(ctx, "dev", &delivery.StacksArgs{
			Stage: dev,
			Stacks: delivery.StacksEntryMap{
				"app": delivery.StacksEntryArgs{Stack: pulumi.String("shop-app/dev"), Source: repo, Directory: pulumi.String("app")},
			},
		}); err != nil {
			return err
		}
		if _, err := delivery.NewJob(ctx, "test", &delivery.JobArgs{
			Stage:     dev,
			Source:    repo,
			Directory: pulumi.String("app"),
			Command:   pulumi.String("npm ci && npm test"),
		}); err != nil {
			return err
		}

		// Stage 2: the release waits here until enough approvers allow it.
		prod, err := delivery.NewStage(ctx, "prod", &delivery.StageArgs{
			Pipeline: pipeline,
			From:     delivery.StageArray{dev},
			When: delivery.RuleArgs{
				Kind:      pulumi.String("approval"),
				Approvers: approvers,
				Count:     pulumi.Int(approvalCount),
			},
		})
		if err != nil {
			return err
		}
		_, err = delivery.NewStacks(ctx, "prod", &delivery.StacksArgs{
			Stage: prod,
			Stacks: delivery.StacksEntryMap{
				"app": delivery.StacksEntryArgs{Stack: pulumi.String("shop-app/prod"), Source: repo, Directory: pulumi.String("app")},
			},
		})
		return err
	})
}
