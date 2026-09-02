---
name: Maintainer / Codex implementation task
about: Define one accepted, focused, coding-agent-ready engineering task
title: "[Agent Task] "
labels: ""
assignees: ""
---

## Objective

<!-- One concrete, observable outcome. -->

## User value

<!-- Why this matters for Vardiel's daily-use diagnostic workflow. -->

## Current behavior and evidence

<!-- Relevant source locations, failing examples, contracts, issues, or research. -->

## Prerequisites

<!-- Issues or pull requests that must be merged first. -->

## In scope

-

## Out of scope

-

## Public contracts affected

<!-- CLI, JSON, MCP, environment variables, tool schemas, case bundles, reports, or none. -->

## Architecture constraints

- Read `AGENTS.md`, `CONTRIBUTING.md`, and the relevant ADRs.
- Preserve provider-neutral boundaries in `internal/agent`.
- Keep the task within Vardiel's responsibilities; do not duplicate Astralith or Kube-Sentinel.
- Work on the maintainer-supplied branch and never push directly to `main`.

## Security and privacy invariants

- No product capability for arbitrary shell or model-generated code execution.
- Read-only by default.
- Model and external data are untrusted.
- Strict JSON and semantic validation.
- No new path that exposes credentials, raw authorization data, raw remote error bodies, or complete sensitive output.
- Preserve SSRF, redirect, TLS, fixed-endpoint, response-size, result-count, timeout, cancellation, and redaction boundaries.
- Record known pre-existing limitations rather than claiming they are already fixed.

## Required implementation

-

## Required tests

- normal success;
- expected unhealthy observation, when applicable;
- malformed, unknown, and trailing input;
- semantic and size boundaries;
- timeout and cancellation;
- deterministic ordering;
- privacy and secret-injection regression cases;
- Registry, CLI, Agent, MCP, telemetry, or storage exposure where applicable.

## Acceptance criteria

- [ ]
- [ ] No unrelated dependency, naming, formatting, or refactor changes.
- [ ] Documentation and examples match implemented behavior.
- [ ] Security, privacy, compatibility, and remaining uncertainty are explained in the pull request.

## Validation

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./cmd/vardiel
```

## Deliverables

- implementation;
- tests;
- documentation or ADR updates when contracts change;
- one focused pull request linked to this issue;
- exact validation results and skipped checks;
- concise AI-assistance and human-verification statement.
