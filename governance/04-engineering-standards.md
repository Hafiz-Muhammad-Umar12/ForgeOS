# DevOS — Engineering Standards

> **Status:** RELEASE CANDIDATE — For Ratification (no production code until ratified)
> **Version:** 1.0-rc1
> **Owner:** Head of Engineering
> **Companion:** Constitution (`01-product-constitution.md`), ADR (`03-adr.md`), Coding (`05-coding-standards.md`), Repo (`06-repository-standards.md`)

---

## 1. Purpose

Defines the **mandatory engineering practices** every contributor follows. These standards implement the Constitutional tenets (especially T6 Observability, T7 Event-Driven, T10 Security) and the Accepted ADRs. Violations block merge.

---

## 2. Architectural Compliance

1. **Ports & Adapters only.** Domain logic depends on interfaces; infrastructure (providers, channels, storage) are adapters (ADR-003, T2/T8).
2. **No cross-context DB access.** Services communicate via the bus or defined ports (ADR-001, T7).
3. **Dependency Injection.** Nothing `new`s its infrastructure dependencies; wiring is at the composition root.
4. **Async/streaming first.** Every long-running operation is async; agent output streams over the bus/SSE/WS.
5. **Bounded contexts** per DDD: Intent, Planning, Orchestration, Agent, Workspace, Delivery, Notification, Observation.

## 3. Testing Requirements

1. **Pyramid:** unit > integration > e2e. Unit/integration on every PR (see Coding §6, DoD §10).
2. **Agent behavior tests** use property-based assertions + golden tasks; real-provider eval is sampled/nightly (non-blocking) (Spec `/specs/08-testing/`).
3. **Coverage floor:** ≥ 80% on non-LLM logic; new code ≥ 80% branch coverage.
4. **No flaky tests in gate.** Flaky tests quarantined, not merged green.

## 4. Code Review

1. Every PR requires **≥ 1 approving review** from a CODEOWNER of the touched area.
2. Reviews verify: spec/ADR conformance, security, observability, tests, no hardcoded config/secrets.
3. Reviews are timely: target first response < 24h.
4. **DCO / sign-off** required (Git Workflow §8).

## 5. Continuous Integration

1. On every PR: lint → format → unit/integration → build images → agent-behavior(scripted) → scan (SBOM + vuln).
2. Required checks must be green to merge.
3. Main branch is always releasable (trunk-based, Git Workflow §8).

## 6. Observability (T6)

1. Every service emits **OTel traces, metrics, structured logs** from first commit.
2. **Agent spans** per run + per tool call; traces must be replayable for any intent.
3. No feature ships without a dashboard panel and an alert on its SLO.
4. Logs are structured JSON; **never log secrets or full PII** (Coding §5).

## 7. Security (T10)

1. Least-privilege by default; OIDC short-lived credentials.
2. Secrets resolved at proxy egress only (T4); never in env/context/agent prompt.
3. Dependency allowlist; all deps scanned.
4. Every mutation audited (`audit_log`).
5. Security review required for auth, secrets, and network-egress changes.

## 8. Performance & Reliability Budgets

1. Channel ACK ≤ provider window (3s Discord / 20s WhatsApp).
2. Workspace cold-start p95 < 5s.
3. Token-stream render p95 < 200ms.
4. Intent success (deploy or safe-fail) ≥ 95%.
5. Regressions on these budgets block merge.

## 9. Feature Management

1. Large/risky features behind **feature flags**, default off.
2. Flags have an owner and a removal plan.

## 10. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |
| CTO | __________ | ☐ Yes ☐ No | ______ |

---

## 11. DevOS Kernel (Mandatory Runtime Authority)

The **DevOS Kernel** is the microkernel (`core/`: message bus, agent/provider registry, dependency-injection container, CQRS, domain aggregates). It is the **only** runtime component permitted to:

1. **Schedule tasks** — dispatch agent runs and drive DAG execution (ADR-002).
2. **Manage agents** — lifecycle, registration, discovery, and assignment.
3. **Manage workspaces** — provisioning and lifecycle (ADR-004).
4. **Manage providers** — registration and routing configuration (ADR-003).
5. **Manage plugins** — load/unload and lifecycle (agents, providers, channels, tools).
6. **Manage runtime lifecycle** — start/stop, scaling coordination, health.

All other services (orchestration, agent-runtime, provider-gateway, workspace-mgr, notification, query) are **tenants of the Kernel, not authorities**. They *request*; the Kernel *authorizes and executes*. Adapters and plugins never self-schedule or self-manage their lifecycle; they expose capabilities to the Kernel.

**Rationale:** A single authority prevents fragmented scheduling/lifecycle logic, enforces Constitutional tenets centrally (Budget T5, HITL T3, Isolation T4, Transparency T11), and keeps adapters swappable. This is the practical realization of the Microkernel + Plugins architecture principle.

**Consequence:** Any code that schedules work, mutates agent/workspace/provider/plugin state, or controls runtime lifecycle outside the Kernel is a violation of this standard and is blocked at review.

---

*End of Engineering Standards v1.0-rc1.*
