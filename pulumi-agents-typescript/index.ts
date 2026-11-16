import * as fs from "fs";
import * as pulumi from "@pulumi/pulumi";
import * as pulumiservice from "@pulumi/pulumiservice";

const config = new pulumi.Config();
const scheduleCron = config.get("scheduleCron") || "0 8 * * 1";

// Optional overrides. Left empty, the organization is the stack's organization
// and the agent's name is <project>-<stack>.
const org = config.get("organization") || pulumi.getOrganization();

// Agent names may only contain lowercase letters, digits and hyphens, and
// are at most 64 characters long.
const name = (config.get("agentName") || `${pulumi.getProject()}-${pulumi.getStack()}`)
    .toLowerCase()
    .replace(/[^a-z0-9-]/g, "-")
    .slice(0, 64);

// The role the agent's tasks assume: read every stack, nothing else.
const role = new pulumiservice.OrganizationRole("role", {
    organizationName: org,
    name: `${name}-reader`,
    description: `Read-only access to stacks for the ${name} agent.`,
    permissions: pulumiservice.buildRolePermissionsOutput({
        stackRules: [{ all: true, scopes: ["stack:read"] }],
    }).permissions,
});

// The agent: its instructions, and the role its tasks run with.
const agent = new pulumiservice.api.agents.Definition("agent", {
    orgName: org,
    name: name,
    displayName: "Stack Health Reporter",
    purpose: "Reports failed and stale stacks.",
    customInstructions: fs.readFileSync("./instructions.md", "utf8"),
    roleId: role.roleId,
    // Tasks act without waiting for a person. Safe here: the role is read-only.
    approvalMode: "auto",
});

// Pulumi Cloud starts a task as the agent on a schedule.
const weeklyReport = new pulumiservice.api.AgentAutomationV2("weeklyReport", {
    orgName: org,
    name: `${name}-weekly`,
    kind: "scheduled",
    scheduleCron: scheduleCron,
    task: {
        agentDefinitionId: agent.definitionID,
        approvalMode: "auto",
        toolExecutionMode: "cloud",
        message: {
            type: "user_message",
            content: "Write this week's stack health report.",
        },
    },
});

// Pulumi Cloud starts a task as the agent when an update of a listed stack fails.
const onFailedUpdate = new pulumiservice.api.AgentAutomationV2("onFailedUpdate", {
    orgName: org,
    name: `${name}-on-failure`,
    kind: "triggered",
    trigger: {
        type: "stack_operation",
        operations: ["update"],
        results: ["failed"],
        // The stacks to watch. Each must already exist in Pulumi Cloud. This one
        // is the stack that defines the agent; add your own.
        stacks: [{ routingProject: pulumi.getProject(), name: pulumi.getStack() }],
        // Failures that follow within this many seconds join the same task.
        quietPeriodSeconds: 60,
    },
    task: {
        agentDefinitionId: agent.definitionID,
        approvalMode: "auto",
        toolExecutionMode: "cloud",
        message: {
            type: "user_message",
            content: "A stack update just failed. Find out why from the update's logs and summarize the cause in a few lines.",
        },
    },
});

export const agentName = name;
export const agentId = agent.definitionID;
export const roleId = role.roleId;
export const weeklyReportId = weeklyReport.automationID;
export const nextRun = weeklyReport.nextExecution;
export const onFailedUpdateId = onFailedUpdate.automationID;
