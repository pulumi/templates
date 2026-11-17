terraform {
  required_providers {
    pulumiservice = {
      source = "pulumi/pulumiservice"
    }
  }
}

variable "schedule_cron" {
  description = "When the weekly report runs (five-field cron, default: Mondays at 08:00 UTC)"
  type        = string
  default     = "0 8 * * 1"
}

# Optional overrides. Left unset, the organization is the stack's organization
# and the agent's name is <project>-<stack>.
variable "organization" {
  description = "The Pulumi Cloud organization that owns the agent"
  type        = string
  default     = null
}

variable "agent_name" {
  description = "The agent's name"
  type        = string
  default     = null
}

locals {
  org = coalesce(var.organization, pulumi.organization)

  # Agent names may only contain lowercase letters, digits and hyphens, and
  # are at most 64 characters long.
  name = substr(replace(lower(coalesce(var.agent_name, "${pulumi.project}-${pulumi.stack}")), "/[^a-z0-9-]/", "-"), 0, 64)
}

data "pulumiservice_build_role_permissions" "read-stacks" {
  stack_rules {
    all    = true
    scopes = ["stack:read"]
  }
}

# The role the agent's tasks assume: read every stack, nothing else.
resource "pulumiservice_organization_role" "role" {
  organization_name = local.org
  name              = "${local.name}-reader"
  description       = "Read-only access to stacks for the ${local.name} agent."
  permissions       = data.pulumiservice_build_role_permissions.read-stacks.permissions
}

# The agent: its instructions, and the role its tasks run with.
resource "pulumiservice_api_agents_definition" "agent" {
  org_name            = local.org
  name                = local.name
  display_name        = "Stack Health Reporter"
  purpose             = "Reports failed and stale stacks."
  custom_instructions = file("${path.module}/instructions.md")
  role_id             = pulumiservice_organization_role.role.role_id
  # Tasks act without waiting for a person. Safe here: the role is read-only.
  approval_mode = "auto"
}

# Pulumi Cloud starts a task as the agent on a schedule.
resource "pulumiservice_api_agent_automation_v2" "weekly-report" {
  org_name      = local.org
  name          = "${local.name}-weekly"
  kind          = "scheduled"
  schedule_cron = var.schedule_cron

  task = {
    agentDefinitionId = pulumiservice_api_agents_definition.agent.definition_id
    approvalMode      = "auto"
    toolExecutionMode = "cloud"
    message = {
      type    = "user_message"
      content = "Write this week's stack health report."
    }
  }
}

# Pulumi Cloud starts a task as the agent when an update of a listed stack fails.
resource "pulumiservice_api_agent_automation_v2" "on-failed-update" {
  org_name = local.org
  name     = "${local.name}-on-failure"
  kind     = "triggered"

  trigger = {
    type       = "stack_operation"
    operations = ["update"]
    results    = ["failed"]
    # The stacks to watch. Each must already exist in Pulumi Cloud. This one
    # is the stack that defines the agent; add your own.
    stacks = [{
      routingProject = pulumi.project
      name           = pulumi.stack
    }]
    # Failures that follow within this many seconds join the same task.
    quietPeriodSeconds = 60
  }

  task = {
    agentDefinitionId = pulumiservice_api_agents_definition.agent.definition_id
    approvalMode      = "auto"
    toolExecutionMode = "cloud"
    message = {
      type    = "user_message"
      content = "A stack update just failed. Find out why from the update's logs and summarize the cause in a few lines."
    }
  }
}

output "agent_name" {
  value = local.name
}

output "agent_id" {
  value = pulumiservice_api_agents_definition.agent.definition_id
}

output "role_id" {
  value = pulumiservice_organization_role.role.role_id
}

output "weekly_report_id" {
  value = pulumiservice_api_agent_automation_v2.weekly-report.automation_id
}

output "next_run" {
  value = pulumiservice_api_agent_automation_v2.weekly-report.next_execution
}

output "on_failed_update_id" {
  value = pulumiservice_api_agent_automation_v2.on-failed-update.automation_id
}
