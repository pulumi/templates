# The software factory, delivered as code

This replaces `npm run orchestrate -- up` with a Pulumi Delivery pipeline. Same
five stacks, same order, same post-deploy steps — but Pulumi Cloud runs them on
every push to `main` instead of your laptop running them against your AWS session.

```
                  ┌── cluster ──┐   ┌─ platform ─┐   ┌──────── services ────────┐   ┌─ demo ─┐
push to main ──▶  │  approval   │──▶│            │──▶│  factory ┐               │──▶│ manual │
                  │  gate       │   │  platform  │   │  agent   ├─▶ bootstrap   │   │  seed  │
                  │  llm        │   │            │   │  triage  ┘   └─▶ verify  │   │        │
                  └─────────────┘   └────────────┘   └──────────────────────────┘   └────────┘
```

## What maps to what

| Today | In the pipeline |
| --- | --- |
| `scripts/lib/plan.ts` `UP_ORDER` | Stage `from` chain |
| `orchestrate.ts` stop-on-first-failure | A failed stage stops the run; later stages never start |
| `pulumi up` per stack, serially | `Stacks` entries — three in one group, updating side by side |
| Hand-edited `Pulumi.dev.yaml` stack pointers | `config` on each `Stacks` entry |
| `pulumi config set projectBrief < factory-prompt.md` | `fs.readFileSync("../factory-prompt.md")` at registration |
| `KUBECONFIG=… tsx scripts/bootstrap.ts` | `Job` **bootstrap** |
| `tsx scripts/compat-check.ts` | `Job` **verify** |
| `scripts/seed-issues.sh < examples/items.txt` | `Job` **seed**, behind a manual gate |
| "This can cost ~$2.2/hr. Continue? [y/N]" | `approval()` gate on the `cluster` stage |
| `--only <stack>` | Retry one stage from the run page |

### Why four stages and not one `Stacks` group

All five stacks could go in a single group — they use `StackReference`s, and
Delivery orders a group by its stack references. But a group is **one unit for
retries**: a `triage` failure would re-run the 45-minute GPU bring-up. Stages
give each link its own retry boundary, and let the expensive one sit behind a gate.

`factory`, `agent` and `triage` do share a group. They reference `platform` (and
`triage` also `llm`), never each other, so they update in parallel. That is three
slots of the org concurrency limit.

### Why the gate is `approval()` and not `manual()`

An approval gate shows the approver a **preview of the stack updates**. For a
routine push the preview says the `llm` stack is unchanged and you click through;
for a change that replaces the GPU node group, you see that before it happens.
`manual()` shows no preview. Approval gates need Pro or Enterprise.

## Before this runs

### 1. Target stack config must be readable by the runner

`Pulumi.*.yaml` is git-ignored in this repo, so a fresh checkout on a hosted
runner has no stack config. Two ways out, and the pipeline assumes the second:

- commit the non-secret `Pulumi.<stack>.yaml` files, or
- let the pipeline write the config, which is what the `config` maps above do.

What the pipeline **cannot** write:

- `aws:defaultTags` — `config` takes plain strings only, and this is an object.
  Leave it in the target stack's config file.
- `triage:jevApiKey`, `llm:hfToken` — `config` rejects secrets by design. Put
  them in an ESC environment imported by the target stack.

### 2. AWS credentials come from ESC, not from your shell

`stacks/llm/creds.sh` uses the local `aicr-test` `credential_process` profile.
That does not exist on a runner. Each target stack needs an ESC environment
built on `fn::open::aws-login`, imported in its config file:

```yaml
environment:
  - aws-login/my-sandbox
```

The pipeline's service principal must be allowed to open it, or the update fails
with "you do not have permission to open environment".

### 3. The EKS API endpoint must be reachable from the runner

`bootstrap`, `verify` and `seed` all use `kubectl` against the cluster. The node
group gets public IPs, so the API endpoint is almost certainly public — worth
confirming before the first run. If it is private, those three jobs cannot run
on a Pulumi-managed runner.

## Register it

```bash
cd delivery
pulumi package add delivery
npm install
cp Pulumi.main.yaml.example Pulumi.main.yaml
pulumi stack init my-org/software-factory-delivery/main
git add -A && git commit -m "Add the delivery pipeline" && git push
pulumi up
```

`pulumi up` exits with code 10. That is expected: the stages are waiting for a run.

The five target stacks must already exist — Delivery updates stacks, it does not
create them. `my-org/llm/dev`, `my-org/platform/dev`, `my-org/factory/dev`,
`my-org/agent/dev`, `my-org/triage/dev`.

## Teardown

The pipeline only brings the factory up. `npm run orchestrate -- down --yes`
stays the way to tear it down, and the ~$2.2/hr bill runs until you do.
