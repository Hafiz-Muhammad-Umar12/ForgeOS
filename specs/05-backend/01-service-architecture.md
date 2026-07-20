# Phase 5 — Backend Service Architecture (Specification)

> **Status:** Draft
> **Depends on:** Phase 1 (Architecture), Phase 2 (API/Bus/Agent/Channel specs)
> **Scope:** All backend services, their responsibilities, inter-service contracts, and operational characteristics.

---

## 1. Purpose & Responsibilities

The backend is a set of **loosely-coupled, event-driven microservices** plus a microkernel (bus + registry + DI). Each service owns one bounded context and communicates only via the Message Bus or well-defined ports. Responsibilities span intent ingestion → orchestration → agent execution → workspace management → delivery → notification.

---

## 2. Service Map

```mermaid
flowchart TB
    subgraph EDGE["Edge"]
        CH[Channel Adapters]
        IG[Intent Ingress]
        GW[API Gateway]
    end
    subgraph CORE["Control Plane"]
        ORCH[Orchestration Svc]
        QRY[Query/Read Svc]
        NOTIF[Notification Svc]
        REG[Registry Svc]
    end
    subgraph RUNTIME["Runtime Plane"]
        AGRT[Agent Runtime Pool]
        PROV[Provider Gateway]
        WSMGR[Workspace Manager]
    end
    subgraph DATA["Data Plane"]
        PG[(PG)] VDB[(Vector)] OBJ[(Object)] RD[(Redis)] BUS[(NATS)]
    end

    CH-->IG-->BUS
    GW-->ORCH
    GW-->QRY
    ORCH-->BUS
    ORCH-->AGRT
    AGRT-->PROV
    AGRT-->WSMGR
    NOTIF-->BUS
    ORCH-->PG
    QRY-->PG
    AGRT-->VDB
    AGRT-->OBJ
    AGRT-->RD
    BUS-->NOTIF
```

---

## 3. Service Responsibilities

### 3.1 Intent Ingress (`apps/ingress`)
- **Responsibility:** Receive raw channel input, ACK within deadline, parse → `IntentCreated`, emit to bus.
- **Tech:** Go, single binary, horizontally scalable.
- **SLA:** ACK < channel deadline (3s Discord / 20s WhatsApp).
- **Dependencies:** Channel adapters (in-process or sidecar), NATS.

### 3.2 API Gateway (`apps/gateway`)
- **Responsibility:** AuthN/Z (JWT/API key/OIDC), rate-limit, routing, CQRS split (command→Orchestration, query→Query Svc), WebSocket/SSE termination.
- **Tech:** Envoy or Kong (or Go if custom). WebSockets proxied to Orchestration/Query.
- **Dependencies:** Auth provider, Redis (rate/budget), NATS.

### 3.3 Orchestration Service (`services/orchestration`)
- **Sub-modules:** Planner, DAG Executor, Agent Team Coordinator, Scheduler, Budget Governor, HITL Gate.
- **Responsibility:** Turn intent into an approved plan, then coordinate agent execution via the bus.
- **Tech:** Rust (perf-critical scheduler) or Go.
- **State:** Stateless; state in PG + bus. DAG execution state in PG `plans`/`tasks`.
- **Dependencies:** Registry, PG, NATS, (via bus) Agent Runtime, Notification.

### 3.4 Agent Runtime (`services/agent-runtime`)
- **Responsibility:** Load agent plugins, run agent loops, invoke tools, stream tokens to bus, manage memory.
- **Tech:** Python (rich LLM ecosystem). Pool of workers, HPA by queue depth.
- **Isolation:** One process per agent run (sidecar) to contain crashes.
- **Dependencies:** Provider Gateway, Workspace Manager, Vector store, NATS, Bus.

### 3.5 Provider Gateway (`services/provider-gateway`)
- **Responsibility:** Abstract LLM/Tool/Deploy providers behind ports; model routing, circuit breaking, fallback, cost tracking.
- **Tech:** Go. Adapter plugins for claude/codex/gemini/aider/openrouter/ollama + vercel/fly/aws/railway.
- **Dependencies:** External providers (HTTP), Redis (circuit state), NATS.

### 3.6 Workspace Manager (`services/workspace-mgr`)
- **Responsibility:** Provision/warm/recycle isolated workspaces (pods/microVMs); expose gRPC/WS to Agent Runtime; attach Secret Proxy.
- **Tech:** Go + Kubernetes operator / Firecracker.
- **Dependencies:** K8s API, Secret Vault (OIDC), object store (snapshots).

### 3.7 Notification Service (`services/notification`)
- **Responsibility:** Consume `task.*`/`deploy.*` events, render via channel adapters, push to originating + linked channels.
- **Tech:** Go. Subscribes to bus; calls Channel Provider `send()`.
- **Dependencies:** NATS, Channel adapters, Redis (session/linkage).

### 3.8 Query/Read Service (`services/query`)
- **Responsibility:** Serve read models (dashboards, task boards, file trees) from materialized views; WebSocket/SSE fan-out.
- **Tech:** Go/Node. Reads PG read models + Redis CRDT docs.
- **Dependencies:** PG, Redis, NATS (for live updates).

### 3.9 Registry Service (`core/registry`)
- **Responsibility:** Agent & provider discovery; capability index; A2A Agent Card exposure.
- **Tech:** Go. Backed by PG + NATS KV.
- **Dependencies:** PG, NATS.

---

## 4. Inter-Service Contracts (Summary)

| Caller → Callee | Mechanism | Payload |
|-----------------|-----------|---------|
| Ingress → Bus | Publish | `IntentCreated` |
| Gateway → Orchestration | gRPC/HTTP | Command (create intent, approve plan) |
| Orchestration → Agent Runtime | Bus `task.assigned` | Task + agent id |
| Agent Runtime → Provider GW | gRPC | `LLMProvider.complete/stream` |
| Agent Runtime → Workspace Mgr | gRPC/WS | `workspace.exec` |
| Any → Notification | Bus event | `task.*`/`deploy.*` |
| Query → Read models | SQL/Redis | Read query |

**Hard rule:** No service calls another service's database. Only via port/bus.

---

## 5. Scalability & HPA

| Service | Scale Trigger | Min | Max |
|---------|---------------|-----|-----|
| Intent Ingress | CPU | 3 | 20 |
| API Gateway | CPU | 3 | 20 |
| Orchestration | CPU | 3 | 15 |
| Agent Runtime | Bus queue depth | 5 | 200 |
| Provider Gateway | CPU | 3 | 30 |
| Workspace Manager | Pending provisions | 2 | 10 |
| Notification | Bus lag | 3 | 20 |
| Query | CPU | 3 | 20 |

---

## 6. Resilience

- **Circuit breakers** on all external provider calls (Provider Gateway).
- **Bulkheads:** Agent Runtime pool isolated per org to prevent noisy-neighbor.
- **Dead-letter:** Failed bus events → DLQ, alerted.
- **Graceful shutdown:** Drain bus consumers before termination.
- **Health:** `/healthz` + `/readyz` on every service; NATS connectivity probe.

---

## 7. Tradeoffs & Risks

| Decision | Risk | Mitigation |
|----------|------|------------|
| Microservices + bus | Distributed complexity | Strong observability; clear ownership |
| Python Agent Runtime | Perf vs ecosystem | Pool + sidecars; Rust for hot paths |
| Per-run process isolation | Overhead | Pre-warmed workers; lightweight sidecars |
| Stateless services | State in PG/bus | Careful transaction boundaries |

---

## 8. Future Extensions

- **Service mesh** (Istio/Linkerd) for mTLS + traffic policy.
- **WebAssembly agent plugins** for near-native speed + isolation.
- **Multi-region active-active** with bus super-cluster.

---

*End of Phase 5.1 — Backend Service Architecture. Deep dives: 5.2 Provider Gateway, 5.3 Workspace Manager.*
