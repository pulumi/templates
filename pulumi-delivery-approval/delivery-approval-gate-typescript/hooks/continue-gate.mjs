#!/usr/bin/env node
// Deployment hook: continues a Pulumi Delivery gate through the REST API.
//
// Usage:
//   PULUMI_ACCESS_TOKEN=... node hooks/continue-gate.mjs <org> <project> <stack> <run-number> [stage] [condition]
//
// Steps (from the Delivery "Continuing a gate with the REST API" docs):
//   1. GET the run.
//   2. Find the stage whose gate is open, and take its gate id.
//   3. POST .../gates/<gate-id>/continue with {} or {"condition": "<rule name>"}.

import { pathToFileURL } from "node:url";

const DEFAULT_API = "https://api.pulumi.com";

export async function continueGate({
    org,
    project,
    stack,
    runNumber,
    stage,
    condition,
    token,
    api = DEFAULT_API,
    fetchImpl = fetch,
}) {
    if (!token) throw new Error("PULUMI_ACCESS_TOKEN is required");

    const runUrl = `${api}/api/preview/stacks/${org}/${project}/${stack}/delivery/runs/${runNumber}`;
    const headers = {
        Authorization: `token ${token}`,
        "Content-Type": "application/json",
    };

    const runRes = await fetchImpl(runUrl, { headers });
    if (!runRes.ok) {
        throw new Error(`reading run ${runNumber} failed: HTTP ${runRes.status}`);
    }
    const run = await runRes.json();

    const open = (run.stages ?? []).find(
        (s) => s.gate?.status === "open" && (!stage || s.name === stage),
    );
    if (!open) {
        return { continued: false, reason: "no open gate matched" };
    }

    const body = condition ? { condition } : {};
    const res = await fetchImpl(`${runUrl}/gates/${open.gate.id}/continue`, {
        method: "POST",
        headers,
        body: JSON.stringify(body),
    });
    if (!res.ok) {
        throw new Error(`continuing gate ${open.gate.id} failed: HTTP ${res.status}`);
    }
    return { continued: true, stage: open.name, gateId: open.gate.id };
}

// CLI entry point.
if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
    const [org, project, stack, runNumber, stage, condition] = process.argv.slice(2);
    if (!org || !project || !stack || !runNumber) {
        console.error("usage: continue-gate.mjs <org> <project> <stack> <run-number> [stage] [condition]");
        process.exit(2);
    }
    try {
        const result = await continueGate({
            org,
            project,
            stack,
            runNumber,
            stage: stage || undefined,
            condition: condition || undefined,
            token: process.env.PULUMI_ACCESS_TOKEN,
            api: process.env.PULUMI_API_URL || DEFAULT_API,
        });
        console.log(JSON.stringify(result));
    } catch (err) {
        console.error(err.message);
        process.exit(1);
    }
}
