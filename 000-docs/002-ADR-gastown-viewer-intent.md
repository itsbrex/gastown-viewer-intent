# ADR-001 — Gastown Viewer Intent MVP Architecture

**Status**: Accepted
**Date**: 2026-01-01
**Repo**: intent-solutions-io/gastown-viewer-intent
**Docs location**: 000-docs/ (flat)

---

## Context

We need a local-first "mission control" viewer for Beads-managed work. Users want a board view, issue details (children + dependencies), a dependency graph export, and optional live updates. The solution must be quick to run in any repo, minimize coupling to Beads internals, and support both terminal and web workflows.

**Background Concepts**:
- **Beads (bd)**: A local issue tracker where work is represented as issues with first-class dependency support. Sequencing is controlled by explicit dependency bonds ("A blocks B"). The `bd ready` command shows unblocked issues.
- **Gastown**: An agent-orchestration architecture (inspiration only). This repo does NOT depend on upstream Gastown code — we build a standalone viewer.

---

## Decision

Adopt a **control-plane architecture**:

1. **Local daemon (Go)** provides a stable HTTP API on localhost:7070
2. **UI clients (TUI + Web)** consume the daemon API only
3. **Beads integration** for MVP shells out to the `bd` CLI (no direct .beads/ parsing)
4. **Documentation** is stored flat in 000-docs/ only

---

## Drivers

- **Fast MVP iteration** while keeping a stable internal contract (API)
- **Reduce breakage risk** from Beads storage/schema changes
- **Support SSH/terminal users and visual users** without duplicating logic
- **Extensibility** for future integrations (SSE, Slack, mobile, multi-repo) without rewriting clients

---

## Key Decisions

### D1: Language and UI Stack
- **Go** for daemon + TUI
- **Bubbletea** for TUI (mature, composable)
- **Vite + React + TypeScript** for web UI

**Rationale**: Go provides single-binary distribution, excellent CLI tooling, and Bubbletea is the de facto standard for Go TUIs. React offers fast development for the web dashboard.

### D2: Data Source Contract
- MVP uses `bd` CLI as the integration contract
- Implement `internal/beads.Adapter` that executes bd commands and parses output
- Graceful degradation on parse errors (return partial data, never panic)

**Rationale**: The bd CLI is stable; .beads/ internals may change. Shelling out is slower but safer for MVP.

### D3: API-First Internal Boundary
- Daemon exposes `/api/v1/*` endpoints (health, issues, board, graph, events)
- Clients MUST NOT call `bd` directly
- All data flows through the daemon

**Rationale**: Single parsing surface, consistent caching opportunity, enables remote clients later.

### D4: Live Updates
- **Preferred**: SSE at `GET /api/v1/events`
- **Acceptable MVP fallback**: Polling in web/TUI

**Decision**: Implement SSE in daemon. If SSE proves complex for clients, polling is acceptable for MVP launch.

### D5: Documentation Standard
- All documentation files live flat in `000-docs/`
- No `docs/` directory, no subfolders
- Files are numbered: `001-PRD-*.md`, `002-ADR-*.md`, `003-API-*.md`

**Rationale**: Flat structure is scannable, sortable, and grep-friendly.

---

## Consequences

### Positive
- **Separation of concerns**: Integration logic (daemon) is separate from presentation (TUI/Web)
- **One parsing surface**: Only daemon parses bd output; clients get clean JSON
- **Easier testing**: Mock Adapter for unit tests, test API contract independently
- **Future-proof**: Can add desktop wrapper, mobile remote client, or multi-repo support

### Negative
- **Requires running daemon**: Extra process to manage
- **CLI parsing can be brittle**: Must invest in robust error handling and graceful degradation
- **SSE adds complexity**: May defer to polling for MVP simplicity

---

## Alternatives Considered

### A1: Web-only app that reads .beads/ directly
**Rejected**: High coupling to storage format; harder to keep stable; duplicates logic if TUI added later.

### A2: Fork beads_viewer (bv) and extend it
**Rejected for MVP**: Fork maintenance burden; product shape differs (control plane + multi-clients). May revisit post-MVP as an optional insights provider.

### A3: Single monolithic TUI that calls bd directly
**Rejected**: Blocks web UI reuse; no stable contract; harder to add mobile/remote later.

### A4: gRPC instead of HTTP/JSON
**Rejected for MVP**: HTTP is simpler, curl-debuggable, and sufficient for local use.

---

## Implementation Notes

### Fail-Fast Requirements
Daemon should fail fast with actionable errors when:
- `bd` is missing from PATH → "ERROR: bd CLI not found. Install from..."
- `.beads/` not initialized → 503 + "Run 'bd init' first"

### Graceful Degradation
Adapter should degrade gracefully:
- Return partial data when fields are missing (with warnings)
- Never panic on unexpected bd output
- Log parse errors for debugging

### Startup Flow
```
1. Check bd available → fail fast if not
2. Check .beads/ exists → return 503 on health if not
3. Start HTTP server on :7070
4. Accept client connections
```

---

## Follow-ups (Post-MVP)

1. **bv integration**: Add optional insights provider via beads_viewer
2. **Persistence cache (sqlite)**: Faster startup on large graphs
3. **Remote mode**: Mobile client talks to daemon over Tailscale/LAN
4. **Write operations**: Create/update issues via API
5. **Multi-repo**: Dashboard for multiple Beads repos
6. **Authentication**: For team deployments

---

## Related Documents

- `000-docs/001-PRD-gastown-viewer-intent.md` — Product requirements
- `000-docs/003-API-gastown-viewer-intent.md` — API specification

---

## Document History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0.0 | 2026-01-01 | Claude | Initial ADR |
