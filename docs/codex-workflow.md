# Codex execution workflow

This document defines how Codex tasks are prepared, executed, reviewed, and merged in Vardiel.

## Operating model

Use the following sequence for every implementation task:

```text
accepted issue
  -> dedicated branch
  -> implementation and tests
  -> focused pull request
  -> CI and review
  -> merge
  -> next issue
```

Do not combine roadmap items simply because they touch adjacent packages. The purpose of the sequence is to make security boundaries, regressions, and migration effects reviewable.

## Before coding

Codex must:

1. Read `AGENTS.md`.
2. Read the assigned issue in full.
3. Read `docs/adr/0001-vardiel-project-identity.md`.
4. Read the relevant parts of `docs/architecture.md`.
5. Read `docs/development-roadmap.md` and `docs/v0.1-scope.md`.
6. Inspect the current implementation, tests, and public examples affected by the task.
7. State a concise implementation plan before editing.
8. Identify public contracts and safety boundaries that must remain unchanged.
9. Verify any current SDK or protocol facts against official primary sources.

If the issue is ambiguous, stop and ask for clarification rather than silently broadening the task.

## Branch naming

Use one of:

```text
chore/<short-task>
fix/<short-task>
feat/<short-task>
security/<short-task>
docs/<short-task>
```

Do not work directly on `main`.

## Scope rules

- One issue maps to one branch and one pull request.
- Do not perform drive-by refactors.
- Do not reformat unrelated files.
- Do not upgrade dependencies outside a dependency task.
- Do not rename public identifiers outside the dedicated rename task.
- Do not copy source from reference projects.
- Do not add speculative abstractions that are not required by the issue.
- Do not leave dead compatibility layers without a documented migration need.
- Preserve deterministic output ordering.

## Security review questions

Before implementation and again before opening the pull request, answer:

1. Can model-controlled input reach a shell, code interpreter, arbitrary path, arbitrary method, or arbitrary infrastructure endpoint?
2. Can a remote error body, log line, environment value, credential, path, prompt, or tool result escape into a public error, telemetry, stdout, MCP response, or case report?
3. Are timeouts, byte limits, item limits, attempts, and Agent steps bounded?
4. Does DNS resolution or redirect behavior create an SSRF or rebinding path?
5. Does the tool definition truthfully describe read-only, idempotent, destructive, and open-world behavior?
6. Can execution success be confused with an unhealthy observed state?
7. Can a finding or conclusion be produced without supporting evidence identifiers?
8. Does a new capability duplicate Astralith or Kube-Sentinel responsibilities?
9. Does failure remain safe and observable?

A pull request must not merge while any answer is unknown or inadequately tested.

## Testing expectations

Tests should cover, where applicable:

- normal success
- expected unhealthy observations
- malformed JSON
- unknown fields and trailing JSON values
- empty and oversized values
- invalid semantic values
- boundary counts and timeouts
- cancellation
- deterministic ordering
- missing configuration
- permission denied and unavailable dependencies
- server-side non-2xx or malformed responses
- secret, path, token, header, and prompt-injection strings
- telemetry and public-output non-leakage
- compatibility with CLI, Agent Registry, and MCP exposure

Prefer table-driven tests and local deterministic test servers. Do not require external network services in the default test suite.

## Validation commands

Run the commands relevant to the task and record the exact result:

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./cmd/vardiel
```

During the identity-migration task, the build path changes from `./cmd/opspilot` to `./cmd/vardiel`. Use the path appropriate to the commit under test.

Also run focused package tests while iterating. Full validation is still required before requesting review.

## Pull request structure

Use `.github/pull_request_template.md` and include:

- linked issue
- summary
- exact scope and non-goals
- architecture impact
- security and privacy impact
- public contract and migration impact
- test coverage
- validation output
- unresolved risks or uncertainty

A pull request description must distinguish implemented behavior from future roadmap items.

## Review and merge

- Keep the pull request in draft while tests or scope are incomplete.
- Do not enable auto-merge before a human has reviewed security-sensitive changes.
- Address review comments with new commits; do not silently rewrite shared history.
- Re-run affected tests after review changes.
- Prefer squash merge for focused task branches unless preserving internal commit structure has clear value.
- Delete the task branch after merge.

## First task sequence

The intended initial sequence is:

1. behavior-preserving OpsPilot-to-Vardiel rename
2. safe public error boundary
3. Tool Contract v2
4. evidence and local case bundle
5. Linux host diagnostics
6. deterministic analyzers
7. typed diagnostic plans
8. three built-in playbooks
9. model factory and OpenAI-compatible adapter
10. `v0.1.0` release preparation

Do not skip the public-error and tool-contract foundations to reach feature work faster.
