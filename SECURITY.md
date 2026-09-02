# Security Policy

Vardiel interacts with operational systems and treats security boundaries as part of its public contract. Please report suspected vulnerabilities privately so they can be investigated and fixed before technical details become public.

## Supported versions

Vardiel has not published a stable release yet.

| Version or branch | Security support |
| --- | --- |
| `main` | Receives security fixes during early development; not a stability guarantee |
| Latest future tagged release | Will be supported after the first release is published |
| Historical commits, archived baselines, and abandoned branches | Not supported |

This table will be updated when release and backport policies exist.

## What to report

Examples include:

- arbitrary shell or code execution;
- an SSRF, DNS rebinding, redirect, or unsafe target-validation bypass;
- exposure of credentials, authorization headers, environment values, private paths, raw infrastructure errors, or sensitive tool output;
- a privilege escalation or unauthorized mutating operation;
- a bypass of standing policy, target allowlists, preconditions, action locks, attempt/concurrency budgets, cooldowns, circuit breakers, validation, or rollback handling;
- action replay, duplicate recovery, cross-user or cross-host authorization confusion, or a model/provider response that can widen action authority;
- path traversal, symlink escape, unsafe file permissions, or case-bundle disclosure;
- a tool contract or protocol path that is falsely advertised as read-only;
- a dependency or build-pipeline compromise with a credible impact on Vardiel users.

A normal configuration question, unsupported environment, or public-data correctness bug may use the bug report form instead.

## Autonomous-recovery security model

Current released behavior is read-only. ADR 0002 permits later focused issues to
add bounded local recovery, but it does not authorize arbitrary commands or make
current CLI/MCP tools mutating.

A future automatic action is authorized only by explicit administrator standing
policy for the exact host or group, target, action, arguments, time window, and
budgets. Model output, events, remote responses, and user prompts cannot create
or widen that grant. The privileged executor must recheck authorization and
preconditions under a per-target lock, bound attempts and concurrency, validate
health independently, record a sanitized audit event, and stop on an open
circuit breaker. Irreversible or broad operations remain human-gated.

Successful recovery records should not expose raw Journal or application logs,
credentials, environment values, private paths, complete command lines, or raw
provider/systemd errors. A failure to notify a remote service must not cause the
local daemon to repeat a mutating action.

## Private reporting

Preferred method:

1. Open the repository **Security** tab.
2. Choose **Report a vulnerability** when GitHub private vulnerability reporting is available.
3. Submit the report without creating a public issue.

Private reporting entry point:

https://github.com/Nesoriel/vardiel/security/advisories/new

If that entry point is unavailable, contact the repository maintainer through a private method listed on their GitHub profile. When no private method is available, open a minimal public issue that asks for a private contact channel. Do not include vulnerability details, proof-of-concept code, secrets, affected hostnames, or private logs in that issue.

## Report contents

Please include, when safely possible:

- affected commit, branch, or release;
- affected component and deployment mode;
- impact and realistic attacker prerequisites;
- reproducible steps or a sanitized proof of concept;
- whether credentials or production data were accessed;
- suggested mitigation or fix;
- your preferred disclosure timeline and attribution.

Use synthetic credentials and infrastructure names. Remove unrelated sensitive data from logs and screenshots.

## Response and disclosure

The maintainers aim to:

- acknowledge a report within 7 calendar days;
- provide an initial severity and scope assessment within 14 calendar days;
- send an update at least every 14 days while a confirmed issue remains unresolved.

These are targets, not contractual service levels. Complex or incomplete reports may take longer.

For a confirmed vulnerability, maintainers will coordinate a fix, tests, release or advisory, and disclosure date with the reporter. The default goal is coordinated disclosure within 90 days, but active exploitation, user safety, upstream coordination, or patch availability may require a shorter or longer timeline.

Do not publicly disclose an unresolved vulnerability without coordinating with maintainers. After remediation, the project may publish a GitHub Security Advisory that credits the reporter if requested.

## Good-faith research

Good-faith research should avoid:

- accessing, modifying, or retaining other people's data;
- disrupting services or exhausting resources;
- persistence, lateral movement, social engineering, or denial of service;
- testing systems you do not own or have permission to assess.

Stop testing and report immediately if sensitive data or unsafe control of a real system becomes possible.
