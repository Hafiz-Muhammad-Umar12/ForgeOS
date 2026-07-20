# Phase 2.1 — API Contracts (Specification)

> **Status:** Draft
> **Depends on:** Phase 1 (Architecture)
> **Scope:** REST API, WebSocket/SSE streaming, SDK surface, and cross-channel intent contract for DevOS.

---

## 1. Purpose & Responsibilities

The API layer is the **uniform contract** through which every channel (Desktop, Web, Mobile, WhatsApp, Telegram, Discord, Slack, Voice, REST) interacts with DevOS. It must:
- Accept intents in a channel-agnostic envelope.
- Expose project/task/agent/workspace/deployment lifecycle operations.
- Stream all agent output (tokens, artifacts, status) over WebSocket/SSE.
- Be versioned, documented, and backward-compatible.

---

## 2. API Surface Overview

| Surface | Protocol | Purpose |
|---------|----------|---------|
| Public REST | HTTPS/JSON | Commands & queries from external clients |
| WebSocket | `wss://` | Real-time streaming of agent output, task updates |
| SSE | `text/event-stream` | Server-push for channels that can't hold WS (e.g., webhook bots) |
| SDK | TS/Python | Typed client for Desktop/Web/Mobile/REST integrators |
| Webhook | HTTPS | Inbound from messaging channels (Discord/Slack/Telegram/WhatsApp/Voice) |

**Base URL:** `https://api.devos.ai/v1` (self-host: `<ingress>/v1`)

---

## 3. Authentication

| Mechanism | Use | Header |
|-----------|-----|--------|
| Bearer JWT | REST/WS from authenticated apps | `Authorization: Bearer <jwt>` |
| API Key | Service/CI integration | `X-DevOS-Key: <key>` |
| Channel Signature | Inbound webhooks | Per-channel (HMAC / token) |
| OIDC Exchange | Cloud deploy creds | `Authorization: Bearer <oidc-jwt>` |

All mutations require `scope` claim: `intent:write`, `project:read`, `deploy:execute`, etc.

---

## 4. REST Endpoints

### 4.1 Intents
```
POST   /v1/intents                  Create intent (NL → plan)
GET    /v1/intents/:id              Get intent status
GET    /v1/intents/:id/stream       SSE stream of intent events
DELETE /v1/intents/:id              Cancel intent
```

**POST /v1/intents** — Request
```json
{
  "channel": "discord",
  "channelUserId": "u_123",
  "projectId": "proj_ecommerce",
  "text": "Build an ecommerce platform with Stripe payments",
  "attachments": [],
  "priority": "normal",
  "autoApprove": false
}
```
**Response `201`**
```json
{
  "id": "intent_8fa2",
  "status": "planning",
  "createdAt": "2026-07-20T10:00:00Z",
  "planId": "plan_91c3",
  "links": {
    "stream": "/v1/intents/intent_8fa2/stream",
    "plan": "/v1/plans/plan_91c3"
  }
}
```

### 4.2 Projects
```
POST   /v1/projects                Create project
GET    /v1/projects                List (paginated, filter)
GET    /v1/projects/:id            Get project
PATCH  /v1/projects/:id            Update metadata
DELETE /v1/projects/:id            Archive
GET    /v1/projects/:id/files      File tree (CRDT snapshot)
POST   /v1/projects/:id/files      Write file
GET    /v1/projects/:id/git/log    Git history
```

### 4.3 Plans (HITL)
```
GET    /v1/plans/:id               Get plan DAG
POST   /v1/plans/:id/approve       Approve (HITL gate)
POST   /v1/plans/:id/reject         Reject with feedback
```

### 4.4 Tasks
```
GET    /v1/tasks?projectId=&status=   List tasks
GET    /v1/tasks/:id                  Get task + agent log
POST   /v1/tasks/:id/retry            Retry failed task
```

### 4.5 Agents
```
GET    /v1/agents                     List registered agent plugins
GET    /v1/agents/:id/capabilities    Capability flags
POST   /v1/agents/:id/invoke          Direct invoke (advanced)
```

### 4.6 Workspaces
```
GET    /v1/workspaces                 List active workspaces
POST   /v1/workspaces                 Provision (warm)
GET    /v1/workspaces/:id/status      Health
DELETE /v1/workspaces/:id             Recycle
POST   /v1/workspaces/:id/exec        Run CLI command (streamed)
```

### 4.7 Deployments
```
POST   /v1/deployments                Deploy artifact
GET    /v1/deployments/:id            Status
GET    /v1/deployments/:id/logs       Build/deploy logs
POST   /v1/deployments/:id/rollback   Rollback
```

### 4.8 Monitoring
```
GET    /v1/projects/:id/metrics       Health metrics
GET    /v1/projects/:id/alerts        Active alerts
GET    /v1/projects/:id/logs          Aggregated logs
```

---

## 5. WebSocket Streaming Contract

**Connect:** `wss://api.devos.ai/v1/stream?intentId=intent_8fa2`
**Auth:** JWT in `Sec-WebSocket-Protocol` or query.

**Server→Client event frames:**
```json
{ "type": "intent.status",   "intentId": "...", "status": "executing" }
{ "type": "plan.proposed",   "planId": "...", "dag": { "...": "..." } }
{ "type": "task.assigned",   "taskId": "...", "agent": "frontend" }
{ "type": "agent.token",     "taskId": "...", "agent": "frontend", "delta": "function App()" }
{ "type": "artifact.published", "taskId":"...", "kind":"code", "path":"src/App.tsx" }
{ "type": "review.passed",   "taskId": "...", "score": 0.92 }
{ "type": "deploy.completed","url": "https://ecom.devos.app" }
{ "type": "intent.completed","status":"success", "summary":"..." }
{ "type": "error",           "code":"BUDGET_EXCEEDED", "message":"..." }
```

**Client→Server:**
```json
{ "type": "plan.approve",   "planId": "..." }
{ "type": "plan.reject",    "planId": "...", "feedback": "Use PostgreSQL not Mongo" }
{ "type": "task.cancel",    "taskId": "..." }
```

---

## 6. SSE Contract (Channel Webhooks / Lightweight Clients)

`GET /v1/intents/:id/stream`
```
event: intent.status
data: {"status":"executing"}

event: agent.token
data: {"taskId":"t1","agent":"frontend","delta":"export "}

event: deploy.completed
data: {"url":"https://ecom.devos.app"}
```

---

## 7. Webhook Inbound (Channel → DevOS)

Each channel adapter exposes `POST /webhooks/:channel` (e.g., `/webhooks/discord`). DevOS verifies signature, calls `ChannelProvider.parse()`, emits `IntentCreated`.

**Canonical Intent Envelope (internal, post-parse):**
```json
{
  "id": "intent_8fa2",
  "channel": "discord",
  "channelUserId": "u_123",
  "channelGuildId": "g_456",
  "projectId": null,
  "text": "Build an ecommerce platform",
  "locale": "en",
  "attachments": [],
  "createdAt": "2026-07-20T10:00:00Z",
  "traceId": "trace_abc"
}
```

---

## 8. SDK Surface (TypeScript)

```typescript
import { DevOS } from "@devos/sdk";

const client = new DevOS({ token: process.env.DEVOS_TOKEN });

// Create intent and stream
const intent = await client.intents.create({
  text: "Build an ecommerce platform",
  projectId: "proj_ecommerce",
});

for await (const evt of client.intents.stream(intent.id)) {
  if (evt.type === "plan.proposed") {
    await client.plans.approve(evt.planId);   // HITL
  }
  if (evt.type === "agent.token") process.stdout.write(evt.delta);
  if (evt.type === "deploy.completed") console.log("Live:", evt.url);
}
```

---

## 9. Error Model

```json
{
  "error": {
    "code": "PLAN_REJECTED|BUDGET_EXCEEDED|WORKSPACE_UNAVAILABLE|PROVIDER_DOWN|FORBIDDEN",
    "message": "Human-readable",
    "traceId": "trace_abc",
    "retryable": false
  }
}
```
HTTP status mapping: 400 (validation), 401 (auth), 403 (scope), 409 (conflict), 429 (rate/budget), 503 (provider/workspace down).

---

## 10. Rate Limiting & Quotas

| Tier | Intents/min | Concurrent agents | Token budget/month |
|------|-------------|-------------------|--------------------|
| Free | 5 | 2 | 1M |
| Pro | 30 | 8 | 20M |
| Team | 200 | 32 | 200M |
| Enterprise | custom | custom | custom |

Rate state stored in Redis; `429` includes `Retry-After` and `X-Budget-Remaining`.

---

## 11. Versioning

- URL version (`/v1`) + `Accept: application/vnd.devos.v1+json`.
- Breaking changes → new major version; old version supported ≥ 12 months.
- Event schema versioned via `schemaVersion` field in every frame.

---

## 12. Future Extensions

- GraphQL gateway for complex client queries.
- Webhook subscriptions for user-defined events (`POST /v1/subscriptions`).
- Bidirectional voice streaming (`wss://.../voice`) with interim transcripts.

---

*End of Phase 2.1 — API Contracts.*
