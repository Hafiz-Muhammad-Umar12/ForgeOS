# DevOS Kernel (`core/`)

The **DevOS Kernel** is the microkernel (`core/`): the composition root and the
sole runtime authority for the platform (SDD §11, governance
04-engineering-standards.md §11). This directory currently contains the
**Component 2 skeleton** — the foundation other milestones build on.

## Scope of this component

Implemented in this component (production-quality, tested, stdlib-only):

| Package | Responsibility |
|---------|----------------|
| `di` | Reflection-based dependency-injection container (singleton/transient, circular-dependency detection). |
| `lifecycle` | Ordered component startup, graceful reverse-order shutdown, health aggregation, `Run` loop. |
| `config` | Kernel config loading from env (prefix) + `.env`-style files, with validation. |
| `event` | `EventEnvelope[T]` (specs/02-specification/04-message-bus.md §4) and the canonical event catalog + subject helpers. |
| `registry` | **Service-registry port** (`Registry` interface) + `ServiceInfo`/`Capability` types. |
| `scheduler` | **Scheduler / `KernelAuthority` ports** + `Task`/`ScheduleRequest`/`TaskState` types. |
| `domain` | Bounded-context strong identifiers (`OrgID`, `ProjectID`, …) and the `Aggregate` contract. |
| `kernel` | Composition root: loads config, builds the DI container, wires core services, exposes the `Kernel` facade. |
| `cmd/kernel` | Binary entry point (signal-aware; runs until interrupted). |

## Deliberately NOT implemented (deferred per scope)

- **NATS integration** (`core/bus`) — transport is a later milestone (ADR-001).
- **Registry client** (NATS-KV-backed) — only the `Registry` *port* exists (ADR-003, SDD §09).
- **Scheduler implementation** — only the *ports* exist (ADR-002).
- **CQRS** (`core/cqrs`) — deferred.
- Business logic, authentication, providers, agents, and the API Gateway.

These gaps are intentional: this component establishes the contracts and
runtime plumbing so later milestones provide swappable implementations behind
the interfaces defined here.

## Architecture & dependency direction

```
cmd/kernel ──▶ kernel (composition root)
                 ├─▶ config   (loads kernel config)
                 ├─▶ di       (wires services)
                 └─▶ lifecycle (manages components)
scheduler ──▶ registry   (port types only; no cycle)
domain     (shared value types)
```

Dependency rule (governance/06-repository-standards.md §3): `core` depends only
on `packages/*` (none yet) and the Go standard library. Tenant services depend
on `core`; `core` never depends on a tenant service. The `registry` and
`scheduler` packages define interfaces that implementations (delivered later)
satisfy — callers depend on the interface, not the concrete store.

## Build & test

Requires Go 1.23+ (pinned via `mise.toml` at repo root).

```bash
cd core
go build ./...
go test ./...
go vet ./...
gofmt -l .
```

From the repo root (Go workspace, `go.work`):

```bash
go work sync
go build ./...
go test ./...
```

## Bootstrap flow

1. `cmd/kernel` calls `kernel.New()`, which loads `config.Load` and builds the
   DI `Container`, registering `*config.Config` and `*di.Container` as
   singletons.
2. A `lifecycle.Manager` is created (empty in the skeleton — tenant components
   are registered in later milestones).
3. `Kernel.Run(ctx)` starts the manager and blocks until a signal/context
   cancellation, then stops gracefully.

## Conventions

- Every exported identifier has a godoc comment (revive `exported`).
- Tests are co-located (`*_test.go`) and use fakes for ports (no real NATS).
- No `TODO`s or placeholder code; the skeleton is complete for its scope.

See also: `specs/01-architecture/01-system-architecture.md`, `sdd/11-kernel.md`,
`governance/03-adr.md` (ADR-001/002/003), and `planning/03-build-order.md`.
