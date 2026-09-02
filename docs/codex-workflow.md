# Codex execution workflow

This document defines how maintainers delegate, Codex executes, and reviewers merge tasks in Vardiel. Shared contribution rules live in `CONTRIBUTING.md`; repository-wide coding-agent constraints live in `AGENTS.md`.

## Operating model

Vardiel uses an execution model in which the maintainer concentrates on product direction, architecture, priorities, and review while Codex performs routine repository work.

```text
accepted issue or explicit task
-> Codex aligns or creates the task branch
-> plan
-> implementation and tests
-> commits
-> push
-> draft pull request
-> CI and review
-> follow-up commits
-> ready for review
-> merge when explicitly delegated
-> issue closure and branch cleanup when delegated
```

Do not combine roadmap items simply because they touch adjacent packages. The sequence keeps security boundaries, regressions, and migration effects reviewable without making the maintainer perform mechanical Git work.

## Delegation levels

### Standard task delegation

Assigning an issue or an explicit implementation task authorizes Codex to complete the routine lifecycle:

- inspect repository and worktree state;
- fetch `origin`;
- select, create, or switch to the task branch under the branch rules below;
- fast-forward a non-diverged task branch;
- edit, test, and document the scoped change;
- create commits;
- push the task branch;
- open or update a draft pull request;
- mark the pull request ready after implementation and validation complete;
- push additional commits in response to review feedback.

Codex should not stop merely because the current workspace is on an obsolete clean branch. It should safely align the workspace to the assigned task.

### Completion delegation

Merge is not implied by standard task delegation. The maintainer or issue may explicitly authorize Codex to complete through merge.

With completion delegation, Codex may:

- merge the pull request after required checks pass and required conversations are resolved;
- close the linked issue when the pull request does not close it automatically;
- delete the completed remote task branch;
- report the final merge commit and post-merge CI state.

Completion delegation does not authorize ruleset bypass, direct `main` pushes, self-approval, repository settings, releases, or destructive cleanup. A security-sensitive change should receive an explicit final merge instruction after the maintainer has had the opportunity to review its scope and risk.

## Branch preparation

Use branch selection in this order:

1. a branch explicitly named by the maintainer;
2. a branch named in the assigned issue;
3. an existing branch clearly linked to the issue;
4. otherwise, create `<type>/<issue-number>-<short-slug>` from current `origin/main`.

Allowed types are:

```text
security
fix
feat
docs
test
ci
build
refactor
chore
```

Before changing branches, inspect:

```bash
git status --short --branch
git remote -v
git log --oneline --decorate -n 5
git fetch origin
```

A normal safe alignment may include:

```bash
git switch <existing-local-task-branch>
git switch --track origin/<existing-remote-task-branch>
git switch -c <new-task-branch> origin/main
git merge --ff-only origin/<task-branch>
git push -u origin <new-task-branch>
```

These examples describe allowed outcomes, not a requirement to execute every command.

Stop and ask when:

- the worktree contains unrelated changes that may be overwritten;
- the current branch has unpushed commits of unknown ownership;
- the required exact branch is missing and the task does not authorize creating it;
- local and remote histories diverged;
- checkout or fast-forward would overwrite files;
- authentication or permissions prevent routine task operations.

Do not solve these conditions with `git reset --hard`, `git clean`, force push, or history rewriting. Do not automatically stash unknown work. A named stash is permitted only for changes Codex created during the current task, and it must be restored or removed before handoff.

## Before coding

Codex must:

1. reach the correct task branch and confirm it is not `main`;
2. read `AGENTS.md`, `CONTRIBUTING.md`, and the assigned issue in full;
3. read relevant ADRs, architecture sections, roadmap, and v0.1 scope;
4. inspect affected implementation, tests, public examples, and known issues;
5. state a concise implementation plan;
6. identify contracts and safety boundaries that must remain unchanged;
7. verify unstable SDK, protocol, and platform facts against official primary sources.

If scope or instructions materially conflict, stop and ask rather than silently broadening the task.

## Repository operations that remain prohibited

Routine task delegation never authorizes Codex to:

- push or commit directly to `main`;
- use a GitHub App or integration ruleset bypass;
- force push, rebase shared review history, or rewrite published commits;
- delete or rename unrelated branches;
- change remotes, rulesets, settings, secrets, permissions, visibility, or repository ownership;
- create tags or releases outside an explicitly assigned release task;
- approve its own pull request or present an automated review as independent human approval;
- merge without completion delegation.

Technical access is not permission for destructive or administrative work. Conversely, normal branch, commit, push, and pull-request operations are expected parts of implementation.

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

## Commit and pull-request workflow

Codex is expected to perform these routine steps itself:

1. review the diff and remove unrelated changes;
2. create one or more clear, reviewable commits;
3. push the task branch;
4. open a draft pull request linked to the issue;
5. keep the pull request body accurate as implementation evolves;
6. wait for or inspect CI;
7. fix failures with additional commits;
8. respond to review comments with additional commits;
9. re-run affected validation;
10. mark the pull request ready when complete.

Use `.github/pull_request_template.md`. The pull request must distinguish implemented behavior from future work and include scope, non-goals, contracts, security/privacy impact, tests, validation, AI assistance, and unresolved uncertainty.

Do not ask the maintainer to perform commits, pushes, branch setup, or pull-request creation when standard task delegation already authorizes Codex to do them safely.

## Review, merge, and cleanup

- Passing CI is necessary but not sufficient for security-sensitive work.
- Preserve review history; use follow-up commits rather than force-pushing rewritten history.
- Re-run affected tests after review changes.
- Under completion delegation, merge only when repository rules permit it, required checks pass, required discussions are resolved, and no unresolved risk requires a maintainer decision.
- Prefer squash merge for focused task branches unless commit structure has durable value.
- Delete the completed task branch when completion delegation includes cleanup.
- Report the merge commit and post-merge CI result.

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