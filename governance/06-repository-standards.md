# DevOS — Repository Standards

> **Status:** DRAFT — For Approval (no production code until ratified)
> **Version:** 1.0-draft
> **Owner:** Head of Engineering
> **Companion:** Git Workflow (`08-git-workflow.md`), Coding (`05-coding-standards.md`), Engineering (`04-engineering-standards.md`)

---

## 1. Repository Model

DevOS is a **single monorepo** (`ai-native-os`). Rationale: shared `packages/`, atomic cross-cutting changes, unified CI. (Future: split publishable packages if scale demands — decide via ADR.)

## 2. Directory Layout (authoritative)

```
ai-native-os/
├── apps/            # deployable binaries
│   ├── desktop/     # Electron + React
│   ├── web/         # Next.js PWA
│   ├── mobile/      # React Native
│   ├── ingress/     # Intent Ingress service
│   └── gateway/     # API Gateway config
├── services/        # backend microservices
│   ├── orchestration/
│   ├── agent-runtime/
│   ├── workspace-mgr/
│   ├── provider-gateway/
│   ├── notification/
│   └── query/
├── plugins/         # all swappable implementations
│   ├── agents/      # one folder per agent
│   ├── providers/llm|deploy/
│   ├── channels/    # whatsapp|telegram|discord|slack|voice|rest
│   └── tools/
├── core/            # microkernel (bus, registry, di, cqrs, domain)
├── packages/        # shared libs
│   ├── contracts/   # port/type definitions (single source)
│   ├── ui-kit/
│   ├── crdt-sync/
│   └── sdk/
├── infra/           # k8s | terraform | helm
├── specs/           # engineering specs (Phases 0–9)
├── product/         # PRD + product docs
├── governance/      # this folder
├── tests/           # unit | integration | e2e | load | chaos
└── README.md
```

## 3. Dependency Direction (no cycles)

```
core ← services ← apps
core ← plugins ← services
packages/* ← (apps | services | plugins)
```
Rules:
1. `core` depends on nothing internal except `packages/contracts`.
2. `services` depend on `core` + `packages`, never on another `service`'s internals.
3. `plugins` implement `packages/contracts`; never imported by `core`.
4. `apps` depend on `packages` + `services` clients, never on `plugins` internals directly.
5. **No circular dependencies** between any two of these layers. Enforced by build graph + CI check.

## 4. Module Ownership (CODEOWNERS)

- Every top-level dir has a CODEOWNERS entry.
- PRs touching a dir require review from that dir's owner.
- `core/`, `governance/`, `specs/` require CTO or designated architect review.

## 5. Naming & Conventions

- Folders: `kebab-case`. Packages: scoped `@devos/*`.
- Service names match dir: `orchestration`, `agent-runtime`, etc.
- Protobuf/JSON schemas versioned (`schemaVersion`); contract changes via `packages/contracts` PR.

## 6. Dependency Management

- Lockfiles committed (Go `go.sum`, Python `uv.lock`/`poetry.lock`, Node `package-lock.json`).
- Dependency **allowlist**; new deps require review + security scan.
- No unpinned `latest`/`*` ranges in production manifests.
- Generated code (protobuf, mocks) committed or generated in CI; mark `// Code generated — DO NOT EDIT`.

## 7. License & Headers

- Repo licensed under the chosen OSS/commercial license (decide via ADR).
- Source files carry a short SPDX header where required by license.

## 8. Secrets in Repo

- **No secrets in git.** Use secret manager + secret proxy (Constitution T4).
- `.env.example` only; real `.env` git-ignored and scanned by pre-commit.

## 9. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ |

*End of Repository Standards v1.0-draft.*
