# AGENTS.md

Repository-wide instructions for coding agents working on Vardiel.

## Scope and precedence

This root file applies to the entire repository. A more deeply nested `AGENTS.md` may add package-specific instructions for its subtree, but it must not weaken the safety invariants in this file or an accepted architecture decision record (ADR).

Use this precedence when repository instructions differ:

1. the explicit maintainer request and assigned GitHub issue define the requested outcome, scope, task branch, and completion authority;
2. the nearest applicable `AGENTS.md` defines implementation constraints;
3. accepted ADRs and `docs/architecture.md` define architectural boundaries;
4. `CONTRIBUTING.md` defines the shared contributor workflow;
5. other documentation and examples describe current behavior.

Do not silently resolve a material conflict. Stop, describe it, and ask the maintainer. Treat written safety rules as requirements, not as proof that the current implementation already satisfies them; inspect the code and linked issues.

## Project intent

Vardiel is a Go-first, evidence-driven operations diagnostic agent for personal servers, homelabs, and small infrastructure environments.

Its default posture is bounded and read-only:

```text
collect structured evidence
-> make deterministic findings
-> explain uncertainty
-> propose, but do not silently execute, changes
```

Vardiel is the independent continuation of OpsPilot and is not a fork of Netiarius. Reference projects may inform design, but their source must not be copied without an explicit, license-compatible decision.

## Delegated repository work

The maintainer intends coding agents to complete routine repository work, not merely edit files. Assigning an issue or an explicit implementation task grants **standard task delegation** unless the maintainer narrows it.

Standard task delegation authorizes the coding agent to:

- inspect the current worktree, branch, remotes, and ahead/behind state;
- fetch current refs from `origin`;
- switch to the explicitly assigned existing task branch;
- create a local tracking branch for an already-existing `origin/<task-branch>`;
- when no task branch exists, create and push one deterministic issue branch from current `origin/main`;
- fast-forward a local task branch to its upstream when histories have not diverged;
- implement the assigned scope and add tests and documentation;
- create clear commits containing only its task changes;
- push the task branch;
- open and update a draft pull request linked to the issue;
- mark the pull request ready after implementation and required validation are complete;
- respond to review feedback with additional commits and update the pull request description.

Branch selection order is:

1. a branch explicitly named by the maintainer;
2. a branch named in the assigned issue;
3. an existing branch clearly linked to the issue;
4. otherwise, a deterministic new branch using `<type>/<issue-number>-<short-slug>`, where type is one of `security`, `fix`, `feat`, `docs`, `test`, `ci`, `build`, `refactor`, or `chore`.

A stale or obsolete current branch is not a reason to stop when the worktree can be safely aligned to the assigned task branch.

Before changing branches, inspect at least:

```bash
git status --short --branch
git remote -v
git log --oneline --decorate -n 5
```

Do not overwrite or discard unrelated work. Stop and report the exact condition when:

- the worktree contains changes not created for the current task and checkout may affect them;
- the current branch contains unpushed commits of unknown ownership;
- the assigned remote branch is missing when the task requires that exact existing branch;
- local and remote task-branch histories have diverged;
- checkout or fast-forward would overwrite files;
- repository access is insufficient for the required routine operation.

Do not use `git reset --hard`, `git clean`, force push, or history rewriting to make alignment succeed. Do not automatically stash unknown or pre-existing work. A coding agent may use a named stash only for changes it created during the current task, and it must restore or remove that stash before handoff.

**Completion delegation** is separate from standard task delegation. A coding agent may merge a pull request, close the linked issue, and delete the completed task branch only when the maintainer or assigned issue explicitly authorizes completion through merge. Merge authorization is valid only after required checks pass, required review threads are resolved, the branch is current under repository rules, and no unresolved risk requires maintainer judgment.

Routine task delegation never authorizes:

- direct commits or pushes to `main`;
- use of a GitHub App or integration ruleset bypass;
- force push or rewriting shared history;
- deletion or renaming of unrelated branches;
- changing remotes, repository settings, rulesets, secrets, permissions, or visibility;
- creating tags or releases unless the task explicitly covers release work;
- approving the agent's own pull request or representing automated analysis as independent human approval.

A technical capability is not permission for destructive or administrative work, but routine Git and pull-request operations are expected parts of an assigned implementation task.

## Required preflight

Before editing:

1. identify the assigned issue and branch under the delegation rules above;
2. safely align the worktree to the task branch and confirm it is not `main`;
3. read the assigned issue in full;
4. read this file, `CONTRIBUTING.md`, the relevant ADRs, and affected architecture sections;
5. inspect the current implementation, tests, public examples, and known follow-up issues;
6. state a concise plan, affected contracts, and safety boundaries;
7. verify unstable SDK, protocol, and platform facts against official primary sources.

## Change discipline

- One issue maps to one focused branch and one pull request.
- Keep drive-by refactors, unrelated formatting, dependency upgrades, generated artifacts, and speculative abstractions out of the change.
- Do not add compatibility layers without a documented user or integration need.
- Preserve public contracts unless the issue explicitly authorizes a breaking change.
- Update documentation, examples, schemas, and migration notes when public behavior changes.
- Prefer small, reviewable commits. Do not rewrite shared history after a branch has been pushed or reviewed.
- Do not commit credentials, `.env` files, local case data, coverage output, binaries, editor state, or production configuration.
- Do not describe roadmap work as implemented behavior.

## Non-negotiable safety invariants

These constraints apply to product code, tests, examples, and documentation:

- Do not add a product capability that executes arbitrary shell commands or model-generated code.
- Model output, tool arguments, provider responses, remote API data, logs, files, and persisted case content are untrusted input.
- Infrastructure tools remain read-only unless a later issue introduces deterministic policy evaluation, explicit approval, bounded execution, post-change validation, and rollback planning.
- Validate every tool name and argument with strict schemas and semantic checks. Reject unknown fields, trailing JSON values, unsafe paths, unsupported methods, and ambiguous identifiers.
- Do not expose raw provider or infrastructure error bodies, credentials, authorization data, environment values, kubeconfig contents, ServiceAccount tokens, private paths, or complete sensitive tool output through the model context, CLI, MCP, telemetry, case bundles, reports, or tests.
- Preserve SSRF defenses, DNS re-resolution checks, redirect limits, TLS verification, fixed endpoint allowlists, response-byte limits, result-count limits, step budgets, timeouts, and cancellation.
- Do not expose arbitrary infrastructure API paths, HTTP methods, PromQL, LogQL, Kubernetes resource types or selectors, Docker endpoints, or filesystem paths as model-controlled arguments.
- Keep protocol streams clean. MCP stdio uses stdout only for protocol frames.
- Treat Docker sockets, kubeconfigs, tokens, and local control sockets as privileged capabilities even when exposed operations are read-only.
- Never weaken a boundary merely to make a test pass.

When a known baseline limitation conflicts with a target invariant, preserve the limitation's tracking issue, avoid extending it, and fix it only in the issue assigned for that purpose.

## Architecture ownership

Vardiel owns:

- provider-neutral Agent Runtime and model adapters;
- bounded tool packs and infrastructure adapters;
- diagnostic orchestration, typed plans, playbooks, analyzers, evidence, findings, reports, and local case bundles;
- policy metadata, non-executable action proposals, CLI, MCP, and machine-readable contracts.

Vardiel does not own:

- host inventory, Ansible scheduling, web authentication, GitOps apply records, centralized persistence, or approval UI; those belong to Astralith;
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
- concurrency with the race detector.

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

## Pull request lifecycle

Under standard task delegation, the coding agent is expected to commit, push, open or update the pull request, keep its description accurate, and push review fixes without asking the maintainer to perform routine Git steps.

- Link the pull request to the assigned issue and use `.github/pull_request_template.md`.
- Keep the pull request in draft while scope, implementation, or validation is incomplete.
- Mark it ready when the required work and checks are complete.
- Re-run affected tests after review changes.
- Preserve review context by adding commits rather than force-pushing rewritten history.
- Do not enable auto-merge unless completion delegation explicitly permits it.
- Merge and remote task-branch deletion follow the completion-delegation rules above.

## Completion handoff

The final handoff and pull request description must include:

- implemented behavior and explicit non-goals;
- files and public contracts changed;
- security, privacy, and compatibility impact;
- tests added or updated;
- exact validation commands and concise results;
- skipped checks, unresolved uncertainty, and follow-up issues;
- commit, pull-request, CI, and merge state;
- material AI assistance and the human verification still required.

Human contribution and conduct guidance lives in `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md`. Vulnerabilities must follow `SECURITY.md`.