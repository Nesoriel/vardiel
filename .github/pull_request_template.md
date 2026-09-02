## Summary

<!-- What user or maintainer problem does this pull request solve? -->

## Linked issue

Closes #

## Change type

- [ ] Bug fix
- [ ] Feature
- [ ] Security hardening
- [ ] Documentation or contributor experience
- [ ] Tests, build, or CI
- [ ] Refactor with no intended behavior change

## Scope

<!-- Exact implemented changes. -->

## Explicit non-goals

<!-- What this pull request intentionally does not change. -->

## Architecture and public contracts

<!-- Affected packages, ADRs, CLI, JSON, MCP, environment variables, schemas, reports, or compatibility. Write "none" with a reason when applicable. -->

## Security and privacy impact

<!-- Discuss model-controlled input, shell/code execution, SSRF, credentials, raw errors, logs, output projection, telemetry, permissions, timeouts, and bounds. -->

- [ ] This change does not add arbitrary model-controlled shell or code execution.
- [ ] Tool and protocol metadata remain truthful.
- [ ] This change does not introduce a new raw-error, credential, path, token, or sensitive-output disclosure path.
- [ ] Existing SSRF, redirect, TLS, endpoint, size, count, timeout, cancellation, and redaction boundaries are preserved.
- [ ] Relevant attacker-controlled strings and failure paths are tested.
- [ ] Known pre-existing limitations are linked rather than presented as fixed.

Explain any unchecked item or mark it not applicable with a reason.

## Tests

<!-- New and changed tests, including unhealthy-state, failure, boundary, cancellation, privacy, and deterministic-output cases. -->

## Validation

<!-- Paste only commands actually run and concise results. State why any expected check was skipped. -->

```text
go mod tidy

git diff --exit-code -- go.mod go.sum

gofmt -w .

test -z "$(gofmt -l .)"

go vet ./...

go test -race -coverprofile=coverage.out ./...

go build ./cmd/vardiel
```

## Evidence and output examples

<!-- For behavior changes, include sanitized representative JSON, CLI, MCP, or report output. -->

## AI-assisted development

- [ ] No material AI assistance was used.
- [ ] AI tools materially assisted this change.

Tools used:

Human review and verification performed:

<!-- Do not include private prompts, chain-of-thought, credentials, or sensitive repository data. -->

## Migration, rollout, and rollback

<!-- Required for breaking, deployment, persistence, or mutation-related changes. Otherwise explain why not applicable. -->

## Unresolved risks or uncertainty

<!-- Remaining limitations, unverified assumptions, follow-up issues, or "none identified" with justification. -->

## Review checklist

- [ ] I read `CONTRIBUTING.md` and the applicable `AGENTS.md`.
- [ ] The pull request addresses one focused issue.
- [ ] Unrelated formatting, dependency, naming, and refactor changes are excluded.
- [ ] Documentation and examples distinguish current behavior from roadmap work.
- [ ] I reviewed generated changes and source attribution.
- [ ] I did not commit secrets, production data, binaries, coverage output, or local configuration.
