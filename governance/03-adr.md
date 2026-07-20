# DevOS — Architecture Decision Records (ADR)

> **Status:** DRAFT — For Approval (no production code until ratified)
> **Version:** 1.0-draft
> **Owner:** CTO (ADR custodian)
> **Companion:** Product Constitution (`01-product-constitution.md`), Engineering Specs (`/specs/01-architecture/`)
> **Supersedes:** The ADR summary embedded in `/specs/01-architecture/01-system-architecture.md` §7 (that summary is informational; this register is authoritative).

---

## 1. Purpose

This register is the **authoritative record of accepted architectural decisions** for DevOS. It binds engineering: implementations must conform to Accepted ADRs. Conflicts between an ADR and the prose specs resolve in favor of the ADR.

## 2. ADR Lifecycle & Process

| Status | Meaning |
|--------|---------|
| `Proposed` | Drafted, under discussion; not yet binding. |
| `Accepted` | Approved by CTO; binding on implementation. |
| `Deprecated` | No longer recommended; retained for history. |
| `Superseded` | Replaced by a newer ADR (cite the successor). |

**Process:**
1. Propose via PR adding `ADR-0XX-title.md` (or append here pre-ratification) using the template in §4.
2. CTO reviews; may request changes.
3. Acceptance requires **CTO sign-off** (CEO + Head of Product notified). Tenet-touching ADRs (see Constitution Article II) require unanimous ratification.
4. On acceptance: set `Status: Accepted`, record in the Index (§3), and bump the register version.
5. To change a decision: propose a new ADR that `Supersedes` the old one; never silently edit an Accepted ADR's Decision.

## 3. Index

| ID | Title | Status | Summary |
|----|-------|--------|---------|
| ADR-001 | Event-Driven Core with NATS JetStream | Accepted | All cross-context comms via pub-sub bus |
| ADR-002 | Stateful DAG Orchestration | Accepted | Workflows = DAGs in Pregel-style super-steps |
| ADR-003 | Provider Abstraction via Ports | Accepted | LLM/Tool/Deploy/Channel/Vector as interfaces |
| ADR-004 | Workspace Isolation via Containers/MicroVMs | Accepted | Sealed pods + Secret Proxy |
| ADR-005 | CRDT Client Sync (Yjs + Edge Relay) | Accepted | Local-first, offline-capable multi-device |
| ADR-006 | Channel Adapters Funnel to Uniform Intent | Accepted | One `IntentCreated` envelope; backend channel-agnostic |
| ADR-007 | Human-in-the-Loop as First-Class Gate | Accepted | Approve/Reject are DAG nodes, not exceptions |
| ADR-008 | Budget Governor | Accepted | Per-org/project/user token & cost ceilings |

## 4. Template

```markdown
# ADR-0XX — <Title>
- Status: Proposed | Accepted | Deprecated | Superseded by ADR-0YY
- Date: YYYY-MM-DD
- Deciders: <roles>
- Related: <Constitution tenets, PRD sections, specs>

## Context
What problem/forces motivate this decision?

## Decision
What we decided, stated as a rule.

## Consequences
Positive, negative, and tradeoffs. What becomes easier/harder.

## Alternatives Considered
What we rejected and why.
```

---

## 5. Accepted ADRs

### ADR-001 — Event-Driven Core with NATS JetStream
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T7; Spec `/specs/01-architecture/` §7; `/specs/02-specification/04-message-bus.md`
- **Context:** Cross-context coupling via direct calls does not scale and prevents replay/audit. Research (MetaGPT pub-sub, AutoGen v0.4 event-driven, CrewAI Flows) favors a bus.
- **Decision:** All cross-context communication uses **NATS JetStream** (pub-sub + durable streams + KV + object store). No service calls another's datastore.
- **Consequences:** +Scalable, replayable, auditable. −Operational complexity (mitigated by managed NATS).
- **Alternatives:** Kafka (heavier ops), Redis Streams (weaker replay), RabbitMQ (no native streaming).

### ADR-002 — Stateful DAG Orchestration
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T7; Spec `/specs/01-architecture/` §9.2
- **Context:** Linear pipelines are too rigid; pure LLM-agent negotiation is non-deterministic and unauditable. ChatDev v2 (DAG YAML) and LangGraph (super-steps) validate graph execution.
- **Decision:** Workflows are **DAGs executed in super-steps**; nodes may be agents, tools, humans, or subgraphs. Supports map-reduce, conditional edges, human nodes.
- **Consequences:** +Deterministic, debuggable, composable. −Graph engine complexity (mitigated by subgraph isolation).
- **Alternatives:** Linear pipeline (rigid), free-form agent negotiation (unauditable).

### ADR-003 — Provider Abstraction via Ports
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T2, T8; Spec `/specs/02-specification/03-agent-protocol.md`, `/specs/05-backend/02-provider-gateway.md`
- **Context:** Vendor lock-in is an existential risk; every major framework ships a pluggable backend.
- **Decision:** `LLMProvider`, `ToolProvider`, `DeployProvider`, `ChannelProvider`, `VectorProvider` are **interfaces**; adapters are plugins. Capability flags express feature differences.
- **Consequences:** +No lock-in, +swappable. −Lowest-common-denominator risk (solved by capability flags + passthrough).
- **Alternatives:** Hardcode providers (rejected — violates T2).

### ADR-004 — Workspace Isolation via Containers/MicroVMs
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T4; Spec `/specs/05-backend/03-workspace-manager.md`
- **Context:** Agent code execution must be safe. Docker sandboxing is consensus (SWE-agent, OpenHands, ChatDev).
- **Decision:** Each workspace = **sealed pod** (FS, git, pkg-mgr, CLI, browser, DB, secret-proxy). Pre-warmed pool for <5s cold start. Firecracker microVMs as v2 security upgrade.
- **Consequences:** +Safe execution, +reproducible. −Cold-start latency (mitigated by warm pool).
- **Alternatives:** Full VMs (slow), shared FS (insecure — rejected).

### ADR-005 — CRDT Client Sync (Yjs + Edge Relay)
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T9; Spec `/specs/01-architecture/` §7.5, `/specs/06-frontend/01-frontend-architecture.md`
- **Context:** Multi-device, offline-capable, real-time collaboration requires local-first sync. Yjs/Liveblocks/PartyKit prove CRDT viability.
- **Decision:** Client shared state syncs via **Yjs CRDTs over an edge WebSocket relay**. Servers are supporting infra, not source of truth.
- **Consequences:** +Offline-first, +realtime, +multi-device parity. −CRDT memory overhead (mitigated by bounded docs + server compaction).
- **Alternatives:** Centralized polling (no offline), OT (complex).

### ADR-006 — Channel Adapters Funnel to Uniform Intent
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T1; Spec `/specs/02-specification/05-channel-adapters.md`
- **Context:** No existing platform unifies all nine surfaces; this is our differentiator.
- **Decision:** Every channel emits the same `IntentCreated` event with channel metadata; backend is channel-agnostic. ACK windows (3s Discord / 20s WhatsApp) honored at the adapter.
- **Consequences:** +True multi-surface, +backend simplicity. −Channel-specific nuance lost (mitigated by `ChannelContext` passthrough, not acted on by core).
- **Alternatives:** Per-channel backend logic (rejected — violates T1).

### ADR-007 — Human-in-the-Loop as First-Class Gate
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T3; Spec `/specs/01-architecture/` §7.7
- **Context:** Autonomous agents need safety checkpoints; treating HITL as an exception leads to unsafe flows.
- **Decision:** Approval checkpoints (plan review, deploy, secret access) are **explicit DAG nodes**, not exceptions.
- **Consequences:** +Safe autonomy, +auditable. −Some friction (mitigated by configurable strictness per project).
- **Alternatives:** Opt-in HITL (rejected — unsafe default).

### ADR-008 — Budget Governor
- **Status:** Accepted · **Date:** 2026-07-20 · **Deciders:** CTO
- **Related:** Constitution T5; Spec `/specs/01-architecture/` §7.8, `/specs/07-ai-runtime/02-context-memory.md`
- **Context:** Multi-agent costs ~3× single-agent (ChatDev 22,949 vs 7,183 tokens). Unbounded spend is existential.
- **Decision:** Per-project/per-user **token & cost ceilings** enforced by a governor before dispatch; pauses + notifies on exceed.
- **Consequences:** +Cost control, +predictable margins. −Possible premature abort (mitigated by tunable ceiling + user override).
- **Alternatives:** Post-hoc billing (rejected — no spend control).

---

## 6. Version & Approval

| Version | Date | Change | Approved By |
|---------|------|--------|-------------|
| 1.0-draft | 2026-07-20 | Initial 8 ADRs ratified | _pending CTO_ |

**Sign-off:**
| Role | Name | Accept? | Date |
|------|------|---------|------|
| CTO | __________ | ☐ Yes ☐ No | ______ |

*End of ADR Register v1.0-draft.*
