# ${PROJECT}

${DESCRIPTION}

A Pulumi Delivery pipeline in HCL with two stages:

1. **dev**: a `test` job runs `npm ci && npm test`, then the `dev` stack updates.
2. **prod**: waits at an approval gate. Once enough approvers allow the release, the `prod` stack updates.

## Configure

```bash
pulumi config set --path approvers '["platform-leads"]'
pulumi config set approvalCount 2
```

## Install and run

No install step: the program is the HCL in `main.tf`.

```bash
pulumi stack init dev
pulumi preview
```

## Continue a gate from a hook

`hooks/continue-gate.mjs` continues an open gate through the Delivery REST API. Use it from a deployment hook or the command line:

```bash
PULUMI_ACCESS_TOKEN=... node hooks/continue-gate.mjs <org> <project> <stack> <run-number> [stage] [condition]
```

Run the hook's tests with `node --test hooks/`.

## Note

Pulumi Delivery is in public preview. Install the Delivery package version you have access to before running this program.
