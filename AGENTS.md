# AGENTS.md

Repository-wide instructions for coding agents working on Vardiel.

## Scope and precedence

This root file applies to the entire repository. A more deeply nested `AGENTS.md` may add package-specific instructions for its subtree, but it must not weaken the safety invariants in this file or an accepted architecture decision record (ADR).

Use this precedence when repository instructions differ:

1. the maintainer's explicit instruction and the assigned GitHub issue define the requested outcome, scope, delivery branch, and merge authority;
2. the nearest applicable `AGENTS.md` defines implementation constraints;
3. accepted ADRs and `docs/architecture.md` define architectural boundaries;
4. `CONTRIBUTING.md` defines the shared contributor workflow;
5. other documentation and examples describe current behavior.

Do not silently resolve a material conflict. Stop, describe it, and ask the maintainer. Treat written safety rules as requirements, not as proof that the current implementation already satisfies them; inspect the code and linked issues.

## Project intent

Vardiel is a Go-first, fast, low-interruption Linux incident-response agent. Its target product loop is:

```text
detect locally
-> diagnose deterministically
-> execute an explicitly pre-authorized fixed action when policy permits
-> validate recovery
-> close quietly or escalate with a precise reason
```

Known recovery paths must not depend on a language-model or remote-control-plane round trip. A model may assist unknown diagnosis, but it never grants authority or invents executable operations.

The current implementation is still bounded and read-only. Mutating behavior may be introduced only by focused issues that implement [ADR 0002](docs/adr/0002-fast-autonomous-linux-recovery.md), including standing authorization, typed actions, preconditions, locks, budgets, validation, cooldowns, circuit breakers, audit, and rollback or explicit non-rollback handling.

Vardiel is the independent continuation of OpsPilot and is not a fork of Netiarius. Reference projects may inform design, but their source must not be copied without an explicit, license-compatible decision. Do not adopt arbitrary model-generated Python or shell execution.

## Maintainer and coding-agent roles

The maintainer owns product direction, priorities, issue scope, architecture and security decisions, resolution of material ambiguity, and final merge authorization when a gate is required.

For an accepted task, the coding agent owns the normal mechanical delivery lifecycle by default: repository inspection, task-branch preparation, implementation, tests, focused commits, pushes, pull-request creation and updates, CI follow-up, and review-response commits. The maintainer should not need to perform routine Git operations merely to let the agent continue.

Passing CI does not let an agent self-authorize a merge. Merge authority must come from an explicit maintainer instruction in the current task, issue, or pull-request conversation.

## Required preflight

Before editing:

1. inspect the current branch, worktree status, remotes, and upstream relationship;
2. read the assigned issue in full;
3. read this file, `CONTRIBUTING.md`, the relevant ADRs, and the affected architecture sections;
4. inspect the current implementation, tests, public examples, and known follow-up issues;
5. state a concise plan, affected contracts, and safety boundaries;
6. verify unstable SDK, protocol, and platform facts against official primary sources;
7. align the worktree to the task branch using the authority and safeguards below.

## Task-branch and delivery authority

A named branch in the maintainer instruction or assigned issue is authorization to use that branch. No second permission prompt is required for routine alignment.

For the assigned task, the coding agent may:

- fetch current refs from `origin`;
- switch to the named existing local task branch;
- create a local tracking branch when the named `origin/<task-branch>` already exists;
- create the named task branch from current `origin/main` when it does not yet exist locally or remotely;
- derive a conventional task branch when no name is supplied, using `<type>/<issue-number>-<short-slug>` where practical;
- fast-forward a task branch that has no divergent work;
- rebase the agent's own unpushed local commits onto current `origin/main` before the first push;
- commit focused logical changes and amend only the agent's own unpushed latest commit;
- push the task branch and set its upstream;
- open and update a draft pull request linked to the issue;
- monitor required checks, inspect failures, fix the task branch, and push follow-up commits;
- mark the pull request ready after scope and validation are complete;
- respond to review feedback with additional commits;
- merge the pull request and delete only its own task branch when the maintainer explicitly authorizes merge and all repository gates pass.

Before switching, creating, rebasing, or updating a branch, verify that no unknown work will be lost. Stop and report the exact condition when:

- the worktree contains unrelated or unrecognized changes;
- the current branch contains commits whose ownership or destination is unclear;
- local and remote task-branch history have diverged after publication or review;
- checkout, rebase, or update would overwrite files;
- a merge conflict requires product, architecture, or security judgment;
- required credentials or permissions are unavailable.

A named safety stash is allowed only for changes clearly created by the same agent for the same assigned task. Record it, restore it promptly, and do not leave hidden work at handoff. Never stash, discard, or overwrite unrecognized user work.

The coding agent must not:

- commit or push directly to `main`;
- use a GitHub App, integration, administrator, or ruleset bypass to avoid the pull-request and required-check flow;
- force-push or rewrite a published or reviewed branch unless the maintainer explicitly authorizes that exact operation and reason;
- use `git reset --hard`, `git clean`, or another destructive command on unknown work;
- change remotes, repository settings, rulesets, access controls, tags, or releases unless the assigned task explicitly requires it;
- delete, rename, or rewrite unrelated branches;
- merge merely because checks pass, without explicit merge authorization.

After a pull request is open, prefer non-destructive follow-up commits. If the base branch advances, merge or otherwise update it without erasing review history; ask before any operation that would require a force-push.

## Change discipline

- One issue maps to one focused branch and one pull request.
- Keep drive-by refactors, unrelated formatting, dependency upgrades, generated artifacts, and speculative abstractions out of the change.
- Do not add compatibility layers without a documented user or integration need.
- Preserve public contracts unless the issue explicitly authorizes a breaking change.
- Update documentation, examples, schemas, and migration notes when public behavior changes.
- Prefer small, reviewable commits. Do not rewrite shared history.
- Do not commit credentials, `.env` files, local case data, coverage output, binaries, editor state, or production configuration.
- Do not describe roadmap work as implemented behavior.

## Non-negotiable safety invariants

These constraints apply to product code, tests, examples, and documentation:

- Do not add a product capability that executes arbitrary shell commands or model-generated code.
- Model output, tool arguments, provider responses, remote API data, logs, files, and persisted case content are untrusted input.
- Current infrastructure tools remain read-only. A later issue may add a separate fixed action only with deterministic policy evaluation, explicit standing authorization or per-incident approval, bounded execution, post-change validation, audit, cooldown/circuit breaking, and rollback or an explicit non-rollback classification.
- Validate every tool name and argument with strict schemas and semantic checks. Reject unknown fields, trailing JSON values, unsafe paths, unsupported methods, and ambiguous identifiers.
- Do not expose raw provider or infrastructure error bodies, credentials, authorization data, environment values, kubeconfig contents, ServiceAccount tokens, private paths, or complete sensitive tool output through the model context, CLI, MCP, telemetry, case bundles, reports, or tests.
- Preserve SSRF defenses, DNS re-resolution checks, redirect limits, TLS verification, fixed endpoint allowlists, response-byte limits, result-count limits, step budgets, timeouts, and cancellation.
- Do not expose arbitrary infrastructure API paths, HTTP methods, PromQL, LogQL, Kubernetes resource types or selectors, Docker endpoints, or filesystem paths as model-controlled arguments.
- Keep protocol streams clean. MCP stdio uses stdout only for protocol frames.
- Treat Docker sockets, kubeconfigs, tokens, and local control sockets as privileged capabilities even when the exposed operations are read-only.
- Never weaken a boundary merely to make a test pass.

When a known baseline limitation conflicts with a target invariant, preserve the limitation's tracking issue, avoid extending it, and fix it only in the issue assigned for that purpose.

## Architecture ownership

Vardiel owns:

- provider-neutral Agent Runtime and model adapters;
- the per-host event, incident, diagnosis, recovery, validation, and circuit-breaker state machine;
- bounded observation tools and fixed typed action implementations;
- versioned playbooks, standing-policy evaluation, sanitized local incident audit, CLI, MCP, and machine-readable contracts.

Vardiel does not own:

- fleet inventory, Ansible scheduling, web authentication, centralized policy distribution, notification delivery, GitOps apply records, centralized persistence, or approval UI; those belong to Astralith;
- Alertmanager ingestion, Kubernetes CRDs, controller reconciliation, or controller-driven remediation; those belong to Kube-Sentinel.

Integrate through explicit MCP, API, or JSON contracts rather than duplicating another project's control plane.

## Go engineering rules

- Keep `internal/agent` provider-neutral. Provider SDKs belong under `internal/models/<provider>`.
- Pass `context.Context` through blocking or cancellable work; do not replace caller cancellation with a background context.
- Ensure every goroutine has an explicit lifetime, cancellation path, bounded work, and testable shutdown behavior.
- Return errors for expected failures; reserve panics for impossible programmer errors. Wrap internal causes with `%w` while keeping public error content sanitized.
- Add doc comments for exported packages, types, functions, methods, constants, and variables.
- Use `gofmt`. Keep imports and output ordering deterministic.
- Built-in tools must strictly decode JSON, reject unknown fields and trailing values, perform semantic validation, and return bounded structured JSON.
- Separate tool execution status from the observed health of the target system.
- Findings and conclusions must cite valid evidence identifiers.
- Prefer event-driven systemd, Journal, PSI, or kernel interfaces to tight polling for the local fast path.
- Keep known recovery deterministic and independent of model/provider availability.
- Fixed actions must recheck policy and preconditions under a per-target lock, enforce timeout and attempt budgets, validate health independently, and stop on cooldown or an open circuit breaker.
- Prefer fixed argument builders and typed clients over shell strings.
- Prefer the standard library. Add a dependency only when the issue justifies the risk and the pull request documents the choice.
- Bound network calls, retries, concurrency, file reads, API bodies, collection sizes, and free-text fields.
- Avoid package-level mutable state unless ownership and concurrency are explicit.

## Testing expectations

Add focused tests for behavior changed by the issue. Where applicable, cover:

- normal success and expected unhealthy observations;
- malformed input, unknown fields, trailing values, and semantic boundaries;
- timeout, cancellation, unavailable dependencies, and permission failures;
- deterministic ordering and size/count limits;
- secret, token, path, header, remote-error, and prompt-injection strings;
- CLI, Agent Registry, MCP, telemetry, and storage exposure;
- concurrency with the race detector;
- event deduplication, per-target serialization, retry/cooldown/circuit-breaker behavior, and crash-safe duplicate-action prevention;
- policy denial, expiry, conflicting grants, unavailable authorization, rollback or explicit non-rollback handling, and failed post-action validation;
- configured detection, action-start, and total-recovery deadlines without flaky production-performance claims.

Use table-driven tests and deterministic local fakes or test servers. The default test suite must not require internet access, production infrastructure, real credentials, or a writable Docker/Kubernetes environment.

Run focused package tests while iterating. Before review, run:

```bash
go mod tidy
git diff --exit-code -- go.mod go.sum
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test -race -coverprofile=coverage.out ./...
go build ./cmd/vardiel
```

Do not claim a command passed unless it was actually executed. If a check cannot run, state exactly why and what remains unverified.

## Pull-request and merge behavior

- Open one pull request linked to the assigned issue and use `.github/pull_request_template.md`.
- The coding agent should create the pull request as draft early enough to preserve CI and review context, then keep its body current.
- Keep the pull request in draft until implementation, documentation, and required validation are complete.
- Do not present automated analysis as independent human review or approve your own work.
- Address review comments with additional commits and re-run affected tests.
- Mark the pull request ready when it meets the issue acceptance criteria.
- Without explicit merge authorization, stop at a review-ready pull request and hand control back to the maintainer.
- With explicit merge authorization, wait for required checks and thread resolution, merge through the protected pull-request path, verify the resulting state, and delete only the completed task branch when appropriate.

## Completion handoff

The final handoff and pull-request description must include:

- implemented behavior and explicit non-goals;
- files and public contracts changed;
- security, privacy, and compatibility impact;
- tests added or updated;
- exact validation commands and concise results;
- commits pushed, pull-request state, and CI status;
- skipped checks, unresolved uncertainty, and follow-up issues;
- the source of merge authorization, or a clear statement that the pull request is awaiting maintainer review.

Human contribution and conduct guidance lives in `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md`. Vulnerabilities must follow `SECURITY.md`.
