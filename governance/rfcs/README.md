# DevOS — RFC Registry

> Part of Governance v1.0-rc1. See [`../11-rfc-process.md`](../11-rfc-process.md) for the process.

## Purpose
This directory holds all **Request for Comments** proposals. RFCs are the proposal layer for major architectural changes and significant features. Accepted RFCs may spawn ADRs (`../03-adr.md`) and spec updates.

## Folders
| Folder | Contents |
|--------|----------|
| [`Accepted/`](Accepted/) | Ratified, live proposals (the active record). |
| [`Rejected/`](Rejected/) | Proposals declined, with rationale (institutional memory). |
| [`Deprecated/`](Deprecated/) | Superseded proposals, linked to successors. |

## Index

| RFC | Title | Status | Date |
|-----|-------|--------|------|
| — | (none yet — first RFC created on demand) | — | — |

## Rules
- Number sequentially: `RFC-0001`, `RFC-0002`, …
- Move the file **between folders** as status changes; never delete.
- Template: [`000-template.md`](../000-template.md).
- No implementation before `Accepted`.
