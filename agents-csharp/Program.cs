using System.Collections.Generic;
using System.IO;
using System.Text.RegularExpressions;
using Pulumi;
using PulumiService = Pulumi.PulumiService;

return await Deployment.RunAsync(() =>
{
    var config = new Config();
    var scheduleCron = config.Get("scheduleCron");
    if (string.IsNullOrEmpty(scheduleCron))
    {
        scheduleCron = "0 8 * * 1";
    }

    // Optional overrides. Left empty, the organization is the stack's organization
    // and the agent's name is <project>-<stack>.
    var org = config.Get("organization");
    if (string.IsNullOrEmpty(org))
    {
        org = Deployment.Instance.OrganizationName;
    }
    var name = config.Get("agentName");
    if (string.IsNullOrEmpty(name))
    {
        name = $"{Deployment.Instance.ProjectName}-{Deployment.Instance.StackName}";
    }

    // Agent names may only contain lowercase letters, digits and hyphens, and
    // are at most 64 characters long.
    name = Regex.Replace(name.ToLowerInvariant(), "[^a-z0-9-]", "-");
    if (name.Length > 64)
    {
        name = name[..64];
    }

    // The role the agent's tasks assume: read every stack, nothing else.
    var role = new PulumiService.OrganizationRole("role", new()
    {
        OrganizationName = org,
        Name = $"{name}-reader",
        Description = $"Read-only access to stacks for the {name} agent.",
        Permissions = PulumiService.BuildRolePermissions.Invoke(new()
        {
            StackRules =
            {
                new PulumiService.Inputs.RoleStackRuleArgs
                {
                    All = true,
                    Scopes = { PulumiService.RbacStackScope.StackRead },
                },
            },
        }).Apply(result => result.Permissions),
    });

    // The agent: its instructions, and the role its tasks run with.
    var agent = new PulumiService.Api.Agents.Definition("agent", new()
    {
        OrgName = org,
        Name = name,
        DisplayName = "Stack Health Reporter",
        Purpose = "Reports failed and stale stacks.",
        CustomInstructions = File.ReadAllText(Path.Combine(Deployment.Instance.RootDirectory, "instructions.md")),
        RoleId = role.RoleId,
        // Tasks act without waiting for a person. Safe here: the role is read-only.
        ApprovalMode = "auto",
    });

    // Pulumi Cloud starts a task as the agent on a schedule.
    var weeklyReport = new PulumiService.Api.AgentAutomationV2("weeklyReport", new()
    {
        OrgName = org,
        Name = $"{name}-weekly",
        Kind = "scheduled",
        ScheduleCron = scheduleCron,
        Task = new Dictionary<string, object?>
        {
            ["agentDefinitionId"] = agent.DefinitionID,
            ["approvalMode"] = "auto",
            ["toolExecutionMode"] = "cloud",
            ["message"] = new Dictionary<string, object?>
            {
                ["type"] = "user_message",
                ["content"] = "Write this week's stack health report.",
            },
        },
    });

    // Pulumi Cloud starts a task as the agent when an update of a listed stack fails.
    var onFailedUpdate = new PulumiService.Api.AgentAutomationV2("onFailedUpdate", new()
    {
        OrgName = org,
        Name = $"{name}-on-failure",
        Kind = "triggered",
        Trigger = new Dictionary<string, object?>
        {
            ["type"] = "stack_operation",
            ["operations"] = new[] { "update" },
            ["results"] = new[] { "failed" },
            // The stacks to watch. Each must already exist in Pulumi Cloud. This one
            // is the stack that defines the agent; add your own.
            ["stacks"] = new[]
            {
                new Dictionary<string, object?>
                {
                    ["routingProject"] = Deployment.Instance.ProjectName,
                    ["name"] = Deployment.Instance.StackName,
                },
            },
            // Failures that follow within this many seconds join the same task.
            ["quietPeriodSeconds"] = 60,
        },
        Task = new Dictionary<string, object?>
        {
            ["agentDefinitionId"] = agent.DefinitionID,
            ["approvalMode"] = "auto",
            ["toolExecutionMode"] = "cloud",
            ["message"] = new Dictionary<string, object?>
            {
                ["type"] = "user_message",
                ["content"] = "A stack update just failed. Find out why from the update's logs and summarize the cause in a few lines.",
            },
        },
    });

    return new Dictionary<string, object?>
    {
        ["agentName"] = name,
        ["agentId"] = agent.DefinitionID,
        ["roleId"] = role.RoleId,
        ["weeklyReportId"] = weeklyReport.AutomationID,
        ["nextRun"] = weeklyReport.NextExecution,
        ["onFailedUpdateId"] = onFailedUpdate.AutomationID,
    };
});
