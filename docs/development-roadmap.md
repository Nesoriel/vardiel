# Vardiel development roadmap

This roadmap optimizes for one outcome: restore common Linux failures quickly
and safely without interrupting the operator. Every implementation item still
requires one focused issue, branch, pull request, and review.

## Product north star

Vardiel is successful when a known incident follows this path:

```text
detect locally
-> diagnose deterministically
-> execute a pre-authorized fixed action
-> validate recovery
-> record quietly
```

Humans are notified only when authority, risk, confidence, validation, or
recovery budgets require a decision. Known recovery must not depend on a model
or remote control plane.

The product grows along three core incident families: error troubleshooting,
environment failure recovery, and deployment troubleshooting. New work is
prioritized by measured incident frequency and recovery value, not by tool or
integration count.

## Current baseline

The repository currently provides:

- a provider-neutral bounded Agent Runtime and Ark adapter;
- one shared Registry used by Agent, CLI, and MCP;
- machine-readable CLI and MCP stdio surfaces;
- privacy-safe JSONL events and optional OpenTelemetry traces;
- bounded read-only DNS, HTTP, TLS, Docker, Kubernetes, Prometheus, and Loki
  diagnostics;
- strict input validation, SSRF protection, output projection, timeouts, and
  privacy tests.

It does **not** yet provide an always-on daemon, Linux event subscriptions,
standing recovery policy, mutating actions, automatic recovery, or notification
delivery. The baseline archive is `archive/opspilot-baseline-2026-07-16` at
`464349280efb48b64acd5ed1195187045e66cd5c`.

## Immediate sequence

### 0. Accept the direction reset

- merge the documentation and ADR work tracked by issue #33;
- rewrite tracking issue #22 to use this roadmap after the ADR is accepted;
- keep current runtime behavior read-only until action issues are implemented.

### 1. Finish the narrow public-error boundary

Issue #19 remains valuable security hardening because provider and
infrastructure text must not leak through future incident records or
notifications. Finish its exact scoped contract without adding case storage,
Tool Contract v2, new tools, or recovery behavior.

This is the last horizontal prerequisite before the first recovery slice.

### 2. Ship one systemd recovery vertical slice

Create one implementation issue for `systemd-service-unavailable` that delivers
the smallest complete loop:

- an always-on `vardield` process with clean startup, cancellation, reconnect,
  and shutdown;
- systemd unit failure events through a typed adapter, with a bounded fallback
  check when signals are unavailable;
- one deduplicated incident state machine and per-target lock;
- minimal bounded unit, listener, and endpoint observations collected in
  parallel where independent;
- one local versioned knowledge entry matched by the stable unit failure class;
- explicit standing policy for an allowlisted unit and one fixed restart action;
- preconditions, timeout, attempt budget, cooldown, and circuit breaker;
- post-action unit and endpoint validation;
- sanitized local incident/audit output;
- quiet success plus one deduplicated, bounded HTTPS webhook escalation path;
- notification retry isolated from recovery so delivery failure cannot repeat
  an action;
- measurement of the ADR 0002 latency budgets;
- no model or remote service call on the known path.

Do not begin with a generic workflow framework. Extract shared contracts only
after this slice proves what must be shared.

### 3. Establish the minimal operational knowledge pack

Store repository-reviewed operations manuals, known error remedies, and
deployment runbooks as local versioned Markdown entries with bounded metadata,
source provenance, and references to registered observation or playbook IDs.
Start with deterministic scope, tag, error-code, and signature matching plus a
bounded text-search fallback.

Knowledge is advisory and untrusted. It cannot define executable commands or
grant action authority. Measure lookup recall on sanitized incident fixtures
before considering semantic retrieval, a vector database, or a remote knowledge
service.

### 4. Harden the local recovery runtime

After the vertical slice works on a real Linux fixture:

- crash-safe incident resumption and duplicate-action protection;
- least-privilege daemon/action separation and multi-user authorization;
- policy expiry, maintenance windows, concurrency and blast-radius limits;
- action audit retention and privacy rules;
- additional notifier integrations and durable bounded delivery retry;
- install, upgrade, rollback, and service-hardening documentation;
- performance and soak tests for event bursts and flapping services.

### 5. Expand the three core incident families

Add one vertical slice at a time, chosen from real incident frequency:

1. error troubleshooting through stable error signatures, bounded log metadata,
   service/dependency correlation, and cited known remedies;
2. environment failures such as CPU/memory/I/O or filesystem pressure,
   permissions, dependencies, DNS, TLS, and local endpoint differences;
3. deployment troubleshooting through a sanitized release marker and
   before/after service and endpoint observations, without automatic deployment
   rollback until a separate fixed-action review;
4. Docker container restart only after a separate socket-risk and action-policy
   review.

Each slice must include detection, observations, policy, action when safe,
validation, cooldown/circuit breaking, notification behavior, and a real Linux
fixture. Kubernetes mutations remain with Kube-Sentinel.

### 6. Add bounded model reasoning for unknown incidents

Only after deterministic recovery is proven:

- let a provider-neutral model rank hypotheses from bounded observations;
- allow selection only among registered observation steps and playbooks;
- require deterministic validation and policy for every proposed action;
- enforce a strict model latency and cost budget;
- fall back to escalation when the model is unavailable or uncertain;
- evaluate diagnosis quality on sanitized incident fixtures.

The model never grants authority and never emits executable shell or code.

### 7. Add fleet and multi-user integration

Integrate with Astralith for identity, host inventory, policy distribution,
fleet correlation, notification delivery, UI, and centralized audit. Preserve
the per-host daemon as the authority for the local low-latency loop so a control
plane outage cannot block an approved recovery.

## v0.1.0 release

The exact v0.1 contract is in [v0.1-scope.md](v0.1-scope.md). Its release proof
is one reliable single-host systemd recovery loop, not the number of adapters,
schemas, or documents shipped.

## Re-sequenced existing issues

| Issue | Disposition after issue #33 is accepted |
| --- | --- |
| #19 safe public errors | Keep, finish narrowly before the first mutating slice. |
| #20 Tool Contract v2 | Supersede. Define only the observation/action metadata proven necessary by the systemd slice. |
| #21 evidence and case bundles | Defer. The first slice needs a minimal sanitized incident/action journal, not a general case platform. |
| #22 v0.1 tracking | Rewrite around the systemd closed loop and this roadmap. |

Do not close or rewrite these issues merely because this branch exists. Update
them after the ADR and roadmap are merged.

## Explicitly deprioritized

The following work is not on the critical path to the first recovery loop:

- a universal Tool Contract v2 migration across every current adapter;
- a general evidence/case-bundle storage system;
- more model providers, semantic retrieval, or a vector database before simple
  knowledge lookup has a measured recall gap;
- a large generated tool catalog;
- remote SSH inventory or Ansible scheduling inside Vardiel;
- a web UI or centralized authentication inside Vardiel;
- more Kubernetes, Prometheus, or Loki surface area;
- arbitrary shell, scripts, or model-generated code;
- broad Docker or Kubernetes mutation.

Add one only when a proven recovery slice needs it or a maintainer accepts a
separate product reason.

## Release-line outlook

### v0.1 — one fast local recovery loop

One allowlisted systemd service, event-driven detection, standing policy, fixed
restart action, validation, audit, circuit breaking, and quiet success.

### v0.2 — useful host coverage

Error-signature diagnosis, environment-pressure and local-endpoint playbooks,
sanitized deployment correlation, a measured local knowledge pack, hardened
privilege separation, and real-host soak and recovery tests.

### v0.3 — unknown incidents and fleets

Bounded model-assisted diagnosis, Astralith fleet/multi-user integration,
cross-host correlation, and centrally distributed policy without making the
control plane a local recovery dependency.

### v1.0 — production confidence

Multiple proven recovery playbooks with stable policy and incident schemas,
upgrade/rollback support, measurable reliability and latency, security review,
and documented operational limits.

## Metrics that govern priority

Prefer work that improves one of these measures:

- detection latency;
- time to first safe action;
- time to verified recovery;
- autonomous recovery success rate;
- false-action and rollback rate;
- repeat-incident/circuit-breaker rate;
- percentage of incidents requiring a notification;
- operator time required per unresolved incident.

Do not optimize tool count, prompt size, schema count, or integration count as a
proxy for recovery value.
