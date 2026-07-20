# Phase 2.4 — Message Bus & Event Catalog (Specification)

> **Status:** Draft
> **Depends on:** Phase 1 (ADR-001 NATS JetStream), Phase 0 (pub-sub research)
> **Scope:** Event-driven backbone: transport, streams, subject hierarchy, event envelopes, replay, and the full event catalog.

---

## 1. Purpose & Responsibilities

The Message Bus is the **single integration backbone** of DevOS. All cross-context communication flows through it. Responsibilities:
- Deliver events at-least-once with durable replay.
- Provide subject-based pub-sub with wildcard filtering.
- Carry streaming token deltas at high throughput.
- Support request/reply for synchronous-ish needs (e.g., HITL approval wait).

---

## 2. Transport & Topology

- **Engine:** NATS JetStream (managed, multi-AZ).
- **Persistence:** Streams with `WorkQueue` + `File` storage, replication factor 3.
- **Fan-out:** Core pub-sub streams + KV for config + Object Store for large payloads.
- **Clients:** All services connect via the `core/bus` NATS client wrapper.

---

## 3. Subject Hierarchy

```
devos.{domain}.{entity}.{action}
```
Examples:
```
devos.intent.created
devos.plan.proposed
devos.task.assigned
devos.agent.token        (high-volume stream)
devos.artifact.published
devos.review.passed
devos.deploy.completed
devos.budget.exceeded
```

Wildcards: `devos.task.*` (all task events), `devos.>.completed` (all completions).

---

## 4. Event Envelope

```typescript
interface EventEnvelope<T> {
  id: string;                 // unique event id (dedupe)
  type: string;               // "intent.created"
  schemaVersion: number;      // 1
  traceId: string;            // distributed trace
  orgId: string;
  projectId?: string;
  producedBy: string;         // service/agent id
  producedAt: string;         // ISO8601
  payload: T;
}
```

---

## 5. Stream Definitions

| Stream | Subjects | Retention | Max Age | Purpose |
|--------|----------|-----------|---------|---------|
| `INTENTS` | `devos.intent.>` | WorkQueue | 7d | Orchestration trigger |
| `PLANS` | `devos.plan.>` | File | 30d | HITL + audit |
| `TASKS` | `devos.task.>` | File | 30d | Coordination |
| `TOKENS` | `devos.agent.token` | File | 24h | Live streaming (TTL) |
| `ARTIFACTS` | `devos.artifact.>` | File | 90d | Replay/recall |
| `DEPLOYS` | `devos.deploy.>` | File | 90d | Delivery audit |
| `OBS` | `devos.obs.>` | File | 14d | Metrics/alerts |

---

## 6. Full Event Catalog

| Event | Producer | Consumers | Key Payload Fields |
|-------|----------|-----------|--------------------|
| `intent.created` | Ingress | Orchestrator | channel, text, projectId |
| `intent.cancelled` | Gateway/User | Orchestrator | reason |
| `plan.proposed` | Planner | HITL, User, Bus | dag |
| `plan.approved` | HITL Gate | Orchestrator | approvedBy |
| `plan.rejected` | HITL Gate | Orchestrator | feedback |
| `task.assigned` | Coordinator | Agent Runtime | agentId, taskId |
| `task.status` | Agent Runtime | Coordinator, UI | status |
| `agent.token` | Agent Runtime | UI (stream) | delta |
| `artifact.published` | Any agent | Subscribers by topic | kind, path, payload |
| `review.passed` | Reviewer | QA, Deploy | score |
| `review.comment` | Reviewer | Author agent (loop) | comments |
| `test.result` | QA | Reviewer, Deploy | passed, coverage |
| `security.finding` | Security | Architect, Backend | severity, cve |
| `git.commit` | Git agent | Deploy, Monitor | sha, message |
| `deploy.completed` | Deployment | Notification, Monitor | url, status |
| `task.failed` | Any | Coordinator, Notification | error |
| `budget.exceeded` | Governor | Coordinator, User | ceiling |
| `workspace.ready` | Workspace Mgr | Coordinator | wsId, endpoint |
| `agent.registered` | Registry | Coordinator | capabilities |

---

## 7. Request/Reply Pattern (HITL)

Coordinator publishes `plan.proposed`, then issues a **NATS request** with 15-min timeout:
```
request:  devos.plan.approve.request  (payload: planId)
reply:    devos.plan.approve.reply    (payload: {decision:"approve|reject"})
```
If timeout → auto-`reject` with `timeout` reason (safe default).

---

## 8. Replay & Audit

- Every event persisted in its stream → full replay for debugging/ML.
- `audit_log` table populated by an `OBS` stream consumer (CQRS projector).
- Replay CLI: `devos replay --intent intent_8fa2 --from task.assigned`.

---

## 9. Tradeoffs & Risks

| Decision | Risk | Mitigation |
|----------|------|------------|
| NATS JetStream | Ops expertise needed | Managed offering + runbooks |
| At-least-once | Duplicate handling | Idempotent consumers via `event.id` dedupe |
| High-volume tokens | Throughput | Dedicated `TOKENS` stream, TTL 24h |
| Single bus | Bus outage = halt | Multi-AZ, health probes, circuit breakers |

---

## 10. Future Extensions

- Multi-region super-cluster with leaf nodes.
- Schema registry (Protobuf) for payload validation at the bus.
- Dead-letter stream for poison events.

---

*End of Phase 2.4 — Message Bus & Event Catalog.*
