# Vardiel development roadmap

This roadmap defines capability boundaries and sequencing. It is not permission to implement several milestones in one pull request. Every item requires a focused GitHub issue and review.

## Current baseline: Vardiel

The preserved baseline already provides:

- provider-neutral bounded Agent Runtime
- Ark Responses API adapter
- shared typed Registry
- MCP stdio server
- machine-readable CLI
- privacy-safe JSONL lifecycle events
- optional OpenTelemetry tracing
- bounded read-only DNS, HTTP, and TLS tools
- bounded read-only Docker, Kubernetes, Prometheus, and Loki diagnostics
- strict input validation, SSRF protection, output projection, and privacy tests

The baseline archive ref is `archive/opspilot-baseline-2026-07-16` at commit `464349280efb48b64acd5ed1195187045e66cd5c`.

## v0.1.0 — Evidence-driven Linux diagnostics

Goal: deliver a safe daily-use local diagnostic agent and a credible course-project release without opening state-changing infrastructure operations.

### Foundation

- complete behavior-preserving Vardiel identity migration
- establish a unified safe public error boundary
- refresh dependencies in isolated dependency-family pull requests
- introduce Tool Contract v2 with input/output schemas and truthful annotations
- keep Agent, CLI, and MCP on one shared Registry

### Evidence and cases

- define `DiagnosticCase`, `DiagnosticPlan`, `PlanStep`, `Observation`, `Finding`, `DiagnosisReport`, and `ActionProposal`
- assign stable evidence identifiers such as `E1`, `E2`, and `E3`
- require findings and conclusions to reference evidence IDs
- store transparent local case bundles under the user's data directory
- keep raw sensitive data out of public reports and default stdout
- make full traces opt-in rather than default output

### Linux host diagnostics

Add bounded read-only tools for:

- host inventory and operating-system facts
- uptime and load snapshot
- memory and swap snapshot
- pressure stall information when available
- filesystem and inode usage
- process resource hotspots without environment values or full command lines
- failed systemd units
- one systemd unit status
- listening sockets with bounded projections

Implementation rules:

- prefer `/proc`, `statfs`, netlink, and fixed argument builders
- never expose an arbitrary shell string
- never return process environments
- strictly validate unit names and all user-controlled fields
- bound item counts and free-text lengths

### Deterministic analyzers

Add analyzers for:

- host load pressure
- memory and swap pressure
- filesystem and inode pressure
- process hotspots
- systemd service health
- Docker container health
- web endpoint reachability
- Kubernetes Pod health
- Prometheus target health

Analyzers produce deterministic `Finding` objects before the language model explains them. Thresholds belong in validated configuration rather than prompts.

### Typed diagnostic plans

- define a strongly typed plan schema
- validate tool existence, annotations, arguments, dependencies, cycles, attempts, and total budgets
- execute only read-only tools in v0.1
- allow deterministic execution without a configured language model
- preserve partial results when a step fails
- prevent duplicate equivalent calls in one plan

### Built-in playbooks

Ship at least these three versioned playbooks:

1. `host-resource-pressure`
2. `systemd-service-unavailable`
3. `web-endpoint-unreachable`

Each playbook must include deterministic steps, applicability metadata, budgets, expected evidence, and end-to-end fixtures.

### Model adapters and user entry point

- retain Ark support behind the provider-neutral interface
- add one OpenAI-compatible adapter
- support a local OpenAI-compatible endpoint without changing the core Runtime
- introduce `vardiel diagnose` as the primary user command
- keep `tool`, `case`, `playbook`, `doctor`, and `mcp` machine-readable

### v0.1 explicit non-goals

- arbitrary shell execution
- model-generated code execution
- automatic service or container restart
- Kubernetes mutations
- remote SSH inventory or Ansible scheduling
- web UI
- multi-agent orchestration
- RAG or vector database
- raw Journal, container, or application log ingestion
- autonomous scheduled remediation

## v0.2.0 — Integration and external tool surfaces

Goal: make Vardiel reusable from other Nesoriel systems and trusted agent runtimes.

Planned scope:

- MCP client support for approved external tool servers
- capability discovery and per-server policy
- signed or pinned external tool-server configuration
- stable case/report interchange format
- Astralith integration through MCP, API, or JSON contracts
- import of Astralith Evidence Packs without duplicating inventory or scheduling
- optional remote diagnostic transport supplied by Astralith rather than embedded host management
- richer Agent Skill packaging
- case-bundle migration and schema-version support

Non-goals remain a standalone web control plane and a competing host inventory.

## v0.3.0 — Policy, approvals, and action proposals

Goal: prepare safe state changes without enabling broad autonomous mutation.

Planned scope:

- typed `ActionProposal` contracts
- deterministic policy evaluation
- explicit approval state and expiry
- target, scope, blast-radius, maintenance-window, and capability checks
- dry-run or diff support where the underlying system provides it
- pre-change snapshots and rollback plans
- post-action validation plans
- audit events that do not leak sensitive payloads
- integration hooks for Astralith approval UI and Kube-Sentinel advisory flows

Execution support may remain disabled for some or all proposals in this milestone.

## v1.0.0 — Selected bounded actions

Goal: support a deliberately small set of production-quality mutations.

Candidate actions must satisfy all of the following:

- fixed typed interface; no arbitrary shell or API path
- deterministic policy evaluation
- explicit scoped approval
- idempotency or duplicate protection
- preconditions and blast-radius limits
- bounded timeout and retries
- post-change validation
- rollback or safe failure semantics
- complete audit trail
- dedicated integration tests in realistic environments

The exact action list will be selected from validated daily-use scenarios. Broad self-healing is not a v1.0 requirement.

## Quality gates for every milestone

- focused issue and pull request boundaries
- clean `go mod tidy`
- clean formatting and `go vet`
- race-enabled test suite
- buildable CLI
- malformed-input and cancellation tests
- privacy and secret-injection regression tests
- deterministic fixture output
- updated architecture and migration documentation
- no unreviewed weakening of existing safety boundaries

## Release discipline

- use semantic versioning beginning with `v0.1.0`
- publish release notes with compatibility and security impact
- provide checksums for released binaries
- keep the previous stable binary and case-schema migration notes available
- do not describe roadmap items as implemented capabilities
