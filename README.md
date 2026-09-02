# Vardiel

Vardiel is a fast, low-interruption Linux incident-response agent implemented in Go. Its goal is to detect, diagnose, and safely recover common server failures locally, then notify an operator only when policy, risk, uncertainty, or failed validation requires attention.

> Status: early development. Current code is still read-only and request-driven: it includes a bounded Agent Runtime, an Ark Responses API adapter, an MCP stdio server, privacy-safe runtime events, optional OpenTelemetry tracing, and machine-readable network, Docker, Kubernetes, Prometheus, and Loki diagnostics. The always-on daemon and automatic recovery loop described below are accepted direction, not implemented behavior.

## Design goals

- Event-driven local detection and verified recovery with measurable deadlines.
- A deterministic fast path for known failures with no model or control-plane round trip.
- Fixed typed actions under explicit standing authorization, strict bounds, validation, audit, cooldowns, and circuit breakers.
- Bounded model assistance for unknown incidents without model-generated shell, code, actions, or authority.
- One per-host runtime that works alone and can integrate with Astralith for fleet and multi-user coordination.
- Structured, sanitized, deterministic contracts rather than prompt-only safety.

## Current capabilities

- Provider-neutral agent loop with bounded steps and structured tool errors.
- Volcengine Ark Responses API adapter through CloudWeGo Eino `agenticark`.
- Official MCP Go SDK stdio server backed by the same tool registry as the CLI and Agent Runtime.
- Strongly defined tool registry with duplicate and schema validation.
- Read-only `dns_lookup`, SSRF-aware `http_probe`, and certificate-aware `tls_inspect` tools.
- Read-only Docker Engine, container-list, and redacted container-inspect diagnostics over a local Unix socket.
- Read-only Kubernetes server, node, Pod-list, and redacted Pod-inspect diagnostics through client-go v0.36.2.
- Read-only Prometheus build/runtime, active-target, and constrained metric-snapshot diagnostics through fixed `/api/v1` endpoints.
- Read-only Loki readiness, build-info, and constrained stream-summary diagnostics without reading log lines.
- Shared network guard that resolves and validates every dial target before connecting.
- Machine-readable CLI intended for agents and automation.
- JSONL lifecycle events with run IDs, step numbers, durations, and sanitized error classes.
- Optional OTLP/HTTP traces for Agent runs, model calls, and tool executions.
- Configuration validation and provider-error secret redaction.

## Target recovery loop

The target architecture uses two paths:

```text
Linux event
-> deterministic incident correlation
-> known playbook + standing policy
-> fixed action
-> independent health validation
-> quiet close or precise escalation
```

Known recoveries stay local and do not call a model. Unknown incidents may use a
bounded model to rank hypotheses or select from registered playbooks, but the
same deterministic policy engine controls every action.

The first planned vertical slice is one explicitly allowlisted systemd service:
observe a failed unit, collect minimal unit/listener/endpoint facts, perform one
pre-authorized restart, validate stable health, close quietly on success, and
send one bounded HTTPS webhook escalation on timeout, failed validation,
exhausted attempts, or flapping. The target budgets are P95 event intake within
1 second, approved action start within 5 seconds, and verified recovery or
escalation within 30 seconds on the documented fixture.

Read [ADR 0002](docs/adr/0002-fast-autonomous-linux-recovery.md), the
[architecture](docs/architecture.md), and the [v0.1 scope](docs/v0.1-scope.md)
for the exact authority and safety boundary.

## Toolchain

Vardiel tracks the latest stable Go toolchain. The current baseline is Go 1.26.5.

```bash
go version
go mod tidy
go test ./...
go build ./cmd/vardiel
```

## Configure Ark

Copy `.env.example` into your preferred secret-management workflow and provide at least:

```bash
export ARK_MODEL_ID='ep-xxxxxxxxxxxxxxxx'
export ARK_API_KEY='your-api-key'
```

`ARK_THINKING` accepts `auto`, `enabled`, or `disabled`. Credentials are read from the environment and must not be committed.

## Run the agent

```bash
go run ./cmd/vardiel agent run \
  'Check Kubernetes Pods, Prometheus targets, and Loki stream availability.'
```

The command writes the final structured result to stdout. The Ark model can select from the registered read-only tools.

### Stream lifecycle events

```bash
go run ./cmd/vardiel agent run --events=jsonl \
  'Inspect example.com.' \
  2>events.jsonl
```

JSONL events are written to stderr, so agents can consume the final result from stdout independently. Events intentionally omit prompts, tool arguments, tool results, credentials, and raw provider errors.

### Export OpenTelemetry traces

```bash
export VARDIEL_OTEL_ENABLED=true
export OTEL_SERVICE_NAME=vardiel
export OTEL_EXPORTER_OTLP_ENDPOINT='http://localhost:4318'

go run ./cmd/vardiel agent run 'Inspect example.com.'
```

The OTLP/HTTP exporter follows standard OpenTelemetry environment variables. Telemetry initialization or shutdown failures do not fail the diagnostic run.

## Run as an MCP server

Build the binary and expose the registered read-only tools over stdio:

```bash
go build -o ./bin/vardiel ./cmd/vardiel
./bin/vardiel mcp stdio
```

A typical MCP client configuration is:

```json
{
  "mcpServers": {
    "vardiel": {
      "command": "/absolute/path/to/vardiel",
      "args": ["mcp", "stdio"],
      "env": {
        "VARDIEL_HTTP_ALLOW_PRIVATE": "false",
        "VARDIEL_TLS_ALLOW_PRIVATE": "false",
        "VARDIEL_DOCKER_SOCKET": "/var/run/docker.sock",
        "VARDIEL_KUBECONFIG": "/absolute/path/to/kubeconfig",
        "VARDIEL_KUBERNETES_CONTEXT": "production-readonly",
        "VARDIEL_PROMETHEUS_URL": "https://prometheus.example.com",
        "VARDIEL_PROMETHEUS_BEARER_TOKEN_FILE": "/absolute/path/to/prometheus-token",
        "VARDIEL_LOKI_URL": "https://loki.example.com",
        "VARDIEL_LOKI_BEARER_TOKEN_FILE": "/absolute/path/to/loki-token",
        "VARDIEL_LOKI_TENANT_ID": "operations"
      }
    }
  }
}
```

The server publishes each Registry tool with its existing JSON Schema and explicit read-only/idempotent annotations. Tool results are returned as text and, when the result is a JSON object, as MCP structured content. stdout is reserved exclusively for MCP protocol frames; warnings and failures go to stderr.

## Use the tool runtime directly

```bash
go run ./cmd/vardiel tool list

go run ./cmd/vardiel tool run dns_lookup \
  '{"host":"example.com"}'

go run ./cmd/vardiel tool run http_probe \
  '{"url":"https://example.com"}'

go run ./cmd/vardiel tool run tls_inspect \
  '{"host":"example.com","port":443}'

go run ./cmd/vardiel tool run docker_engine_info '{}'

go run ./cmd/vardiel tool run docker_container_list \
  '{"all":true,"limit":100}'

go run ./cmd/vardiel tool run docker_container_inspect \
  '{"container":"web"}'

go run ./cmd/vardiel tool run kubernetes_cluster_info \
  '{"node_limit":100}'

go run ./cmd/vardiel tool run kubernetes_pod_list \
  '{"namespace":"operations","limit":100}'

go run ./cmd/vardiel tool run kubernetes_pod_inspect \
  '{"namespace":"operations","pod":"web-0","event_limit":50}'

go run ./cmd/vardiel tool run prometheus_server_info '{}'

go run ./cmd/vardiel tool run prometheus_target_list \
  '{"limit":100}'

go run ./cmd/vardiel tool run prometheus_metric_snapshot \
  '{"metric":"up","matchers":{"job":"node"},"aggregation":"sum","group_by":["instance"],"limit":100}'

go run ./cmd/vardiel tool run loki_server_info '{}'

go run ./cmd/vardiel tool run loki_stream_summary \
  '{"matchers":{"namespace":"operations","service_name":"api"},"lookback_minutes":60,"limit":100}'
```

### Docker diagnostic boundary

`VARDIEL_DOCKER_SOCKET` defaults to `/var/run/docker.sock`. It accepts an absolute filesystem path or a `unix:///absolute/path` URI. Remote `tcp://`, `http://`, `https://`, `ssh://`, and relative targets are rejected.

The Docker client negotiates the daemon API through `/version`, then performs bounded GET requests. It does not invoke the Docker CLI and does not expose mutating operations.

`docker_container_inspect` returns a deliberate diagnostic projection rather than raw `docker inspect` output. It omits environment values, commands, arguments, raw labels, health-check output, Docker log paths, bind-mount and volume source paths, and free-text OCI/runtime errors.

Engine warning text is also omitted; only `warning_count` is returned. Container runtime errors are represented by `error_present` without returning the raw text.

Access to a Docker Unix socket is still a privileged host capability. Vardiel's read-only implementation does not turn the socket itself into a read-only security boundary. Only grant the process access to a trusted local socket, and do not mount that socket into untrusted containers.

### Kubernetes diagnostic boundary

Vardiel uses the official Kubernetes client-go v0.36.2. Kubernetes configuration is initialized lazily, so missing credentials do not prevent non-Kubernetes tools or the MCP server from starting.

When running outside a cluster, set `VARDIEL_KUBECONFIG` to an absolute kubeconfig path. `VARDIEL_KUBERNETES_CONTEXT` optionally selects a context. When running inside Kubernetes without an explicit kubeconfig, Vardiel uses the mounted ServiceAccount token and CA.

Before constructing a Kubernetes client, Vardiel rejects kubeconfigs that contain HTTP API servers, `insecure-skip-tls-verify`, proxy URLs, exec credential plugins, legacy auth-provider plugins, or user impersonation.

The model cannot provide a kubeconfig path, API server URL, arbitrary resource type, selector, API path, or HTTP method in tool arguments.

`kubernetes_pod_inspect` returns a deliberate projection rather than a raw Pod object. It omits environment values, commands, arguments, labels, annotations, volume source details, Secret and ConfigMap references, Pod logs, and all free-text condition, container-state, Pod-status, and Event messages. Event output is aggregated by `type` and `reason` only.

Apply the included minimum RBAC objects for an in-cluster deployment:

```bash
kubectl apply -f deploy/kubernetes/vardiel-readonly-rbac.yaml
```

The role grants GET on `/version`, GET/LIST on Nodes and Pods, and LIST on Events. It does not grant Secret access or the `pods/log` subresource.

### Prometheus diagnostic boundary

Set `VARDIEL_PROMETHEUS_URL` to the trusted Prometheus base URL. HTTPS is required by default. Internal HTTP endpoints require the explicit `VARDIEL_PROMETHEUS_ALLOW_HTTP=true` opt-in. Optional bearer authentication uses an absolute path in `VARDIEL_PROMETHEUS_BEARER_TOKEN_FILE`; the token is read for each request to support rotation and is never returned.

The client disables ambient proxies and redirects, requires TLS 1.2 or newer for HTTPS, bounds response bytes and timeouts, and only calls fixed read-only `/api/v1` endpoints. It does not expose configuration, flags, rules, alerts, label enumeration, series enumeration, admin APIs, or arbitrary paths.

`prometheus_metric_snapshot` does not accept raw PromQL. Vardiel generates a bounded instant query from a validated metric name, up to eight exact-match diagnostic labels, one of `none`, `sum`, `avg`, `min`, `max`, or `count`, up to five grouping labels, and a hard series limit. Query parameters are submitted in a POST form rather than the URL.

Prometheus output is projected before it reaches the Agent. Scrape URLs, discovered labels, arbitrary target and metric labels, target error text, runtime hostname and working directory, API warning/info text, and raw server errors are omitted. Only warning and info counts are retained.

### Loki diagnostic boundary

Set `VARDIEL_LOKI_URL` to a trusted Loki base URL. HTTPS is required by default. Internal HTTP requires `VARDIEL_LOKI_ALLOW_HTTP=true`. Optional bearer authentication uses an absolute token file, and multi-tenant deployments may set one validated `VARDIEL_LOKI_TENANT_ID`.

The client disables ambient proxies and redirects, requires TLS 1.2 or newer for HTTPS, and only calls `/ready`, `/loki/api/v1/status/buildinfo`, and POST `/loki/api/v1/series`. It does not expose log lines, arbitrary LogQL, query/query-range, tail, labels or label-value enumeration, push, delete, config, metrics, rings, rules, arbitrary paths, or arbitrary methods.

`loki_stream_summary` requires at least one exact-match label from a fixed diagnostic allowlist. Lookback is limited to 1–360 minutes, result count is limited to 1–500, and the selector plus timestamps are sent in a POST form rather than the URL. Returned streams are projected, deduplicated, deterministically sorted, and locally truncated. Unknown labels, file paths, and log content are omitted.

Private, loopback, link-local, multicast, and unspecified HTTP/TLS targets are blocked by the generic network tools by default. Prometheus and Loki use explicitly configured trusted endpoints and do not accept URLs from tool arguments.

## Community

- Read [CONTRIBUTING.md](CONTRIBUTING.md) before proposing a non-trivial change.
- Follow [SECURITY.md](SECURITY.md) for private vulnerability reports; do not disclose security details in a public issue.
- Participation is governed by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).
- Coding agents must follow [AGENTS.md](AGENTS.md) and the assigned issue.
- Public bugs and feature requests should use the repository issue forms.

## Roadmap

The authoritative development sequence is maintained in [docs/development-roadmap.md](docs/development-roadmap.md), [docs/v0.1-scope.md](docs/v0.1-scope.md), and [ADR 0002](docs/adr/0002-fast-autonomous-linux-recovery.md). Tracking issue #22 must be rewritten after the direction reset in [issue #33](https://github.com/Nesoriel/vardiel/issues/33) is accepted.

The immediate sequence is:

1. accept and synchronize the autonomous Linux recovery direction;
2. finish the narrowly scoped safe public error boundary;
3. ship the `systemd-service-unavailable` daemon-to-validation vertical slice;
4. harden local policy, privilege separation, crash recovery, audit, notification, and performance evidence;
5. add new host recovery playbooks one complete slice at a time;
6. add bounded model reasoning and Astralith fleet integration after local recovery is proven.

The universal Tool Contract v2 and general case-bundle work are deferred until a
working recovery slice proves which contracts are actually needed. See
[docs/architecture.md](docs/architecture.md) for current and target boundaries.
Roadmap entries are plans, not implemented capabilities.
