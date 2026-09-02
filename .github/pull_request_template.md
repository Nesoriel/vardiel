## Summary

<!-- What problem does this pull request solve? -->

## Linked issue

Closes #

## Scope

<!-- List the exact implemented changes. -->

## Explicit non-goals

<!-- State what this pull request intentionally does not change. -->

## Architecture impact

<!-- Affected packages, contracts, ownership boundaries, or ADRs. -->

## Security and privacy impact

<!-- Address model-controlled input, shell/code execution, SSRF, credentials, raw errors, logs, output projection, telemetry, timeouts, and bounds. Write "none" only with a reason. -->

- [ ] No arbitrary shell or model-generated code execution was introduced.
- [ ] Read-only and mutation annotations remain truthful.
- [ ] Public errors and outputs do not expose raw infrastructure or provider data.
- [ ] Existing SSRF, redirect, TLS, endpoint, size, count, timeout, and redaction boundaries were preserved.
- [ ] Secret/path/token/prompt-injection regression cases were considered.

## Public contract and migration impact

<!-- CLI, environment variables, JSON schema, MCP metadata, case schema, binary names, or compatibility. -->

## Tests

<!-- New and changed tests, including failure and privacy paths. -->

## Validation

<!-- Paste the exact commands run and concise results. Do not claim unexecuted checks. -->

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

<!-- For behavior changes, include sanitized representative JSON or report output. -->

## Unresolved risks or uncertainty

<!-- State remaining limitations, follow-up issues, or "none identified" with justification. -->

## Review checklist

- [ ] I read `AGENTS.md` and the accepted identity ADR.
- [ ] The pull request addresses one focused issue.
- [ ] Unrelated formatting, dependency, naming, and refactor changes are excluded.
- [ ] Tool names and arguments are validated in code.
- [ ] Execution success is not confused with observed system health.
- [ ] Findings and conclusions cite evidence IDs where applicable.
- [ ] Documentation and examples match implemented behavior.
- [ ] Roadmap items are not described as already implemented.
