# ADR 0001: Vardiel project identity and responsibility boundaries

- Status: Accepted
- Date: 2026-09-02
- Decision owners: Nesoriel maintainers

## Context

OpsPilot already contains a provider-neutral Go Agent Runtime, a shared typed tool registry, an Ark model adapter, an MCP stdio server, privacy-aware observability, and bounded read-only diagnostics for network, Docker, Kubernetes, Prometheus, and Loki.

The next stage is not another short-lived prototype. The project is intended to become a daily-use operations diagnostic agent for personal servers, homelabs, and small infrastructure environments. At the same time, Netiarius demonstrates useful product ideas around explicit plans, playbooks, diagnosis/repair/validation flow, and human approval.

Nesoriel also maintains Astralith and Kube-Sentinel. Without explicit boundaries, these projects could duplicate host management, scheduling, diagnosis, approval, and Kubernetes remediation responsibilities.

## Decision

### Independent continuation

The project will be renamed from OpsPilot to **Vardiel** while preserving its Git history and existing Apache-2.0 license.

Vardiel is the independent continuation of OpsPilot. It is not a GitHub fork of Netiarius and will not track Netiarius as an upstream dependency.

Netiarius may be used as a design reference for:

- explicit diagnostic planning
- playbook-guided investigation
- diagnosis, repair proposal, and validation stages
- human approval and confirmation interaction
- user-readable tool execution traces

Netiarius source code will not be copied. Vardiel will not adopt arbitrary model-generated Python execution or in-process `exec()` as an execution model.

### Product definition

Vardiel is an evidence-driven operations diagnostic agent. It collects bounded structured observations through typed tools, applies deterministic analyzers, and uses language models to explain evidence, form hypotheses, and recommend next steps.

The default behavior is read-only. Vardiel may produce reviewable action proposals, but state-changing operations are outside v0.1.

### Safety invariants

Vardiel will preserve these invariants:

1. No arbitrary shell command interface.
2. No model-generated code execution.
3. Read-only defaults.
4. Model output and external data are untrusted.
5. Strict schemas and semantic validation for tool calls.
6. Bounded steps, timeouts, response sizes, and result counts.
7. Public errors are stable and sanitized; internal causes do not cross public boundaries.
8. Tool outputs are projected and redacted before reaching models or clients.
9. Conclusions and findings are traceable to evidence identifiers.
10. Future mutations require deterministic policy, explicit approval, validation, audit, and rollback planning.

### Responsibility boundaries

#### Vardiel

Vardiel owns:

- Agent Runtime and model-provider adapters
- diagnostic orchestration and typed plans
- diagnostic playbooks
- evidence, findings, reports, and local case bundles
- deterministic analyzers
- bounded diagnostic tool packs
- MCP, CLI, and machine-readable integration contracts
- policy metadata and action proposals

#### Astralith

Astralith owns:

- host inventory and host groups
- Ansible Runner and scheduled remote execution
- web UI, authentication, and centralized persistence
- GitOps desired state, diffs, and apply records
- centralized approvals and audit presentation

Vardiel may integrate with Astralith through stable MCP, API, or JSON contracts. Vardiel will not build a competing inventory, scheduler, or web control plane.

#### Kube-Sentinel

Kube-Sentinel owns:

- Alertmanager event ingestion
- Kubernetes CRDs and controller reconciliation
- event-driven remediation gates
- controller-managed snapshots
- bounded automatic Kubernetes recovery

Vardiel may provide advisory diagnosis or action proposals to Kube-Sentinel, but it will not become a second Kubernetes controller.

## Consequences

### Positive

- Existing Go safety work and tests remain the technical foundation.
- The project gains a distinct long-term identity without inheriting another project's release cycle.
- Plan and playbook features can be added without weakening the typed-tool execution model.
- Responsibility boundaries reduce duplication across Nesoriel projects.
- Vardiel can be used directly from the CLI, by external agents through MCP, and later by Astralith.

### Costs

- The rename affects module paths, binary names, environment variables, deployment resources, documentation, and integration examples.
- Typed plans and evidence contracts require more up-front design than free-text prompts.
- Some useful operational data will remain intentionally unavailable until a safe projection and redaction policy exists.

### Rejected alternatives

#### Fork Netiarius and continue from its Python codebase

Rejected because Vardiel already has a stronger Go execution foundation, typed tools, protocol adapters, bounded APIs, and privacy tests. Adopting the Python CodeAgent execution path would weaken the intended security model.

#### Merge Vardiel into Astralith

Rejected because the reusable diagnostic runtime and MCP surface should remain independent from the web control plane, inventory, scheduler, and GitOps system.

#### Merge Vardiel into Kube-Sentinel

Rejected because Kube-Sentinel is an event-driven Kubernetes controller, while Vardiel is a general diagnostic runtime spanning hosts, containers, clusters, and observability systems.

## Migration policy

The rename is a dedicated behavior-preserving task. It must not include feature work or dependency upgrades.

Until the rename task is merged, references to OpsPilot in code are expected. After the rename, new public identifiers use:

- project and product: `Vardiel`
- binary and MCP server name: `vardiel`
- Go module: `github.com/Nesoriel/vardiel`
- environment prefix: `VARDIEL_`

No long-term compatibility aliases are required because OpsPilot has not published a stable release. Any exception must be justified by an existing integration or automated test.
