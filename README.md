# Gastown Viewer Intent

> Local-first Mission Control for Beads + Gastown-style agent swarms.

## What It Does

- **Board View**: Kanban-style board showing issues by status
- **Issue Details**: Deep-dive into any issue with children and dependency visualization
- **Graph Export**: Dependency graph in multiple formats (DOT, JSON)
- **Events Stream**: Real-time SSE feed of Beads state changes
- **TUI Client**: Terminal UI for keyboard-driven navigation
- **Web UI**: Browser-based dashboard with React

## Quickstart

### Prerequisites

- Go 1.22+
- Node.js 20+
- [Beads](https://github.com/intent-solutions-io/beads) (`bd` CLI in PATH)

### Run

```bash
# Start the daemon + web dev server
make dev

# Or run components separately:
go run ./cmd/gvid              # Daemon on :7070
cd web && npm run dev          # Web UI on :5173

# TUI client (requires daemon running)
go run ./cmd/gvi-tui
```

### Verify

```bash
curl http://localhost:7070/api/v1/health
# {"status":"ok","beads_initialized":true}
```

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Gastown Viewer Intent                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   ┌──────────────┐      ┌──────────────┐      ┌──────────────┐  │
│   │   gvi-tui    │      │   Web UI     │      │  External    │  │
│   │  (Bubbletea) │      │ (React+Vite) │      │   Clients    │  │
│   └──────┬───────┘      └──────┬───────┘      └──────┬───────┘  │
│          │                     │                     │          │
│          └─────────────────────┼─────────────────────┘          │
│                                │                                 │
│                                ▼                                 │
│                    ┌───────────────────────┐                    │
│                    │       gvid Daemon     │                    │
│                    │  (HTTP API + SSE)     │                    │
│                    │   localhost:7070      │                    │
│                    └───────────┬───────────┘                    │
│                                │                                 │
│                                ▼                                 │
│                    ┌───────────────────────┐                    │
│                    │    Beads Adapter      │                    │
│                    │   (shells to `bd`)    │                    │
│                    └───────────┬───────────┘                    │
│                                │                                 │
│                                ▼                                 │
│                    ┌───────────────────────┐                    │
│                    │     .beads/ state     │                    │
│                    │   (managed by bd)     │                    │
│                    └───────────────────────┘                    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/health` | GET | Health check + beads status |
| `/api/v1/issues` | GET | List all issues (supports filters) |
| `/api/v1/issues/:id` | GET | Get single issue with children/deps |
| `/api/v1/board` | GET | Board view (issues grouped by status) |
| `/api/v1/graph` | GET | Dependency graph (JSON) |
| `/api/v1/events` | GET | SSE stream of state changes |

## Roadmap

Tracked via Beads issues in this repo:

**Epic**: `gastown-viewer-intent-btp` — Gastown Viewer Intent MVP

| Issue | Title | Blocked By |
|-------|-------|------------|
| `.1` | Domain model + event schema | — |
| `.2` | Beads adapter via bd CLI | .1 |
| `.3` | Daemon HTTP API + SSE events | .2 |
| `.4` | TUI client (Bubbletea) | .3 |
| `.5` | Web UI (Vite+React) | .3 |
| `.6` | Dev tooling + docs | — |
| `.7` | MVP demo + sanity checks | .4, .5, .6 |

```
Dependency Graph:
.1 → .2 → .3 → .4 → .7
           ↘ .5 → .7
      .6 ────────→ .7
```

Run `bd ready` to see unblocked work. Run `bd blocked` to see dependencies.

## Project Structure

```
gastown-viewer-intent/
├── 000-docs/              # Documentation (flat)
├── cmd/
│   ├── gvid/              # Daemon binary
│   └── gvi-tui/           # TUI binary
├── internal/
│   ├── api/               # HTTP handlers
│   ├── beads/             # Beads adapter (bd CLI wrapper)
│   ├── model/             # Domain types
│   └── store/             # State management
├── web/                   # React + Vite frontend
├── scripts/               # Dev scripts
├── Makefile               # Build targets
└── README.md
```

## License

MIT — See [LICENSE](LICENSE)

## Disclaimer

This project is not affiliated with the original Beads or Gastown authors. It is an independent viewer/dashboard implementation that integrates with those tools via their public CLIs.
