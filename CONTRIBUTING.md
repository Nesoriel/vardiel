# Contributing to Vardiel

Thank you for helping improve Vardiel. Contributions may include bug reports, design discussion, documentation, tests, and code.

Participation in this repository is governed by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). Security vulnerabilities must be reported through [SECURITY.md](SECURITY.md), not through a public issue.

## Choose the right channel

- Use the bug report form for a reproducible defect.
- Use the feature request form for a user problem or proposed capability.
- Use the maintainer/Codex task template only for scoped implementation work that already has an accepted direction.
- For a security vulnerability, follow the private-first process in `SECURITY.md`.

Search existing issues before opening a new one. Do not include credentials, private infrastructure data, customer data, raw authorization headers, kubeconfig content, or sensitive logs in any public issue.

## Before writing code

For anything beyond a small typo or clearly isolated documentation correction, open or join an issue first. This avoids duplicate work and lets maintainers confirm scope, architecture ownership, and compatibility expectations.

Read:

1. `README.md`;
2. this file;
3. `AGENTS.md` if an AI coding agent will participate;
4. the assigned issue;
5. relevant ADRs and `docs/architecture.md`.

Vardiel deliberately separates responsibilities from Astralith and Kube-Sentinel. Vardiel owns the per-host incident and recovery loop; Astralith may own identity, inventory, fleet policy distribution, UI, and notification delivery; Kube-Sentinel owns Kubernetes controller remediation. A proposal that duplicates those control-plane or controller responsibilities should be discussed before implementation.

## Development setup

Use the Go version declared in `go.mod`.

```bash
git clone https://github.com/Nesoriel/vardiel.git
cd vardiel
go mod download
go test ./...
go build ./cmd/vardiel
```

Tests should run without external infrastructure. Some manual integrations require explicitly configured Docker, Kubernetes, Prometheus, Loki, or Ark access; never use production credentials in a contribution or test fixture.

## Branches and commits

Do not work directly on `main`. Use a focused branch such as:

```text
fix/<issue-number>-<short-description>
feat/<issue-number>-<short-description>
security/<issue-number>-<short-description>
docs/<issue-number>-<short-description>
test/<issue-number>-<short-description>
chore/<issue-number>-<short-description>
```

Keep commits reviewable and use an imperative summary. Conventional prefixes such as `feat:`, `fix:`, `security:`, `docs:`, `test:`, `refactor:`, `ci:`, `build:`, and `chore:` are encouraged but not mechanically required.

Do not:

- force-push shared or reviewed branches without explicit maintainer approval;
- mix dependency upgrades with unrelated work;
- reformat unrelated files;
- commit generated binaries, coverage output, `.env` files, secrets, or production configuration;
- create a compatibility layer without a documented migration need;
- discard or hide another contributor's unrecognized work.

## Coding standards

Vardiel is a security-sensitive Go application. Contributions should:

- remain formatted with `gofmt`;
- keep `internal/agent` provider-neutral;
- pass `context.Context` through cancellable operations;
- give every goroutine a bounded lifetime and shutdown path;
- use errors for expected failures and preserve wrapped internal causes without exposing unsafe public text;
- document exported APIs;
- use strict JSON decoding and semantic validation;
- keep output deterministic and bounded;
- use typed clients or fixed argument builders instead of arbitrary shell strings;
- keep known recovery paths deterministic and independent of model/provider availability;
- use event-driven Linux/systemd interfaces instead of tight polling where available;
- prefer the standard library unless a dependency materially reduces risk or implements an adopted integration.

Product features that execute arbitrary model-controlled shell commands or code are not accepted. A mutating action requires an accepted focused issue and ADR 0002's fixed-action boundary: explicit standing or per-incident authorization, strict targets and arguments, preconditions, a per-target lock, bounded attempts and concurrency, independent validation, audit, cooldown/circuit breaking, and rollback or explicit non-rollback handling. Current CLI and MCP tools remain read-only unless a later contract explicitly says otherwise. Read-only infrastructure access does not make the underlying credential, socket, or kubeconfig unprivileged.

## Tests

Add tests at the same abstraction level as the change. Relevant cases include success, unhealthy-but-successfully-observed state, malformed input, boundary values, timeout, cancellation, deterministic ordering, privacy, and failure paths. Recovery work must also cover event deduplication, policy denial/expiry/conflict, duplicate-action prevention, concurrency, failed validation, retry/cooldown/circuit breaking, crash recovery, and configured latency budgets.

Prefer table-driven tests and local deterministic fakes or servers. The default suite must not require internet access, real credentials, or production systems.

Before requesting review, run:

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./cmd/vardiel
```

Record the commands actually run. Do not report an unexecuted check as successful.

## Pull requests

Open one focused pull request per issue and complete the repository template. A good pull request:

- links the issue;
- explains user-visible and internal behavior;
- states explicit non-goals;
- identifies architecture, security, privacy, and compatibility impact;
- includes tests and exact validation results;
- updates documentation and examples when public behavior changes;
- names remaining uncertainty and follow-up work.

Keep the pull request in draft while behavior or validation is incomplete. Maintainers may request that a large change be split before review.

## AI-assisted contributions

AI-assisted development is welcome when it improves the contribution rather than replacing contributor responsibility.

The contributor remains accountable for every submitted line and must:

- review generated code and documentation;
- verify licenses and source attribution;
- run the applicable tests;
- disclose material AI assistance in the pull request;
- describe human verification performed;
- avoid sending repository secrets, private infrastructure data, or undisclosed third-party code to an external model.

You do not need to publish private prompts or chain-of-thought. A concise statement of the tools used and the verification performed is sufficient.

### Maintainer-delegated delivery automation

Maintainers may delegate the normal mechanical repository workflow to a coding agent. Subject to `AGENTS.md`, an assigned agent may prepare the task branch, commit, push, open and update the pull request, monitor and fix CI, mark the pull request ready, and respond to review comments.

This delegation is intentional: maintainers may focus on goals, architecture, risk, and acceptance rather than routine Git operations. The agent does not need a separate permission prompt for each ordinary branch, commit, push, or pull-request action when the task has been assigned.

Delegation does not remove accountability or protected-branch gates. A coding agent must not push directly to `main`, bypass repository rules, discard unknown work, rewrite shared history without explicit approval, or self-authorize merge. Merge may be delegated only through an explicit maintainer instruction and only after required checks and review conditions pass.

## Review and merge

Passing CI is necessary but not sufficient for a security-sensitive change. Maintainers review scope, architecture, threat boundaries, tests, and public contracts.

Automated or AI review is advisory and is not independent human approval. Review threads must be resolved without erasing prior context. Focused branches are normally squash-merged and deleted after merge.

When a maintainer explicitly delegates merge execution, the coding agent may perform the merge through the protected pull-request path, verify the resulting state, and delete only the completed task branch. Otherwise the agent should leave a review-ready pull request for the maintainer.

## Licensing

By submitting a contribution, you agree that it may be distributed under the repository's [Apache License 2.0](LICENSE). Contributors are responsible for identifying copied or adapted material and preserving required notices.
