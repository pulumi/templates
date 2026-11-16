# AICR chatbot demo

A public chatbot page whose answers come from an open-weight LLM served by
**NVIDIA Dynamo** on an **AWS EKS** cluster with two L4 GPUs — and whose entire
GPU software stack (GPU Operator, DRA driver, KAI scheduler, Prometheus,
cert-manager, agentgateway, Grove, Dynamo platform) is deployed by **one
`ClusterStack` resource** from the Pulumi NVIDIA AICR provider.


```
internet → ELB → chatbot (FastAPI page)  → Dynamo frontend → 2 × vLLM workers (1 L4 each)
                                            └─ all of it on the AICR-validated stack
```

**Cost: ~$2.2/hr while up (~$53/day).** `pulumi destroy` when done.

## Layout

| Path | What |
|---|---|
| `app/` | The chat application: `main.py` (FastAPI, streams tokens from the OpenAI-compatible endpoint), `static/index.html`, `Dockerfile`, tests |
| `infra/` | The Pulumi program (Python): `cluster.py` → `gpu_stack.py` → `gate.py` → `model.py` → `chatbot.py`, wired in `__main__.py`; scripts; tests |

## Prerequisites

- AWS credentials for the sandbox account (`. infra/creds.sh` uses the `my-aws-profile` profile)
- Pulumi CLI (logged in to Pulumi Cloud), Python 3.12, Docker with buildx, kubectl, helm
- The provider checkout at `../pulumi-nvidia-aicr` — the demo targets the unreleased
  0.4.x line (AICR SDK v0.21.1, Dynamo 1.4.2). `setup.sh` installs its Python SDK
  from that path and its plugin from `bin/`. When 0.4.0 is on PyPI, swap the path in
  `infra/requirements.txt` for `pulumi-labs-nvidia-aicr>=0.4.0,<0.5.0`.

## Bring it up (do this ahead of the demo, not live)

```bash
cd infra
./setup.sh                 # venv, SDK, provider plugin (one-time)
. ./creds.sh               # AWS session; runs `aws login` if expired
./preflight.sh             # quota, AMI, instance offerings, toolchain -- read-only
pulumi stack select dev    # or: pulumi stack init dev
pulumi preview             # ~67 resources, no cost
pulumi up                  # ~30 min: EKS ~12, nodes+driver ~7, AICR stack ~6, model image+download ~8
./verify.sh                # the money shots, in escalating order
```

`pulumi up` ends with `chat_url`. Open it, ask a question.

Config knobs (`pulumi config set …`): `instanceType` (`g6.2xlarge`; fallback
`g5.2xlarge`), `nodeCount` (2), `model` (`Qwen/Qwen2.5-1.5B-Instruct`),
`modelReplicas` (2), `k8sVersion` (1.35), `validate` (false), `hfToken` (secret,
optional), `owner` (tag value — load-bearing, see gotchas).

## The demo flow (~10 minutes)

1. **The program** — `infra/__main__.py`: five layers, read top to bottom. The
   product is step 2, `gpu_stack.py`: five inputs (`h100 / eks / inference /
   dynamo / ubuntu`) resolve NVIDIA's validated recipe in-process and become ~16
   Helm releases in dependency order, each a real Pulumi resource.
2. **The deployment** — Pulumi Console stack view: one `ClusterStack` fanning out
   into namespaces and Helm releases; the `DynamoGraphDeployment` and the chat
   app on top of it. `pulumi stack output`: `recipe_name`, `recipe_version`,
   `deployed_components`.
3. **The chat page** — `pulumi stack output chat_url`. Ask something. Tokens
   stream; the footer shows tokens/s.
4. **Prove the GPUs** — `./verify.sh`: two nodes advertising `nvidia.com/gpu`,
   the Helm inventory, Dynamo pods on two different nodes, `nvidia-smi` inside a
   worker, a raw `curl` against the OpenAI-compatible API.
5. **Both GPUs busy** — `./bench.sh 8` fires eight concurrent questions; it
   prints aggregate tok/s and each worker's peak GPU utilization during the
   run. Measured: ~72 tok/s for a single chat, ~465 tok/s aggregate at 8
   concurrent clients with both L4s at 100%.
6. **Optional closers** — `pulumi config set model Qwen/Qwen2.5-3B-Instruct &&
   pulumi preview` shows the targeted diff (one CR changes); `pulumi config set
   validate true && pulumi up` runs AICR's deployment + conformance checks
   against the live cluster (expect a mixed verdict on L4: the recipe is
   H100-tuned and two components are skipped by design).


## Tear down

```bash
pulumi destroy             # ~10-12 min
./leakcheck.sh             # ELBs, k8s-elb-* security groups, volumes, ECR, roles, ASGs
```

## Tests

```bash
cd infra && ./venv/bin/python -m pytest tests -q     # program under Pulumi mocks + gate script against a fake kubectl
cd app && python3 -m venv .venv && .venv/bin/pip install -r requirements-dev.txt && .venv/bin/python -m pytest tests -q
```

Run the app locally without a GPU: `python tests/stub_llm.py 8000 &` then
`LLM_BASE_URL=http://127.0.0.1:8000 uvicorn main:app --port 8080`.

## Gotchas

- **`Owner` tags are load-bearing** in the sandbox account. The cleanup lambda
  deletes untagged EKS clusters, ASGs, IAM roles, ECR repos, classic ELBs and
  security groups after an hour. The program tags everything it creates, the ASG
  explicitly (cluster tags do not propagate), and the ELB through a Service
  annotation. Don't remove them.
- **AWS sessions expire mid-day**: `. ./creds.sh` before any `pulumi` command.
- **A first `pulumi up` can hit a webhook race** inside the recipe: with
  `skipAwait`, one component's CRs can beat a dependency's admission webhook.
  Known provider issue; the remedy is to re-run `pulumi up`. (Not seen on the
  2026-09-18 deployment.) The model layer is protected by the gate
  (`wait-for-dynamo.sh`), which waits for the Dynamo operator, its CRD and
  webhook, and GPU capacity before creating the `DynamoGraphDeployment`.
- **`InsufficientInstanceCapacity` for g6.2xlarge** happens per AZ (1a and 1d on
  2026-09-18). The ASG spans five AZs and retries, so wait a few minutes before
  switching to `g5.2xlarge`. Pulumi may report the ASG as failed while the ASG
  itself is still succeeding in another AZ; re-run `pulumi up`.
- **Helm releases with install hooks (cert-manager, kube-prometheus-stack) fail
  if no node has joined yet** — `skipAwait` does not skip hooks. A re-run
  upgrades the failed releases in place once nodes exist.
- **The recipe is H100-validated; the L4 runs the identical software path** on
  cheaper silicon (~$1/hr vs ~$13/hr per node). The provider accepts no low-end
  accelerator value.
- **The chat page is public and unauthenticated** (that's the brief). Bounded by
  `max_tokens`, a short history window and a prompt character budget; the LLM
  API itself is not exposed. Destroy right after the demo.
- **`destroy` leftovers**: Helm keeps some CRDs by policy (harmless, they go with
  the cluster). If a Helm release refuses to uninstall (a controller already
  gone), `pulumi state delete <urn>` and destroy again — the cluster deletion
  takes everything in it.
- **An orphaned 50 GB EBS volume after destroy** (Prometheus PVC): the program
  deletes all PVCs at destroy time (`release-pvcs.sh`, before the stack is
  uninstalled) — verified on the 2026-09-18 teardown, `leakcheck.sh` came back
  clean. If one ever shows up, `aws ec2 delete-volume --volume-id <vol>`.
- **`failed to delete release: <name>` for several Helm releases at once**:
  Pulumi uninstalls the recipe's releases in parallel and a burst of them can
  fail within seconds (seen with cert-manager, nfd, grove, agentgateway-crds on
  2026-09-18), leaving Helm status `uninstalling` and the pods untouched. Nothing
  is stuck — re-run `pulumi destroy`; the second pass uninstalls them cleanly.
- **`DependencyViolation` deleting the node security group**: a stray
  `aws-K8S-<instance-id>` network interface (VPC CNI secondary ENI) left behind
  by a terminated node still carries the SG. `leakcheck.sh` lists it; delete it
  with `aws ec2 delete-network-interface --network-interface-id <eni>` and run
  `pulumi destroy` again. Seen on the first teardown (2026-09-18).
- **Version warnings on destroy** (`using an older version of package … than the
  specified program version`) are expected: pulumi-eks and the AICR provider
  create resources with their own bundled `aws`/`kubernetes` versions, and
  Pulumi deletes each resource with the version that created it.
