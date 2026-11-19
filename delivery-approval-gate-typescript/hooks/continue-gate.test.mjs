import { test } from "node:test";
import assert from "node:assert/strict";
import http from "node:http";
import { continueGate } from "./continue-gate.mjs";

// Starts a fake Pulumi API that serves one run and records the calls made to it.
function startFakeApi(run) {
    const calls = [];
    const server = http.createServer((req, res) => {
        let body = "";
        req.on("data", (c) => (body += c));
        req.on("end", () => {
            calls.push({ method: req.method, url: req.url, auth: req.headers.authorization, body });
            if (req.method === "GET" && req.url.endsWith("/delivery/runs/7")) {
                res.writeHead(200, { "Content-Type": "application/json" });
                return res.end(JSON.stringify(run));
            }
            if (req.method === "POST" && req.url.endsWith("/gates/gate-prod/continue")) {
                res.writeHead(200, { "Content-Type": "application/json" });
                return res.end("{}");
            }
            res.writeHead(404);
            res.end();
        });
    });
    return new Promise((resolve) => {
        server.listen(0, "127.0.0.1", () => {
            const { port } = server.address();
            resolve({ api: `http://127.0.0.1:${port}`, calls, close: () => server.close() });
        });
    });
}

const run = {
    stages: [
        { name: "dev", gate: null },
        { name: "prod", gate: { id: "gate-prod", status: "open" } },
    ],
};

test("continues the open gate with an empty body", async () => {
    const fake = await startFakeApi(run);
    try {
        const result = await continueGate({
            org: "acme", project: "shop-delivery", stack: "main", runNumber: "7",
            token: "test-token", api: fake.api,
        });
        assert.deepEqual(result, { continued: true, stage: "prod", gateId: "gate-prod" });
        const post = fake.calls.find((c) => c.method === "POST");
        assert.equal(post.url, "/api/preview/stacks/acme/shop-delivery/main/delivery/runs/7/gates/gate-prod/continue");
        assert.equal(post.body, "{}");
        assert.equal(post.auth, "token test-token");
    } finally {
        fake.close();
    }
});

test("sends the rule name as condition for multi-rule gates", async () => {
    const fake = await startFakeApi(run);
    try {
        await continueGate({
            org: "acme", project: "shop-delivery", stack: "main", runNumber: "7",
            condition: "dashboards", token: "t", api: fake.api,
        });
        assert.equal(fake.calls.find((c) => c.method === "POST").body, '{"condition":"dashboards"}');
    } finally {
        fake.close();
    }
});

test("does nothing when no gate is open for the requested stage", async () => {
    const fake = await startFakeApi(run);
    try {
        const result = await continueGate({
            org: "acme", project: "shop-delivery", stack: "main", runNumber: "7",
            stage: "dev", token: "t", api: fake.api,
        });
        assert.equal(result.continued, false);
        assert.equal(fake.calls.some((c) => c.method === "POST"), false);
    } finally {
        fake.close();
    }
});

test("throws when the token is missing", async () => {
    await assert.rejects(
        continueGate({ org: "a", project: "b", stack: "c", runNumber: "1", token: "" }),
        /PULUMI_ACCESS_TOKEN is required/,
    );
});
