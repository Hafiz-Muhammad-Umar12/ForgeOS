# DevOS — Governance Documents

> **Status:** RELEASE CANDIDATE — For Ratification
> **Owner:** CTO
> **Binding Rule:** These documents are the **supreme project governance**, ranking above PRD, specs, and ADRs. **No production code may be written until all ten are approved.**

---

## Foundation Documents (11 + index)

| # | Document | Purpose | Ranks |
|---|----------|---------|-------|
| 1 | [Product Constitution](01-product-constitution.md) | Immutable, inviolable tenets every decision must obey | Supreme |
| 2 | [DevOS Manifesto](02-devos-manifesto.md) | Vision & philosophy — what we believe | Inspirational |
| 3 | [Architecture Decision Records](03-adr.md) | Authoritative accepted architectural decisions (ADR-001…008) | Binding on eng |
| 4 | [Engineering Standards](04-engineering-standards.md) | Mandatory engineering practices (testing, CI, observability, security) | Binding on eng |
| 5 | [Coding Standards](05-coding-standards.md) | Concrete code conventions (Go/Python/TS/SQL) | Binding on eng |
| 6 | [Repository Standards](06-repository-standards.md) | Monorepo layout, dependency direction, ownership | Binding on eng |
| 7 | [Documentation Standards](07-documentation-standards.md) | Docs-as-code, structure, diagrams, lifecycle | Binding on all |
| 8 | [Git Workflow](08-git-workflow.md) | Trunk-based, Conventional Commits, PR/review rules | Binding on all |
| 9 | [Release Strategy](09-release-strategy.md) | Versioning, GitOps, canary, rollback | Binding on eng |
| 10 | [Definition of Done](10-definition-of-done.md) | The checklist an item must meet to be "Done" | Binding on all |
| 11 | [RFC Process](11-rfc-process.md) | Request for Comments workflow for major changes/features | Governance process |

---

## Hierarchy

```
Product Constitution (supreme)
   ├── DevOS Manifesto (philosophy)
   ├── ADR Register (accepted decisions)
   └── Standards (Eng / Coding / Repo / Docs / Git / Release / DoD)
          └── implemented by → PRD + Engineering Specs + Code
```

Any conflict resolves **upward** (Constitution > ADR > Standards > Specs/PRD).

---

## Consolidated Approval Gate

Per the user's instruction, **no production code starts until these are approved.** Sign below to ratify all eleven as the mandatory governance for the project.

| Role | Name | Approve all 10? | Date | Notes |
|------|------|-----------------|------|-------|
| CEO | __________ | ☐ Yes ☐ No | ______ | |
| CTO | __________ | ☐ Yes ☐ No | ______ | |
| Head of Product | __________ | ☐ Yes ☐ No | ______ | |
| Head of Engineering | __________ | ☐ Yes ☐ No | ______ | |

**Conditions of approval:**
- ☐ Approved as-is (v1.0-rc1 → v1.0).
- ☐ Approved with changes: ________________________ (tracked as revisions, re-versioned).
- ☐ Rejected: ________________________ (revise and re-present).

> Once all four sign **Approve**, the governance set is ratified and implementation planning may begin. Individual documents also carry their own sign-off lines for granular approval; the consolidated gate above is the final gate.

---

*End of Governance Index. Next: implementation planning (post-approval).*
