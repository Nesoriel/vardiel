# Codex execution workflow

This document defines how maintainers prepare, run, review, and merge Codex tasks in Vardiel. Shared contribution rules live in `CONTRIBUTING.md`; repository-wide coding-agent constraints live in `AGENTS.md`.

## Operating model

Use this sequence for every implementation task:

```text
accepted issue
-> maintainer-created task branch
-> plan
-> implementation and tests
-> focused pull request
-> CI and human review
-> merge
-> branch deletion
-> next issue
```

Do not combine roadmap items simply because they touch adjacent packages. The sequence keeps security boundaries, regressions, and migration effects reviewable.

## Before coding

Codex must:

1. confirm it is on the maintainer-supplied branch and not `main`;
2. read `AGENTS.md`, `CONTRIBUTING.md`, and the assigned issue in full;
3. read the relevant ADRs, architecture sections, roadmap, and v0.1 scope;
4. inspect affected implementation, tests, public examples, and known issues;
5. state a concise plan;
6. identify contracts and safety boundaries that must remain unchanged;
7. verify unstable SDK, protocol, and platform facts against official primary sources.

If scope or instructions materially conflict, stop and ask rather than silently broadening the task.

## Repository operations

Codex must not:

- create, switch, rename, delete, or force-update branches unless the issue explicitly authorizes it;
- push or commit directly to `main`;
- use a GitHub App or integration ruleset bypass to avoid the pull request and required checks;
- create tags, releases, or repository settings;
- change remotes or rewrite shared history;
- merge or enable auto-merge for its own security-sensitive work.

A technical permission is not project authorization.

## Scope rules

- One issue maps to one branch and one pull request.
- No drive-by refactors or unrelated formatting.
- No dependency upgrade outside a dependency task.
- No public rename outside an identity task.
- No copied source from reference projects without an explicit license-compatible decision.
- No speculative abstraction that is not required by the issue.
- No dead compatibility layer without a documented migration need.
- Preserve deterministic output ordering and current public contracts unless the issue says otherwise.

## Security review questions

Answer these before implementation and again before review:

1. Can model-controlled input reach a shell, code interpreter, arbitrary path, arbitrary method, or arbitrary infrastructure endpoint?
2. Can a remote body, error, log, environment value, credential, path, prompt, or tool result escape into public output, telemetry, MCP, or persisted artifacts?
3. Are steps, attempts, concurrency, time, bytes, and collection sizes bounded?
4. Can DNS, redirects, or alternate resolved addresses create an SSRF or rebinding path?
5. Does tool metadata truthfully describe read-only, idempotent, destructive, open-world, risk, cost, and sensitivity behavior?
6. Is execution success separate from observed target health?
7. Can a finding or conclusion exist without valid evidence references?
8. Does the change duplicate Astralith or Kube-Sentinel responsibilities?
9. Does every failure mode remain safe, observable, cancellable, and testable?

Do not treat a written invariant as evidence that the existing code already implements it. Link known baseline gaps and avoid expanding them.

## Testing expectations

Where applicable, cover:

- normal success and unhealthy observations;
- malformed JSON, unknown fields, and trailing values;
- empty, oversized, and semantically invalid values;
- count, byte, retry, concurrency, and timeout boundaries;
- cancellation and shutdown;
- missing configuration, permission errors, and unavailable dependencies;
- non-2xx, malformed, and adversarial remote responses;
- deterministic ordering;
- secret, token, header, path, and prompt-injection strings;
- CLI, Agent Registry, MCP, telemetry, and persistence exposure.

Prefer table-driven tests and deterministic local fakes or servers. The default test suite must not depend on external network services, real credentials, or production infrastructure.

## Validation

Run focused tests while iterating. Before review, run and report the exact result of:

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./cmd/vardiel
```

If a check cannot run, state why and what remains unverified.

## Pull request and review

Use `.github/pull_request_template.md`. The pull request must distinguish implemented behavior from future work and include scope, non-goals, contracts, security/privacy impact, tests, validation, AI assistance, and unresolved uncertainty.

- Keep the pull request in draft while implementation or validation is incomplete.
- Do not present an automated review as independent human approval.
- Address comments with new commits rather than rewriting review history.
- Re-run affected tests after review changes.
- Prefer squash merge for focused task branches unless commit structure has durable value.
- Delete the task branch after merge.

## Current implementation sequence

The identity migration was completed in pull request #24. Continue in order:

1. issue #19 — safe public error boundary;
2. issue #20 — Tool Contract v2;
3. issue #21 — evidence and local case bundles;
4. Linux host and systemd diagnostics;
5. deterministic analyzers;
6. typed diagnostic plans;
7. three built-in playbooks;
8. model factory and OpenAI-compatible adapter;
9. `v0.1.0` release preparation.

Do not skip the public-error and tool-contract foundations to reach feature work faster.
