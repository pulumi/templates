"""A Pulumi Cloud agent, its role and its automations, defined as code"""

import re
from pathlib import Path

import pulumi
import pulumi_pulumiservice as pulumiservice

config = pulumi.Config()
schedule_cron = config.get("scheduleCron") or "0 8 * * 1"

# Optional overrides. Left empty, the organization is the stack's organization
# and the agent's name is <project>-<stack>.
org = config.get("organization") or pulumi.get_organization()

# Agent names may only contain lowercase letters, digits and hyphens, and
# are at most 64 characters long.
name = config.get("agentName") or f"{pulumi.get_project()}-{pulumi.get_stack()}"
name = re.sub(r"[^a-z0-9-]", "-", name.lower())[:64]

# The role the agent's tasks assume: read every stack, nothing else.
role = pulumiservice.OrganizationRole(
    "role",
    organization_name=org,
    name=f"{name}-reader",
    description=f"Read-only access to stacks for the {name} agent.",
    permissions=pulumiservice.build_role_permissions_output(
        stack_rules=[{"all": True, "scopes": ["stack:read"]}],
    ).permissions,
)

# The agent: its instructions, and the role its tasks run with.
agent = pulumiservice.api.agents.Definition(
    "agent",
    org_name=org,
    name=name,
    display_name="Stack Health Reporter",
    purpose="Reports failed and stale stacks.",
    custom_instructions=Path("instructions.md").read_text(encoding="utf-8"),
    role_id=role.role_id,
    # Tasks act without waiting for a person. Safe here: the role is read-only.
    approval_mode="auto",
)

# Pulumi Cloud starts a task as the agent on a schedule.
weekly_report = pulumiservice.api.AgentAutomationV2(
    "weeklyReport",
    org_name=org,
    name=f"{name}-weekly",
    kind="scheduled",
    schedule_cron=schedule_cron,
    task={
        "agentDefinitionId": agent.definition_id,
        "approvalMode": "auto",
        "toolExecutionMode": "cloud",
        "message": {
            "type": "user_message",
            "content": "Write this week's stack health report.",
        },
    },
)

# Pulumi Cloud starts a task as the agent when an update of a listed stack fails.
on_failed_update = pulumiservice.api.AgentAutomationV2(
    "onFailedUpdate",
    org_name=org,
    name=f"{name}-on-failure",
    kind="triggered",
    trigger={
        "type": "stack_operation",
        "operations": ["update"],
        "results": ["failed"],
        # The stacks to watch. Each must already exist in Pulumi Cloud. This one
        # is the stack that defines the agent; add your own.
        "stacks": [{"routingProject": pulumi.get_project(), "name": pulumi.get_stack()}],
        # Failures that follow within this many seconds join the same task.
        "quietPeriodSeconds": 60,
    },
    task={
        "agentDefinitionId": agent.definition_id,
        "approvalMode": "auto",
        "toolExecutionMode": "cloud",
        "message": {
            "type": "user_message",
            "content": "A stack update just failed. Find out why from the update's logs and summarize the cause in a few lines.",
        },
    },
)

pulumi.export("agentName", name)
pulumi.export("agentId", agent.definition_id)
pulumi.export("roleId", role.role_id)
pulumi.export("weeklyReportId", weekly_report.automation_id)
pulumi.export("nextRun", weekly_report.next_execution)
pulumi.export("onFailedUpdateId", on_failed_update.automation_id)
