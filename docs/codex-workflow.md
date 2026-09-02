# Codex execution workflow

This document defines how maintainers direct Codex tasks and how Codex performs the mechanical delivery lifecycle in Vardiel. Shared contribution rules live in `CONTRIBUTING.md`; repository-wide coding-agent constraints live in `AGENTS.md`.

## Operating model

The maintainer is responsible for deciding what should be built, why it matters, the acceptable architecture and risk boundary, and whether a gated change may merge.

Codex is responsible for carrying the assigned task through the normal repository workflow:

```text
accepted issue and direction
-> inspect repository state
-> prepare or select task branch
-> plan
-> implement and test
-> commit and push
-> open or update draft pull request
-> monitor and fix CI
-> respond to review
-> mark ready
-> merge only when explicitly authorized
-> verify and clean up its task branch
```

The maintainer should not have to switch branches, create routine task branches, write commits, push changes, or open pull requests merely to keep an assigned task moving.

## Before coding

Codex must:

1. inspect the current branch, worktree, remotes, and upstream relationship;
2. read `AGENTS.md`, `CONTRIBUTING.md`, and the assigned issue in full;
3. read the relevant ADRs, architecture sections, roadmap, and v0.1 scope;
4. inspect affected implementation, tests, public examples, and known issues;
5. state a concise implementation plan;
6. identify contracts and safety boundaries that must remain unchanged;
7. verify unstable SDK, protocol, and platform facts against official primary sources;
8. place the worktree on an appropriate task branch without discarding unknown work.

If scope or instructions materially conflict, stop and ask rather than silently broadening the task.

## Task-branch preparation

A task branch named by the maintainer or issue is authoritative and may be used without another confirmation.

Codex may perform the following standard preparation:

- `git fetch origin`;
- switch to an existing named local branch;
- create a local tracking branch from an existing `origin/<task-branch>`;
- create the named task branch from current `origin/main` when it does not exist;
- when no branch is named, derive one using `<type>/<issue-number>-<short-slug>`;
- set the upstream when first pushing;
- fast-forward a branch that has no divergent work;
- before first publication, rebase only its own unpushed local task commits onto current `origin/main`.

Codex must stop rather than improvise when the worktree contains unrelated or unrecognized changes, commit ownership is unclear, checkout would overwrite files, published histories have diverged, or a conflict requires product, architecture, or security judgment.

Do not use destructive cleanup to make preparation succeed. Unknown work must not be discarded, hidden, or overwritten. A named temporary stash is permitted only for clearly owned changes from the same assigned task, and it must be restored or removed before handoff.

## Commits and pushes

Codex is expected to make and push commits as part of completing a task.

- Use small logical commits with imperative summaries; conventional prefixes are encouraged.
- Commit only files within the assigned scope.
- Amend the latest commit only while it is the agent's own unpushed work.
- After a branch is pushed or reviewed, preserve history and use follow-up commits.
- Push only the assigned task branch and set its upstream when necessary.
- Never commit or push directly to `main`.
- Never use a GitHub App, administrator, integration, or ruleset bypass to avoid the pull-request and required-check path.
- Do not force-push a published or reviewed branch unless the maintainer explicitly authorizes that exact operation and reason.

## Pull-request lifecycle

Codex should own the routine pull-request workflow:

1. open a draft pull request linked to the issue;
2. keep the title and body synchronized with the actual implementation;
3. push updates as focused commits;
4. monitor required checks;
5. inspect failed jobs and logs;
6. correct failures that are within scope and push the fix;
7. document skipped or unavailable checks honestly;
8. mark the pull request ready when scope, tests, documentation, and validation are complete;
9. respond to review feedback with additional commits without erasing prior review context.

A draft pull request is preferred early for substantial work because it preserves CI and discussion context. Very small changes may open the pull request after the first complete commit, but they must still use the same protected review path.

## Merge authority

CI success is not merge authorization.

Codex may merge its own task pull request only when all of the following are true:

- the maintainer explicitly authorizes merge in the current instruction, assigned issue, or pull-request conversation;
- required checks pass;
- required review threads are resolved;
- the pull request still matches the assigned scope;
- no unresolved risk requires a maintainer decision;
- the merge uses the protected pull-request path rather than a bypass.

When merge is not explicitly authorized, Codex must stop at a review-ready pull request and provide a concise handoff.

After an authorized merge, Codex may verify the resulting default-branch commit, confirm linked-issue state, and delete only the completed task branch. It must not clean up unrelated branches, tags, or releases.

## Repository operations that still require explicit authorization

Routine task delivery does not authorize Codex to:

- change repository settings, rulesets, access controls, webhooks, secrets, or environments;
- change remotes or repository ownership;
- create or delete tags and releases;
- delete, rename, or rewrite unrelated branches;
- push to `main` or bypass protected-branch rules;
- use `git reset --hard`, `git clean`, or destructive operations on unknown work;
- force-push or rewrite shared history;
- merge without explicit maintainer authorization.

A technical permission is not project authorization.

## Scope rules

- One issue maps to one focused branch and one pull request.
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
8. For a mutating action, where is standing or per-incident authorization checked, and are target, blast radius, preconditions, lock, attempts, cooldown, circuit breaker, validation, audit, and rollback/non-rollback behavior explicit?
9. Can a known recovery complete within its local deadline without a model or remote control plane?
10. Does the change duplicate Astralith or Kube-Sentinel responsibilities?
11. Does every failure mode remain safe, observable, cancellable, and testable?
12. Can knowledge or deployment metadata introduce executable content, widen
    authority, escape retrieval bounds, or appear without source provenance?

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

## Pull-request content

Use `.github/pull_request_template.md`. The pull request must distinguish implemented behavior from future work and include scope, non-goals, contracts, security/privacy impact, tests, validation, AI assistance, migration or rollback needs, and unresolved uncertainty.

The final handoff must also state:

- branch and commit state;
- pull-request URL and draft/ready state;
- required-check status;
- whether merge is authorized, completed, or awaiting the maintainer;
- any repository operation intentionally not performed.

## Current implementation sequence

The identity migration was completed in pull request #24. ADR 0002 resets the
delivery order around a complete recovery loop:

1. issue #33 — accept and synchronize the autonomous Linux recovery direction;
2. issue #19 — finish the narrow safe public error boundary without broadening it;
3. create one issue for the `systemd-service-unavailable` daemon-to-validation vertical slice;
4. harden privilege separation, crash recovery, policy, audit, notification, and real-host performance evidence around that slice;
5. add resource-pressure, filesystem-pressure, and endpoint playbooks one complete slice at a time;
6. add bounded model reasoning for unknown incidents;
7. integrate fleet identity, inventory, policy distribution, and notification through Astralith.

Issue #20's universal Tool Contract v2 and issue #21's general case-bundle work
are not prerequisites for the first recovery slice. After issue #33 merges,
supersede or reshape them around contracts proven necessary by the working
systemd loop. Do not describe this target sequence as implemented behavior.
