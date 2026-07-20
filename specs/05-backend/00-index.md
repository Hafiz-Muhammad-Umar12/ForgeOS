# Phase 5 — Backend (Index)

> Microservices + microkernel for the DevOS control/runtime planes.

| Doc | Title | Covers |
|-----|-------|--------|
| [5.1](01-service-architecture.md) | Service Architecture | Service map, responsibilities, inter-service contracts, HPA, resilience |
| [5.2](02-provider-gateway.md) | Provider Gateway | LLM abstraction, model tiers, fallback, cost ledger |
| [5.3](03-workspace-manager.md) | Workspace Manager | Isolation, warm pool, secret proxy, lifecycle |

**Key rule:** No service touches another's DB; all cross-context via NATS bus or ports.
**Next:** Phase 6 (Frontend), Phase 7 (AI Runtime), Phase 8 (Testing), Phase 9 (Deployment).
