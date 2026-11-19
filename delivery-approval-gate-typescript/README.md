# ${PROJECT}

${DESCRIPTION}

A Pulumi Delivery pipeline with two stages:

1. **dev**: a `test` job runs `npm ci && npm test`. Then the `dev` stack updates.
2. **prod**: waits at an approval gate. Only after enough approvers allow the release does the `prod` stack update.

## Configure

Set the approvers and the number of approvals required:

```bash
pulumi config set --path approvers '["platform-leads"]'
pulumi config set approvalCount 2
```

## Continue a gate from a hook

`hooks/continue-gate.mjs` continues an open gate through the Delivery REST API. Use it from a deployment hook or from the command line:

```bash
PULUMI_ACCESS_TOKEN=... node hooks/continue-gate.mjs <org> <project> <stack> <run-number> [stage] [condition]
```

Run the hook's tests with `npm run test:hooks`.

## Note

`@pulumi/delivery` is part of Pulumi Delivery, which is in public preview. Update the package version in `package.json` to the release you have access to.
