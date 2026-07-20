# Phase 4 — Database & Data Architecture (Specification)

> **Status:** Draft
> **Depends on:** Phase 2.2 (Schemas), Phase 1 (CQRS ADR, Microkernel)
> **Scope:** Migration strategy, data-access patterns (repositories/ports), indexing, read-model projections, sharding, and operational runbooks. *Schemas themselves live in Phase 2.2.*

---

## 1. Purpose & Responsibilities

Beyond schema definition, this document specifies **how data moves**: migrations without downtime, clean repository boundaries (hexagonal), materialized read models for CQRS, and scale-out strategy for millions of developers.

---

## 2. Storage Topology (Recap)

| Store | Engine | Role |
|-------|--------|------|
| Primary | PostgreSQL 16+ | Write aggregates (intents, plans, tasks, deployments) |
| Vector | pgvector or Qdrant | Embeddings (code, memory, intent history) |
| Object | S3 / R2 | Artifacts, builds, snapshots, git bundles |
| Cache/CRDT | Redis 7 | Rate/budget, Yjs docs, sessions |
| Stream | NATS JetStream | Event log (see Phase 2.4) |

---

## 3. Migration Strategy

### 3.1 Tooling
- **Migrator:** `golang-migrate` (CLI + lib), versioned, checksum-verified.
- **Format:** `{version}_{name}.up.sql` / `.down.sql`.
- **Applied tracking:** `schema_migrations` table.

### 3.2 Zero-Downtime Rules
1. **Additive first:** New columns `DEFAULT NULL` or with safe defaults; backfill async.
2. **No `ALTER COLUMN TYPE` in place** on large tables → add new col, dual-write, swap, drop.
3. **Indexes:** `CREATE INDEX CONCURRENTLY` (never block writes).
4. **Drops:** Only after 1+ release with dual-write, gated by feature flag.
5. **Expand/Contract:** Every breaking change goes expand (add) → migrate → contract (remove).

### 3.3 Example: adding `auto_approve` to intents
```sql
-- up
ALTER TABLE intents ADD COLUMN auto_approve BOOLEAN NOT NULL DEFAULT false;
-- (backfill where appropriate)
-- down
ALTER TABLE intents DROP COLUMN auto_approve;
```

---

## 4. Data Access Patterns (Hexagonal)

Domain aggregates never touch SQL directly. They go through **repository ports**:

```typescript
interface IntentRepository {
  save(intent: Intent): Promise<void>;
  findById(id: string): Promise<Intent | null>;
  listByOrg(orgId: string, page: Page): Promise<Intent[]>;
}
interface TaskRepository {
  save(task: Task): Promise<void>;
  findByPlan(planId: string): Promise<Task[]>;
  updateStatus(id: string, status: TaskStatus): Promise<void>;
}
```

- **Implementation:** `PostgresIntentRepository` behind the port.
- **DI:** Injected at runtime (ADR: Dependency Injection). Tests use in-memory fakes.
- **Transactions:** Aggregate-scoped; cross-aggregate coordination via the bus, not 2PC.

---

## 5. CQRS Read Models (Projections)

Command side writes aggregates → bus events → **projectors** build read models.

| Read Model | Source Events | Table / Store |
|------------|---------------|---------------|
| Project dashboard | intent/task/deploy events | `project_view` (Redis/PG) |
| Task board | task.status | `task_board_view` |
| Agent activity | agent.token/status | `agent_activity_view` |
| Billing/budget | token_usage | `org_usage_rollup` |
| Search (NL history) | intent.created | vector `intent_history` |

**Freshness:** Eventual, target < 500ms. Critical paths (HITL approval) read write-model directly to avoid stale approvals.

---

## 6. Indexing Strategy

```sql
-- Hot-path indexes (already in 2.2)
CREATE INDEX idx_intents_org_status ON intents(org_id, status);
CREATE INDEX idx_tasks_plan ON tasks(plan_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_audit_org_time ON audit_log(org_id, created_at DESC);

-- Composite for dashboards
CREATE INDEX idx_tasks_project_status ON tasks(project_id, status);
-- Partial for active only
CREATE INDEX idx_ws_active ON workspaces(org_id) WHERE status IN ('active','warm');
```

- **Token table:** time-partitioned by `created_at` (monthly); TTL + compaction.
- **Audit:** partitioned by month; cold tier to columnar (future).

---

## 7. Sharding & Scale-Out

| Scale Stage | Strategy |
|-------------|----------|
| <10M intents | Single PG instance, read replica |
| 10–100M | Read replicas + `org_id` range partitioning |
| >100M | **Shard by `org_id`** (logical shards in Citus/PostgresML); bus + projectors per shard |
| Vector | Dedicated Qdrant cluster, per-org namespaces |

**Routing:** `org_id` is in every event → deterministic shard routing at the gateway.

---

## 8. Backup & DR

- **PG:** Continuous WAL archiving → object store; PITR ≤ 5 min RPO.
- **Redis:** AOF + snapshot; CRDT docs durable in PG-backed Yjs store.
- **Object:** Versioned buckets; cross-region replication.
- **Bus:** JetStream replicated (RF=3); replay from stream on recover.
- **RTO:** < 15 min (managed failover). **RPO:** < 5 min.

---

## 9. Data Retention & Compliance

| Data | Retention | Note |
|------|-----------|------|
| Agent tokens | 24h (stream) / 90d (archive) | Replay then compact |
| Audit log | 7y (enterprise) / 1y (default) | Compliance |
| Artifacts | Project lifetime + 90d | Snapshot before delete |
| CRDT docs | Project lifetime | Export on archive |

**GDPR:** Right-to-erasure → `org_id` cascade delete job; anonymize audit after 90d if requested.

---

## 10. Tradeoffs & Risks

| Decision | Risk | Mitigation |
|----------|------|------------|
| Eventual read models | Stale reads | Critical paths read write-model; freshness SLA |
| Org sharding | Cross-org queries hard | Analytics via separate columnar warehouse |
| Append-only tokens | Storage growth | Partition + TTL + compaction job |
| Dual-write migrations | Consistency window | Feature-flagged, monitored |

---

## 11. Future Extensions

- **ClickHouse** for audit/observability analytics.
- **Neo4j** for agent-dependency graph visualization.
- **CQRS event-sourcing** for intents (full event log as source of truth, not just audit).

---

*End of Phase 4 — Database & Data Architecture. (Schemas: see Phase 2.2.)*
