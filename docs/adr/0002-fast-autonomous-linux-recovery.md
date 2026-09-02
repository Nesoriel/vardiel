# ADR 0002: Fast autonomous Linux incident recovery

- Status: Accepted
- Date: 2026-09-02
- Decision owners: Nesoriel maintainers
- Supersedes: the read-only product definition and mutation-approval model in
  [ADR 0001](0001-vardiel-project-identity.md)

## Context

Vardiel currently has a safe Go runtime and bounded, read-only diagnostics, but
its roadmap puts several horizontal contracts and integration surfaces before
the product's primary job: restore a failed Linux service quickly and with as
little operator interruption as possible.

The target environment ranges from one personal server to fleets shared by
multiple users. A remote model call, central control plane, or interactive
confirmation cannot be required on the critical path for a known, locally
recoverable failure. At the same time, broad autonomous shell access would turn
an operational incident into a privilege-escalation and data-loss risk.

Netiarius demonstrates useful patterns: one tool runtime, explicit risk
metadata and preflight checks, playbook-guided diagnosis, and a
diagnosis-to-repair-to-validation flow. Its CLI-first confirmation flow,
model-generated playbooks during incidents, and model-generated Python
execution do not meet Vardiel's latency or security goals.

Linux already supplies efficient event and control primitives. systemd exposes
unit state and fixed management operations, the Journal provides change
notifications, and Pressure Stall Information (PSI) supports pollable CPU,
memory, and I/O pressure thresholds. Vardiel should compose these primitives
rather than continuously ask a model what to inspect.

## Decision

### Product mission

Vardiel is a fast, low-interruption Linux incident-response agent. It detects,
diagnoses, and safely recovers common host and service failures, then validates
the result. It asks for human attention only when the configured authority,
risk, confidence, or recovery budget does not permit autonomous completion.

Current released code remains read-only. This ADR authorizes later focused
issues to add bounded mutations; it does not itself change runtime behavior.

### Two-speed incident loop

Known incidents use a deterministic local fast path:

```text
event
-> correlate and deduplicate
-> match a versioned recovery playbook
-> collect the minimum required observations concurrently
-> evaluate standing policy and preconditions
-> execute one fixed typed action
-> validate health
-> close quietly or roll back / stop / escalate
```

The language model is not required on this path.

Unknown or ambiguous incidents use a bounded reasoning path. The model may rank
hypotheses, select from registered observations or playbooks, explain evidence,
and recommend a typed action. It cannot create a tool, action, shell command,
script, API path, or policy grant. The same policy engine controls execution no
matter who selected the action.

### Per-host execution and fleet coordination

An always-on Vardiel daemon will own the incident loop on each Linux host. Local
detection and recovery continue when the network, model provider, or fleet
control plane is unavailable.

Vardiel owns incident state, local observations, playbook execution, policy
evaluation, local actions, validation, circuit breaking, and the sanitized
incident record. Astralith may supply identity, host inventory, fleet dispatch,
central policy distribution, UI, notification delivery, and centralized audit.
Kube-Sentinel continues to own Kubernetes controller reconciliation and
automatic Kubernetes recovery.

### Authorization and action safety

Routine recovery does not require an interactive prompt when an administrator
has explicitly granted standing authorization for the exact host or host group,
target, action, parameters, and limits. Standing authorization is configuration,
not a model inference.

Every action implementation must have:

- a fixed name and typed, strictly validated arguments;
- an allowlisted target and least-privilege execution identity;
- deterministic policy evaluation before execution;
- bounded blast radius, timeout, attempts, and concurrency;
- idempotency or duplicate suppression plus a per-target lock or lease;
- preconditions captured immediately before execution;
- post-action health validation;
- a rollback step when rollback is safe and meaningful;
- an explicit non-rollback classification otherwise;
- cooldown and circuit-breaker behavior for repeated failures;
- a sanitized audit record that does not expose secrets or raw logs.

Irreversible, broad, security-sensitive, or unbounded actions remain
human-gated. A model cannot widen standing authorization.

### Quiet operation

A successfully auto-recovered incident is recorded but does not page the user.
Vardiel escalates when at least one of these conditions holds:

- no approved playbook or standing authorization matches;
- the action is high risk, irreversible, or exceeds configured blast radius;
- evidence is contradictory or confidence is below policy;
- validation fails, rollback fails, or the recovery deadline expires;
- retry, action, or concurrency budgets are exhausted;
- the target repeatedly flaps and opens its circuit breaker;
- probable data loss, credential compromise, or cross-host impact is detected;
- the local daemon cannot safely observe or act.

An escalation must state impact, observations, attempted actions, validation,
current safety state, and the smallest decision required from the operator.

### Performance budgets

The first vertical slice targets, on a supported local systemd host:

- P95 event intake within 1 second of a unit entering a failed state;
- P95 approved recovery action start within 5 seconds of event intake;
- verified recovery or explicit escalation within 30 seconds for the fixture;
- zero model or remote-control-plane calls on the known recovery fast path.

These are product budgets, not documentation promises. Implementations must
publish the measurement environment and test the configured deadlines without
introducing flaky wall-clock assertions.

### First vertical slice

The first implementation milestone is `systemd-service-unavailable` for one
local host and one explicitly allowlisted service:

1. observe the unit failure through systemd;
2. create or update one deduplicated incident;
3. collect bounded unit, listener, and health-check facts;
4. evaluate standing authorization for a fixed systemd restart action;
5. execute at most the configured number of attempts;
6. validate unit and endpoint health;
7. close quietly on success or stop and send one bounded escalation on failure
   or flapping.

This vertical slice comes before a generic fleet platform, large tool catalog,
RAG system, web UI, or universal case-store abstraction.

## Preserved invariants

- No arbitrary shell interface or model-generated code execution.
- Model output, logs, provider data, and remote responses remain untrusted.
- Current strict schemas, SSRF protection, TLS verification, fixed endpoint
  allowlists, redacted projections, byte/count limits, timeouts, cancellation,
  and protocol-stream boundaries remain in force.
- Raw credentials, environment values, authorization headers, private paths,
  provider bodies, and unbounded logs do not enter model or public output.
- Current Docker, Kubernetes, Prometheus, and Loki operations remain read-only
  until their own focused action ADR and implementation tasks are accepted.
- Current CLI and MCP behavior does not become mutating merely because this ADR
  is accepted.

## Consequences

### Positive

- The roadmap delivers a useful recovery loop before broad platform features.
- Common incidents avoid model and network latency.
- Explicit standing policy provides low interruption without granting the model
  open-ended privilege.
- One local loop works on a single server and composes with a future fleet
  control plane.

### Costs

- The daemon and privileged action boundary require stronger lifecycle,
  authorization, audit, concurrency, and crash-recovery tests.
- Not every action has a meaningful rollback; such actions need explicit policy
  and must fail closed.
- Unknown incidents may still require model latency or human intervention.
- Existing roadmap issues must be re-sequenced or replaced after this ADR is
  reviewed.

## Rejected alternatives

### Keep Vardiel read-only until every horizontal contract is complete

Rejected because it delays an end-to-end user outcome and provides no evidence
that the proposed abstractions improve recovery speed or safety.

### Put a language model in every incident decision

Rejected because provider latency and availability are unsuitable for a known
local recovery path, and nondeterministic reasoning cannot grant authority.

### Copy Netiarius or adopt model-generated Python execution

Rejected because Vardiel already has a safer Go foundation, and arbitrary
generated code conflicts with the fixed-action boundary. Netiarius remains a
design reference only.

### Build remote SSH inventory and fleet scheduling into Vardiel first

Rejected because remote coordination adds latency and failure modes to local
recovery and duplicates Astralith responsibilities.

### Notify before every mutation

Rejected for pre-authorized, bounded routine recovery whose rollback or
non-rollback behavior is explicitly accepted because it defeats the
low-interruption goal. High-risk or out-of-policy operations remain human-gated.

## Research references

- [Netiarius system architecture](https://github.com/Joeen-AI-Labs/Netiarius/blob/main/doc/systeminfo.md)
- [Netiarius tool integration guide](https://github.com/Joeen-AI-Labs/Netiarius/blob/main/doc/ADD_Tools.md)
- [systemd service restart and watchdog semantics](https://github.com/systemd/systemd/blob/main/man/systemd.service.xml)
- [systemd unit failure dependencies](https://github.com/systemd/systemd/blob/main/man/systemd.unit.xml)
- [systemd D-Bus API](https://github.com/systemd/systemd/blob/main/man/org.freedesktop.systemd1.xml)
- [Linux Pressure Stall Information](https://docs.kernel.org/accounting/psi.html)
- [systemd Journal change-notification API](https://github.com/systemd/systemd/blob/main/man/sd_journal_get_fd.xml)
- [polkit authorization architecture](https://polkit.pages.freedesktop.org/polkit/polkit.8.html)
