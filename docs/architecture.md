# Architecture

## Status and source of truth

Vardiel's accepted target architecture is a fast, event-driven Linux incident
response and recovery agent. The decision and safety model are recorded in
[ADR 0002](adr/0002-fast-autonomous-linux-recovery.md).

The current implementation is still a CLI/MCP-oriented, read-only diagnostic
runtime. Sections marked **target** describe planned work, not shipped behavior.

## Core principle

Known local failures must not wait for a language model or remote control plane.
Vardiel uses a deterministic fast path for approved recovery and a bounded
reasoning path for unknown incidents.

```text
Linux event sources                         Optional fleet plane
systemd D-Bus / Journal / PSI / timers      Astralith
            |                               identity / inventory / policy
            v                                      |
    per-host vardield  <---------------------------+
    - correlate and deduplicate incidents
    - run typed playbooks
    - evaluate standing policy
    - execute fixed actions
    - validate, cool down, circuit-break
            |
       +----+-----------------+
       |                      |
       v                      v
 deterministic fast path   bounded reasoning path
 no model call             model ranks hypotheses or selects
 known approved recovery   registered observations/playbooks
       |                      |
       +----------+-----------+
                  v
          quiet close or escalation
```

Remote connectivity may improve a diagnosis, but loss of the network, model
provider, or fleet plane must not prevent a local known recovery.

## Performance budgets

For the first supported `systemd-service-unavailable` fixture:

| Stage | Target |
| --- | --- |
| systemd failure to local event intake | P95 <= 1 second |
| event intake to approved action start | P95 <= 5 seconds |
| event intake to verified recovery or escalation | <= 30 seconds |
| known fast-path model/control-plane calls | zero |

Implementations must make stage timestamps observable, preserve a total
deadline, and test configured budgets with deterministic fakes. A real Linux
systemd fixture supplies release evidence; CI timing alone is not a production
latency claim.

## Target per-host components

### Event sources

Prefer kernel and service-manager notifications to tight polling:

- systemd D-Bus unit and job state changes;
- bounded structured Journal fields where a unit-state signal is insufficient;
- pollable Linux PSI thresholds for CPU, memory, and I/O pressure;
- bounded timers for filesystems, endpoints, and compatibility fallbacks.

Events are untrusted inputs. Each source must bound payloads, normalize a stable
event type, record monotonic and wall-clock time, and survive reconnects without
replaying an unbounded history.

### Incident coordinator

The coordinator owns a small state machine:

```text
observed -> triaging -> authorized -> acting -> validating -> recovered
                    +-> blocked ---------------------------> escalated
                                   +-> rollback -> escalated or recovered
```

It deduplicates equivalent events, serializes actions per target, preserves one
total deadline, and records why each transition occurred. Repeated failures
open a circuit breaker instead of starting an action loop.

### Playbook engine

A playbook is a versioned typed recovery workflow, not free-form model text. It
declares:

- applicability and required observations;
- concurrent observation groups and dependencies;
- action candidates and validation checks;
- per-step and total deadlines;
- retry, cooldown, and circuit-breaker rules;
- rollback or explicit non-rollback behavior;
- escalation conditions.

Known playbooks execute without a model. A model may choose among registered
playbooks for an unknown incident, but validation and policy remain
deterministic.

### Policy engine

Standing authorization is explicit administrator configuration. Policy binds
an identity or service principal to an exact host or group, target, action,
parameter limits, time window, attempt budget, concurrency budget, and maximum
blast radius.

Policy evaluation fails closed. Missing, ambiguous, expired, conflicting, or
unavailable policy never becomes an implicit grant. Multi-user requests are
authorized before they cross the privileged action boundary.

### Action executor

Actions are fixed typed implementations. There is no arbitrary shell, script,
model-generated code, arbitrary filesystem path, or arbitrary infrastructure
API method.

Before execution, the executor rechecks current preconditions and obtains a
per-target lock or lease. During execution it enforces cancellation, timeout,
attempt, concurrency, and output bounds. Afterwards it runs independent health
validation. An action with no safe rollback must declare that fact and be
allowed explicitly by policy.

The first planned action is a bounded systemd restart for an allowlisted unit.
It does not generalize to arbitrary `systemctl` arguments.

### Validation, audit, and notification

Execution success is not recovery. The validator checks the intended health
contract, such as stable unit state plus a bounded TCP or HTTP probe.

Every incident records sanitized timestamps, state transitions, stable reason
codes, selected playbook, policy decision, action outcome, validation outcome,
and escalation reason. Raw provider bodies, logs, credentials, environment
values, private paths, and complete command lines are excluded.

Successful routine recovery is quiet. Notification occurs only for an
out-of-policy, high-risk, ambiguous, repeated, failed, expired, or potentially
data-destructive incident. The first standalone transport is one configured,
bounded HTTPS webhook; notification delivery may later be delegated to
Astralith. Delivery failure is retried independently and never repeats a
recovery action. The local incident still reaches a durable terminal state when
delivery is unavailable.

## Model boundary

The model is optional for known recoveries. It may:

- rank hypotheses from bounded observations;
- choose among registered observation steps or playbooks;
- explain deterministic findings and uncertainty;
- recommend a registered typed action for policy evaluation.

It may not create an observation source, tool, action, shell command, script,
API path, target, credential, policy grant, or final success claim. Model output
is untrusted and cannot directly drive the privileged executor.

`internal/agent` remains provider-neutral. Provider adapters stay under
`internal/models/<provider>`.

## Single-host, fleet, and multi-user boundaries

The per-host daemon owns the low-latency loop. A standalone installation can
manage one host with locally configured policy.

Astralith may own:

- user identity, authentication, tenancy, and host inventory;
- host grouping and central policy distribution;
- fleet-wide incident correlation and dispatch;
- notification delivery, UI, and centralized audit presentation;
- remote execution transport when local agents are not available.

Vardiel owns the diagnostic and recovery semantics executed on a host. It does
not embed a second inventory, general scheduler, Ansible runner, or web control
plane. Kube-Sentinel continues to own Kubernetes CRDs, controller
reconciliation, and automatic Kubernetes recovery.

## Current implementation boundaries

The following boundaries describe shipped code and remain in force while the
target daemon is built.

### Packages and public surfaces

- `internal/agent`: provider-neutral messages, registry, bounded agent loop,
  and lifecycle events.
- `internal/models`: provider adapters; provider SDKs do not enter
  `internal/agent`.
- `internal/tools`: strictly decoded, bounded, read-only tools.
- `internal/dockerapi`, `internal/kubeapi`, `internal/promapi`, and
  `internal/lokiapi`: fixed read-only adapters with redacted projections.
- `internal/mcpserver`: adapter from the shared Registry to the official MCP Go
  SDK; it is not a second execution engine.
- `internal/observability`: privacy-safe JSONL and OpenTelemetry translation.
- `cmd/vardiel`: machine-readable CLI and MCP stdio process boundary.

CLI, Agent, and MCP consume the same Registry. MCP stdio stdout remains
protocol-only. Current tools are read-only and idempotent; a future action API
must be a separate truthful contract and must not silently make existing MCP
tools mutating.

### Network and infrastructure adapters

- Generic HTTP and TLS tools preserve SSRF protection, target re-resolution,
  redirect limits, TLS verification, timeouts, and response limits.
- Docker accepts only a trusted local Unix socket and fixed bounded GET
  operations. Raw inspect data, environment values, commands, labels, health
  output, log paths, mount source paths, warnings, and runtime error text are
  not exposed.
- Kubernetes uses validated least-privilege credentials and fixed queries.
  Unsafe kubeconfig features, arbitrary resources, selectors, paths, methods,
  logs, Secrets, environment values, and free-text status/Event messages remain
  outside the model-facing interface.
- Prometheus and Loki use configured trusted endpoints, disable ambient proxies
  and redirects, enforce TLS and response bounds, and expose only fixed
  read-only APIs with generated allowlisted queries/selectors. Raw PromQL,
  LogQL, log lines, API bodies, credentials, and arbitrary paths remain absent.

These privileged sockets and credentials stay privileged even when Vardiel's
current calls are read-only. No autonomous mutation of Docker, Kubernetes,
Prometheus, or Loki is authorized by ADR 0002.

### Public errors and telemetry

Public errors must use stable sanitized content; internal causes do not cross
CLI, Agent, MCP, telemetry, incident records, or notifications. Unknown errors
use a generic safe fallback.

Default telemetry contains lifecycle metadata, stable codes, durations, and
bounded identifiers only. It excludes prompts, model text, arguments, results,
credentials, raw errors, and logs. Export failures do not fail the local
incident loop.

## Research basis

- [ADR 0002](adr/0002-fast-autonomous-linux-recovery.md) records the decision
  and primary sources.
- Netiarius is a design reference for tool/runtime and
  diagnosis-repair-validation structure only; no source is copied.
- Event-driven systemd, Journal, and PSI integrations are preferred over a new
  generic monitoring framework.
