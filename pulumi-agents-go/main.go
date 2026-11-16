package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/pulumi/pulumi-pulumiservice/sdk/go/pulumiservice"
	"github.com/pulumi/pulumi-pulumiservice/sdk/go/pulumiservice/api"
	"github.com/pulumi/pulumi-pulumiservice/sdk/go/pulumiservice/api/agents"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		scheduleCron := cfg.Get("scheduleCron")
		if scheduleCron == "" {
			scheduleCron = "0 8 * * 1"
		}

		// Optional overrides. Left empty, the organization is the stack's organization
		// and the agent's name is <project>-<stack>.
		org := cfg.Get("organization")
		if org == "" {
			org = ctx.Organization()
		}
		name := cfg.Get("agentName")
		if name == "" {
			name = ctx.Project() + "-" + ctx.Stack()
		}

		// Agent names may only contain lowercase letters, digits and hyphens, and
		// are at most 64 characters long.
		name = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(strings.ToLower(name), "-")
		if len(name) > 64 {
			name = name[:64]
		}

		// The role the agent's tasks assume: read every stack, nothing else.
		permissions, err := pulumiservice.BuildRolePermissions(ctx, &pulumiservice.BuildRolePermissionsArgs{
			StackRules: []pulumiservice.RoleStackRule{{
				All:    pulumi.BoolRef(true),
				Scopes: []pulumiservice.RbacStackScope{pulumiservice.RbacStackScopeStackRead},
			}},
		})
		if err != nil {
			return err
		}
		role, err := pulumiservice.NewOrganizationRole(ctx, "role", &pulumiservice.OrganizationRoleArgs{
			OrganizationName: pulumi.String(org),
			Name:             pulumi.String(name + "-reader"),
			Description:      pulumi.String("Read-only access to stacks for the " + name + " agent."),
			Permissions:      pulumi.ToMap(permissions.Permissions),
		})
		if err != nil {
			return err
		}

		// The agent: its instructions, and the role its tasks run with.
		instructions, err := os.ReadFile("instructions.md")
		if err != nil {
			return fmt.Errorf("reading instructions.md: %w", err)
		}
		agent, err := agents.NewDefinition(ctx, "agent", &agents.DefinitionArgs{
			OrgName:            pulumi.String(org),
			Name:               pulumi.String(name),
			DisplayName:        pulumi.String("Stack Health Reporter"),
			Purpose:            pulumi.String("Reports failed and stale stacks."),
			CustomInstructions: pulumi.String(string(instructions)),
			RoleId:             role.RoleId,
			// Tasks act without waiting for a person. Safe here: the role is read-only.
			ApprovalMode: pulumi.String("auto"),
		})
		if err != nil {
			return err
		}

		// Pulumi Cloud starts a task as the agent on a schedule.
		weeklyReport, err := api.NewAgentAutomationV2(ctx, "weeklyReport", &api.AgentAutomationV2Args{
			OrgName:      pulumi.String(org),
			Name:         pulumi.String(name + "-weekly"),
			Kind:         pulumi.String("scheduled"),
			ScheduleCron: pulumi.String(scheduleCron),
			Task: pulumi.Map{
				"agentDefinitionId": agent.DefinitionID,
				"approvalMode":      pulumi.String("auto"),
				"toolExecutionMode": pulumi.String("cloud"),
				"message": pulumi.Map{
					"type":    pulumi.String("user_message"),
					"content": pulumi.String("Write this week's stack health report."),
				},
			},
		})
		if err != nil {
			return err
		}

		// Pulumi Cloud starts a task as the agent when an update of a listed stack fails.
		onFailedUpdate, err := api.NewAgentAutomationV2(ctx, "onFailedUpdate", &api.AgentAutomationV2Args{
			OrgName: pulumi.String(org),
			Name:    pulumi.String(name + "-on-failure"),
			Kind:    pulumi.String("triggered"),
			Trigger: pulumi.Map{
				"type":       pulumi.String("stack_operation"),
				"operations": pulumi.StringArray{pulumi.String("update")},
				"results":    pulumi.StringArray{pulumi.String("failed")},
				// The stacks to watch. Each must already exist in Pulumi Cloud. This one
				// is the stack that defines the agent; add your own.
				"stacks": pulumi.Array{
					pulumi.Map{
						"routingProject": pulumi.String(ctx.Project()),
						"name":           pulumi.String(ctx.Stack()),
					},
				},
				// Failures that follow within this many seconds join the same task.
				"quietPeriodSeconds": pulumi.Int(60),
			},
			Task: pulumi.Map{
				"agentDefinitionId": agent.DefinitionID,
				"approvalMode":      pulumi.String("auto"),
				"toolExecutionMode": pulumi.String("cloud"),
				"message": pulumi.Map{
					"type":    pulumi.String("user_message"),
					"content": pulumi.String("A stack update just failed. Find out why from the update's logs and summarize the cause in a few lines."),
				},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("agentName", pulumi.String(name))
		ctx.Export("agentId", agent.DefinitionID)
		ctx.Export("roleId", role.RoleId)
		ctx.Export("weeklyReportId", weeklyReport.AutomationID)
		ctx.Export("nextRun", weeklyReport.NextExecution)
		ctx.Export("onFailedUpdateId", onFailedUpdate.AutomationID)
		return nil
	})
}
