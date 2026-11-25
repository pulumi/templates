terraform {
  required_providers {
    delivery = {
      source = "pulumi/delivery"
    }
  }
}

variable "approvers" {
  description = "Pulumi usernames or team names that must approve prod releases"
  type        = list(string)
  default     = ["your-pulumi-username"]
}

variable "approvalCount" {
  description = "Number of approvals required before prod runs"
  type        = number
  default     = 1
}

# A release pipeline: test in dev, then hold the release at an approval gate
# before it updates prod.
resource "delivery_git_hub_source" "repo" {
  repository = "acme/shop"
  branch     = "main"
  trigger = {
    paths = ["app/**", "delivery/**"]
  }
}

resource "delivery_pipeline" "shop" {}

# Stage 1: the test job must pass before the release can move on.
resource "delivery_stage" "dev" {
  pipeline = delivery_pipeline.shop
}

resource "delivery_stacks" "dev" {
  stage = delivery_stage.dev
  stacks = {
    app = { stack = "shop-app/dev", source = delivery_git_hub_source.repo, directory = "app" }
  }
}

resource "delivery_job" "test" {
  stage     = delivery_stage.dev
  source    = delivery_git_hub_source.repo
  directory = "app"
  command   = "npm ci && npm test"
}

# Stage 2: the release waits here until enough approvers allow it.
resource "delivery_stage" "prod" {
  pipeline = delivery_pipeline.shop
  from     = [delivery_stage.dev]
  when = {
    kind      = "approval"
    approvers = var.approvers
    count     = var.approvalCount
  }
}

resource "delivery_stacks" "prod" {
  stage = delivery_stage.prod
  stacks = {
    app = { stack = "shop-app/prod", source = delivery_git_hub_source.repo, directory = "app" }
  }
}
