# DevOS — Product Constitution

> **Status:** RELEASE CANDIDATE — For Ratification (no production code until ratified)
> **Version:** 1.0-rc1
> **Owner:** CTO (ratified by CEO + CTO + Head of Product)
> **Supersedes:** None
> **Companion:** PRD (`/product/PRD.md`), Engineering Specs (`/specs/`)

---

## Preamble

We are building DevOS — an AI-Native Development Operating System. This Constitution establishes the **immutable, governing principles** that bind every product, architecture, and engineering decision. Where the PRD defines *what we build* and the specs define *how*, this Constitution defines *what we are never allowed to violate*.

Articles herein may only be amended through the process in Article VI. Until ratified, this document is provisional.

---

## Article I — Purpose

1. DevOS exists to let any human **build, manage, deploy, and monitor software through natural-language intent**, controllable from any surface they choose.
2. The real product is an **AI orchestration platform**; the IDE is only one of many interfaces.
3. Every decision shall serve the human directing intent — their safety, control, and outcome quality come first.

---

## Article II — Non-Negotiable Tenets

These tenets are **inviolable**. Any proposal contradicting a tenet must be rejected or routed through Amendment (Article VI).

| # | Tenet | Binding Rule |
|---|-------|--------------|
| T1 | **Channel-Agnostic Core** | All control surfaces (Desktop, Web, Mobile, WhatsApp, Telegram, Discord, Slack, Voice, REST) funnel to one uniform intent. Core logic shall not depend on any channel. |
| T2 | **No Vendor Lock-In** | No AI/tool/cloud provider shall be hardcoded. All providers sit behind a port/adapter with capability flags. Switching providers requires zero agent-logic changes. |
| T3 | **Human-in-the-Loop** | Approval gates for plan, deploy, and secret access are first-class. Autonomy never overrides a required human checkpoint. |
| T4 | **Workspace Isolation & Secret Safety** | Agents execute only in isolated workspaces. Raw secrets are never placed in agent context; they are resolved at proxy egress only. |
| T5 | **Budget Governance** | Token and cost spend is capped per org/project/user. Unbounded autonomous spend is prohibited. |
| T6 | **Observability First** | Every service and agent action is traced, metered, and logged from day one. No feature ships blind. |
| T7 | **Event-Driven & Decoupled** | Cross-context communication occurs only via the message bus. No service calls another's datastore. |
| T8 | **Provider Abstraction** | LLM, Tool, Deploy, Channel, and Vector capabilities are expressed as interfaces; concrete implementations are plugins. |
| T9 | **Offline-First Where Possible** | Client state syncs via CRDT; intents queue locally and flush on reconnect. |
| T10 | **Security by Design** | Least-privilege, OIDC short-lived credentials, full audit trail on every mutation. Security is a prerequisite, not a phase. |
| T11 | **AI Transparency** | Every autonomous decision is explainable. For each action the platform must expose: why it happened, which agent performed it, which provider was selected, cost estimate, files modified, and rollback strategy. |
| T12 | **Open Standards** | DevOS prioritizes open standards (MCP, OpenAPI, Git, OCI, OAuth, Webhooks). Proprietary protocols require documented justification. |

---

## Article III — Product Boundaries (What We Are Not)

To preserve focus, DevOS explicitly **renounces** the following scopes:

1. We are **not** a code-hosting platform (we integrate GitHub; we do not replace it).
2. We are **not** a general-purpose chatbot or a human-developer replacement.
3. We are **not** an IDE in the traditional sense — an editor is one interface among nine.
4. We do **not** ship features that violate Article II, regardless of competitive pressure or revenue opportunity.

---

## Article IV — Stakeholder Primacy

1. **Users first:** their data sovereignty, safety, and control outweigh internal metrics.
2. **Developers second:** the platform must remain comprehensible, debuggable, and extensible by its own builders.
3. **Business third:** commercial goals are pursued only within the bounds of Articles II–III.

---

## Article V — Governance & Authority

1. This Constitution is the **supreme governance document**, ranking above PRD, specs, ADRs, and standards.
2. ADRs (Architecture Decision Records) implement these tenets; a conflict resolves in favor of this Constitution.
3. The CTO is the custodian; the CEO + CTO + Head of Product ratify amendments.

---

## Article VI — Amendment Process

1. Any stakeholder may propose an amendment via a PR to this file, referencing rationale and impact.
2. Amendment requires **explicit approval by CEO + CTO + Head of Product**.
3. Accepted amendments increment the Version (e.g., 1.0 → 1.1) and are logged below.
4. Tenets (Article II) require the highest bar: unanimous approval and a recorded justification.

### Amendment Log
| Version | Date | Change | Approved By |
|---------|------|--------|-------------|
| 1.0-draft | 2026-07-20 | Initial constitution | _pending_ |
| 1.0-rc1 | 2026-07-20 | Added T11 AI Transparency, T12 Open Standards; added Article VIII Vision Stability | _pending_ |

---

## Article VII — Ratification

This Constitution takes effect only upon signature below. Until then, no production code may be written.

| Role | Name | Ratify? | Date |
|------|------|---------|------|
| CEO | __________ | ☐ Yes ☐ No | ______ |
| CTO | __________ | ☐ Yes ☐ No | ______ |
| Head of Product | __________ | ☐ Yes ☐ No | ______ |

---

## Article VIII — Vision Stability

1. The DevOS **core product vision** — an AI-Native Development Operating System where software is built, managed, deployed, and monitored through natural-language intent across any surface by a coordinated, safe, provider-agnostic agent team — is **constant**.
2. **Technology may change. Architecture may evolve. Implementation may change. Vision remains constant.**
3. Any proposal to alter the core vision requires **unanimous approval** of the governance ratifiers (CEO + CTO + Head of Product + Head of Engineering) with a recorded, public justification.
4. Vision changes are processed as amendments to this Constitution (Article VI) and version-bumped; they are never made through ADRs or specs alone.

---

*End of Product Constitution v1.0-rc1.*
