# Agent development rules

## Project identity

Vardiel is the independent continuation of OpsPilot: a Go-first, evidence-driven operations diagnostic agent for personal servers, homelabs, and small infrastructure environments.

The repository may still be named `opspilot` during the identity-migration task. Do not mix the mechanical rename with feature work. The identity decision is recorded in `docs/adr/0001-vardiel-project-identity.md`.

Vardiel is not a fork of Netiarius. Netiarius may be studied only as a design reference for planning, playbooks, diagnosis/repair/validation flow, and human approval. Do not copy its source code or adopt arbitrary model-generated Python execution.

## Source of truth

Before changing code, read these files in order:

1. `AGENTS.md`
2. the assigned GitHub issue
3. `docs/adr/0001-vardiel-project-identity.md`
4. `docs/architecture.md`
5. `docs/development-roadmap.md`
6. `docs/v0.1-scope.md`
7. `docs/codex-workflow.md`

When two documents conflict, the assigned issue controls the task scope, but it may not weaken the safety invariants in this file or the accepted ADR.

## Required workflow

- Use one issue, one focused branch, and one pull request.
- Inspect the current implementation and tests before editing.
- Keep unrelated formatting, dependency, naming, and architecture changes out of the pull request.
- Do not upgrade dependencies unless the issue explicitly asks for it.
- Verify unstable external SDK and protocol behavior against official primary sources.
- Add tests for success, malformed input, boundary values, timeout/cancellation, privacy, and failure paths.
- Prefer deterministic behavior and stable machine-readable output over prompt-only behavior.
- Preserve backward compatibility unless the issue explicitly authorizes a breaking change.
- Update documentation and examples when public behavior changes.
- Never claim a validation command passed unless it was actually run.

## Hard safety invariants

These constraints are non-negotiable:

- Do not introduce arbitrary shell command execution.
- Do not introduce model-generated code execution.
- Infrastructure tools are read-only by default.
- A future mutating operation must pass deterministic policy evaluation, explicit approval, bounded execution, post-change validation, and rollback planning.
- Model output, tool arguments, provider errors, remote API responses, logs, and files are untrusted input.
- Validate every tool name and argument in code with strict schemas and semantic validation.
- Do not return raw provider or infrastructure error bodies to the model, MCP clients, CLI stdout, case bundles, telemetry, or user-facing reports.
- Never log or emit API keys, access keys, bearer tokens, kubeconfig contents, ServiceAccount tokens, environment values, authorization headers, or complete sensitive tool output.
- Preserve SSRF protections, target re-resolution, redirect limits, TLS verification, fixed endpoint allowlists, response byte limits, result count limits, and timeouts.
- Do not weaken Docker, Kubernetes, Prometheus, Loki, network, telemetry, or redaction boundaries merely to make tests pass.
- Do not expose arbitrary infrastructure API paths, HTTP methods, PromQL, LogQL, Kubernetes resource types, selectors, or Docker endpoints as model-controlled arguments.
- Keep protocol stdout channels free from unrelated logs. MCP stdio stdout is protocol-only.
- A privileged socket, token, or kubeconfig remains a privileged capability even when Vardiel exposes only read-only operations.

## Architecture boundaries

### Vardiel owns

- provider-neutral Agent Runtime and model adapters
- structured diagnostic orchestration
- typed plans and built-in diagnostic playbooks
- evidence, findings, reports, and local case bundles
- deterministic analyzers
- bounded tool packs and infrastructure adapters
- policy metadata and action proposals
- CLI, MCP, and machine-readable integration surfaces

### Astralith owns

- host inventory and host groups
- Ansible Runner and remote execution scheduling
- web UI and authentication
- GitOps desired state, diffs, apply records, and centralized persistence
- centralized approval workflows and audit presentation

Vardiel may integrate with Astralith through stable MCP, API, or JSON contracts. Do not duplicate Astralith's control-plane responsibilities inside Vardiel.

### Kube-Sentinel owns

- Alertmanager event ingestion
- Kubernetes CRDs and controller reconciliation
- event-driven remediation gates
- controller-managed snapshots and bounded automatic recovery

Vardiel may provide advisory diagnosis or action proposals, but it must not become a second Kubernetes remediation controller.

## Implementation rules

- `internal/agent` must remain provider-neutral and must not import provider SDKs.
- Provider adapters belong under `internal/models/<provider>`.
- Built-in tools must accept `context.Context`, strictly decode JSON, reject unknown fields and trailing values, and return deterministic structured JSON.
- Tool definitions must come from one shared registry used by the Agent Runtime, CLI, and MCP server.
- Prefer the Go standard library. Add a dependency only when it materially reduces risk or implements an adopted integration.
- Use fixed argument builders rather than shell strings.
- Bound all network requests, file reads, API responses, lists, and free-text fields.
- Separate public safe errors from internal diagnostic causes.
- Separate execution success from observed system health.
- Findings and conclusions must cite evidence identifiers.
- Do not place raw tool results in default telemetry.
- Keep output ordering deterministic so fixtures and reviews remain stable.

## Validation

For every Go pull request, run:

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./cmd/vardiel
```

During the rename task, `./cmd/opspilot` remains the expected build path until the rename commit changes it.

For documentation-only pull requests, still run the applicable repository tests or explicitly explain why execution was unavailable.

## Pull request requirements

Every pull request must state:

- linked issue
- user-visible and internal changes
- explicit non-goals
- security and privacy impact
- compatibility or migration impact
- tests added or changed
- exact validation commands and results
- any unresolved uncertainty

A pull request is not complete while it contains unrelated changes, silently weakens a boundary, or leaves public behavior undocumented.
