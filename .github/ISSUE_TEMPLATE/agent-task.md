---
name: Agent implementation task
about: Define one focused Codex-ready engineering task
title: "[Agent Task] "
labels: ""
assignees: ""
---

## Objective

<!-- One concrete outcome. -->

## User value

<!-- Why this matters for Vardiel's daily-use diagnostic workflow. -->

## Context

<!-- Relevant current behavior, packages, ADRs, prior issues, or constraints. -->

## In scope

- 

## Out of scope

- 

## Public contracts affected

<!-- CLI, JSON, MCP, environment variables, tool schemas, case bundles, reports, or none. -->

## Architecture constraints

- Read `AGENTS.md`.
- Read `docs/adr/0001-vardiel-project-identity.md`.
- Preserve provider-neutral boundaries in `internal/agent`.
- Keep the task within Vardiel's responsibilities; do not duplicate Astralith or Kube-Sentinel.

## Security and privacy invariants

- No arbitrary shell or model-generated code execution.
- Read-only by default.
- Model and external data are untrusted.
- Strict JSON and semantic validation.
- Stable sanitized public errors.
- No credentials, raw authorization data, raw remote error bodies, or complete sensitive output in public surfaces.
- Preserve SSRF, redirect, TLS, fixed-endpoint, response-size, result-count, timeout, and redaction boundaries.

## Required implementation

- 

## Required tests

- normal success
- expected unhealthy observation, when applicable
- malformed and unknown input
- boundary values
- timeout and cancellation
- deterministic ordering
- privacy and secret-injection regression cases
- Registry, CLI, Agent, and MCP exposure where applicable

## Acceptance criteria

- [ ] 
- [ ] No unrelated dependency, naming, formatting, or refactor changes.
- [ ] Documentation and examples match the implemented behavior.
- [ ] Security and privacy impact is explained in the pull request.

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

Use the pre-rename build path only during the dedicated identity migration task.

## Deliverables

- implementation
- tests
- documentation or ADR updates when contracts change
- focused pull request linked to this issue
- exact validation results
