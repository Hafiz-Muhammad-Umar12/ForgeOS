# DevOS — Documentation Standards

> **Status:** RELEASE CANDIDATE — For Ratification (no production code until ratified)
> **Version:** 1.0-rc1
> **Owner:** Head of Product / Tech Writing
> **Companion:** Engineering Standards (`04-engineering-standards.md`), Specs (`/specs/`)

---

## 1. Principle: Docs as Code

Documentation is **versioned, reviewed, and owned** like code. It lives in-repo, is updated in the same PR as the change it describes, and is never an afterthought.

## 2. Required Document Types

| Type | Owner | Location | When |
|------|-------|----------|------|
| PRD | CTO/Product | `/product/` | Before any build of a feature area |
| ADR | CTO | `/governance/03-adr.md` | On any significant architectural decision |
| Engineering Spec | Architect | `/specs/` | Before implementation of a subsystem |
| README | Module owner | each top-level dir | On creation; kept current |
| Runbook | On-call/Eng | `/infra/runbooks/` | Before a service reaches production |
| API Docs | Backend | generated from code/schemas | On API change |
| Changelog/Release Notes | Release owner | `/product/` or repo root | On each release |

## 3. Spec & Doc Structure (canonical format)

Every engineering/spec document follows this skeleton for consistency:
```
Purpose → Responsibilities → Architecture → Data Flow →
Dependencies → Interfaces → Diagrams (Mermaid) →
Tradeoffs → Risks → Alternatives → Future Extensions
```
This matches the existing `/specs/` set; new docs conform.

## 4. Diagrams as Code

- All diagrams in **Mermaid** (renderable in GitHub/IDEs), committed as text.
- No binary image diagrams (PNG/excalidraw exports) in specs.
- C4 levels used consistently: Context (L1) → Container (L2) → Component (L3).

## 5. Writing Style

1. **Purpose-first:** lead with what and why.
2. **Concrete over vague:** name services, cite ADR/PRD IDs, show schemas.
3. **Audience-aware:** specs for engineers; PRD for product/leadership; runbooks for on-call.
4. **No duplication drift:** single source of truth; cross-reference, don't copy.
5. **Active voice**, present tense.

## 6. Versioning & Lifecycle

- Docs carry a `Status` (Draft / Approved / Deprecated) and `Version`.
- Deprecated docs are marked, not deleted, for traceability.
- Links use relative paths; broken links fail CI doc-check.

## 7. Ownership

- Each doc has a named owner in its header.
- PRs editing a doc require that owner's review (CODEOWNERS).

## 8. Approval

| Role | Name | Accept? | Date |
|------|------|---------|------|
| Head of Product | __________ | ☐ Yes ☐ No | ______ |

*End of Documentation Standards v1.0-rc1.*
