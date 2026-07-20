# Phase 2 — Engineering Specification (Index)

> Complete engineering specification for DevOS. Each document is self-contained and references Phase 0 (research) and Phase 1 (architecture).

| Doc | Title | Covers |
|-----|-------|--------|
| [2.1](01-api-contracts.md) | API Contracts | REST, WebSocket, SSE, SDK, webhook, error/rate-limit model |
| [2.2](02-data-models.md) | Data Models & DB Schemas | PostgreSQL schema, vector store, object store, cache, CQRS |
| [2.3](03-agent-protocol.md) | Agent Communication Protocol | Plugin contract, lifecycle, topic taxonomy, artifacts, discovery |
| [2.4](04-message-bus.md) | Message Bus & Event Catalog | NATS streams, subjects, envelope, full event catalog, replay |
| [2.5](05-channel-adapters.md) | Channel Adapter Contracts | 9 channels, ACK windows, Discord reference, voice, notifications |

**Cross-reference:** ADRs from Phase 1 (001–008) are implemented here as concrete contracts.
**Next:** Phase 3 (UI/UX), Phase 4 (Database — expands 2.2), Phase 5 (Backend), Phase 6 (Frontend), Phase 7 (AI Runtime), Phase 8 (Testing), Phase 9 (Deployment).
