# Phase 6.2–6.4 — App Technical Architectures (Specification)

> **Status:** Draft
> **Depends on:** Phase 6.1 (Frontend Architecture), Phase 3.2 (UX)
> **Scope:** Per-app technical stack, build, and platform-specific concerns for Desktop, Web, Mobile.

---

## 1. Desktop App (Electron + React)

### 1.1 Stack
- **Shell:** Electron 31+ (main + renderer).
- **UI:** React 18 + `packages/ui-kit`.
- **State:** `packages/state` + `crdt-sync`.
- **Terminal:** `xterm.js` over WebSocket to `workspace.exec`.
- **Native:** `electron-updater` for auto-update; `electron-store` for local prefs.

### 1.2 Process Model
```
main (1) → manages windows, preload, IPC
renderer (per window) → React app (project)
worker → CRDT relay client, offline queue
```
- **Multi-window:** one renderer per project; shared worker for bus connection.
- **Security:** `contextIsolation: true`, `nodeIntegration: false`, sandbox on.

### 1.3 Local Capabilities
- Native terminal pane (full PTY).
- File browser with direct diff viewer (CRDT-synced tree).
- Offline queue persisted to disk; flush on reconnect.
- OS notifications for `deploy.completed` / `task.failed`.

### 1.4 Build
- `electron-builder` → dmg/AppImage/exe.
- Code-signed; auto-update via S3/R2.

---

## 2. Web App (Next.js PWA)

### 2.1 Stack
- **Framework:** Next.js 15 (App Router), React 18.
- **UI:** `packages/ui-kit` (React).
- **Realtime:** WS/SSE via `packages/sdk`.
- **Terminal:** `xterm.js` in browser (streamed from workspace).
- **PWA:** `next-pwa` (service worker + manifest).

### 2.2 Architecture
- **Read paths:** RSC where possible; client components for live panes.
- **Realtime:** Client component subscribes to `/intents/:id/stream` via SDK.
- **Sharing:** `/p/[projectId]` deep links; role-scoped views.

### 2.3 Offline
- SW caches shell + last state (IndexedDB).
- Queued intents flush on reconnect.

### 2.4 Build/Deploy
- Static + edge functions (Vercel/self-host).
- CDN for assets; edge WebSocket termination.

---

## 3. Mobile App (React Native + Expo)

### 3.1 Stack
- **Framework:** React Native 0.74 + Expo SDK 51.
- **UI:** `packages/ui-kit` (RN port) — shared tokens, separate primitives.
- **State:** `packages/state` (RN-compatible) + `crdt-sync`.
- **Realtime:** WS via SDK; push via FCM/APNs.
- **Voice input:** `expo-speech` + ASR.

### 3.2 Capabilities
- Tabs: Chat / Tasks / Files / Settings.
- **Quick approve:** notification action buttons (Approve/Reject).
- **Voice:** mic → ASR → intent.
- **Read-only files:** preview only; "Request change" → intent.
- **Deep link:** `devos://p/[projectId]` opens project.

### 3.3 Push & Notifications
- `expo-notifications` for `deploy.completed`, `plan.proposed`, `task.failed`.
- Action buttons map to `plan.approve` / `plan.reject` WS messages.

### 3.4 Build
- EAS Build (iOS/Android); OTA updates via `expo-updates`.

---

## 4. Cross-App Consistency

| Concern | Mechanism |
|---------|-----------|
| Visual | `packages/ui-kit` tokens shared (RN port for mobile) |
| Types | `packages/contracts` single source |
| State | `packages/state` + `crdt-sync` |
| Real-time | `packages/sdk` (WS/SSE) |
| Offline | CRDT + IndexedDB/disk queue |

---

## 5. Tradeoffs & Risks

| Decision | Risk | Mitigation |
|----------|------|------------|
| Electron weight | Bundle/memory | Lazy load, per-window renderers |
| RN UI port | Divergence from web | Shared token lib; component parity tests |
| PWA offline | Stale state | Conflict UI + versioned merge |
| xterm in browser | Security surface | Workspace egress only; no host shell |

---

## 6. Future Extensions

- **Tauri** Desktop for lighter weight.
- **React Server Components** deeper on Web.
- **Shared design tokens** exported to Figma.

---

*End of Phase 6.2–6.4 — App Technical Architectures.*
