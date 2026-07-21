# DevOS — System Design Document (SDD)

> **Status:** DRAFT — For Approval (no production code until ratified)
> **Version:** 1.0-draft
> **Owner:** CTO / Architecture
> **Companion:** Governance (`/governance/`), Engineering Specs (`/specs/`), PRD (`/product/PRD.md`)
> **Role:** The bridge between **governance** (what we must obey) and **implementation** (how we build). Every major service is specified here to the level required to begin implementation planning.

---

## 1. How to Read This Document

- Each major service has its own file under `/sdd/`.
- Every service section follows a fixed template (§3) and includes **four Mermaid diagrams**: Component, Sequence, State, Data Flow.
- Services are specified, not implemented. No code exists yet.
- Where a service depends on a governance rule or spec, it is cited (e.g., `Constitution T11`, `ADR-003`, `Spec 2.1`).

## 2. Service Inventory

| # | Service | File | Kind | Primary Responsibility |
|---|---------|------|------|------------------------|
| 1 | Intent Ingress | [01-intent-ingress.md](01-intent-ingress.md) | Edge | Accept channel input, ACK, parse → `intent.created` |
| 2 | API Gateway | [02-api-gateway.md](02-api-gateway.md) | Edge | AuthN/Z, rate-limit, routing, WS/SSE |
| 3 | Orchestration | [03-orchestration.md](03-orchestration.md) | Control | Plan, coordinate agents, enforce HITL + budget |
| 4 | Agent Runtime | [04-agent-runtime.md](04-agent-runtime.md) | Runtime | Run agent loops, tools, memory, streaming |
| 5 | Provider Gateway | [05-provider-gateway.md](05-provider-gateway.md) | Control | LLM/Tool/Deploy abstraction, routing, fallback |
| 6 | Workspace Manager | [06-workspace-manager.md](06-workspace-manager.md) | Runtime | Isolated workspaces, secret proxy, lifecycle |
| 7 | Notification | [07-notification.md](07-notification.md) | Edge | Render + push events to channels |
| 8 | Query (Read) | [08-query-service.md](08-query-service.md) | Control | Read models, WS/SSE fan-out |
| 9 | Registry | [09-registry.md](09-registry.md) | Core | Agent/provider discovery, A2A cards |
| 10 | Channel Adapters | [10-channel-adapters.md](10-channel-adapters.md) | Edge (subsystem) | 9-channel translate to/from canonical |
| 11 | DevOS Kernel | [11-kernel.md](11-kernel.md) | Core (microkernel) | Sole authority: schedule/agent/workspace/provider/plugin/runtime |

## 3. Per-Service Template

```
1. Purpose
2. Responsibilities
3. Architecture            (Mermaid Component Diagram)
4. Interaction Sequence    (Mermaid Sequence Diagram)
5. Interfaces (ports)
6. APIs (endpoints / gRPC / signatures)
7. Events (published / consumed)
8. State Machine          (Mermaid State Diagram)
9. Folder Structure
10. Dependencies
11. Data Flow             (Mermaid Data Flow Diagram)
12. Failure Handling
13. Security
14. Scalability
15. Testing Strategy
16. Future Extensions
```

## 4. Cross-Cutting Rules (from Governance)

- **Constitution T1–T12** bind every service (channel-agnostic, no lock-in, HITL, isolation, budget, observability, event-driven, provider abstraction, offline-first, security, transparency, open standards).
- **ADR-001** event bus = NATS JetStream; no cross-context DB access.
- **ADR-003** providers behind ports.
- **ADR-004** workspace isolation.
- **ADR-005** CRDT client sync.
- **ADR-006** uniform intent envelope.
- **ADR-007** HITL as DAG node.
- **ADR-008** budget governor.
- **Eng §11** Kernel is the only runtime authority for scheduling/lifecycle.

## 5. Approval

| Role | Name | Approve? | Date |
|------|------|---------|------|
| CTO | __________ | ☐ Yes ☐ No | ______ |
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |

*Once approved, implementation planning begins. No production code before then.*
