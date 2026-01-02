# CLAUDE.md — Gastown Viewer Intent

> Context file for Claude Code sessions. Self-contained; no prior knowledge required.

## Project Overview

**Gastown Viewer Intent** is a local-first Mission Control for **Beads** (a local issue tracker with dependency support). It provides visual board views, issue details, and dependency graphs via a daemon API, TUI, and Web UI.

## Core Concepts

**Beads (bd)**: A local issue tracker where work is represented as issues with first-class dependency support. Issues can block other issues. The `bd ready` command shows unblocked work. The `bd blocked` command shows the dependency graph.

**Gastown**: An agent-orchestration architecture (inspiration only). This repo is a standalone viewer that integrates with Beads via the `bd` CLI.

## Tech Stack

- **Go**: Daemon (`cmd/gvid`) + TUI (`cmd/gvi-tui`)
- **Bubbletea**: TUI framework
- **Vite + React + TypeScript**: Web UI (`web/`)
- **Beads Adapter**: Shells to `bd` CLI (no direct .beads/ parsing)

## Quick Commands

```bash
# Prerequisites check
bd --version && go version && node -v

# Development
make dev           # Daemon + web in parallel
go run ./cmd/gvid  # Daemon only
go run ./cmd/gvi-tui  # TUI (needs daemon)

# Beads work tracking
bd ready           # Show unblocked issues
bd blocked         # Show dependency graph
bd show <id>       # View issue details
```

## Project Structure

```
gastown-viewer-intent/
├── 000-docs/          # Flat documentation (PRD, ADR, API spec)
├── cmd/gvid/          # HTTP daemon binary
├── cmd/gvi-tui/       # TUI binary
├── internal/api/      # HTTP handlers
├── internal/beads/    # Beads adapter (shells to bd)
├── internal/model/    # Domain types
├── internal/store/    # State management
├── web/               # React frontend
├── Makefile           # Build targets
└── .beads/            # Beads issue database
```

## Beads Work Graph

Epic: `gastown-viewer-intent-btp`

```
.1 (model) → .2 (adapter) → .3 (daemon) → .4 (TUI) → .7 (demo)
                                       → .5 (Web) → .7
                            .6 (tooling) ─────────→ .7
```

Run `bd ready` to see what's unblocked. Pick the top item.

## Fail-Fast Behavior

If `bd` is not installed, the daemon MUST return 503 with:
```json
{"error": "bd CLI not found. Install from https://github.com/intent-solutions-io/beads"}
```

If `.beads/` is not initialized, return 503 with:
```json
{"error": "Beads not initialized. Run 'bd init' in your project directory."}
```

## API Summary

| Endpoint | Description |
|----------|-------------|
| GET /api/v1/health | Health + beads status |
| GET /api/v1/issues | List issues |
| GET /api/v1/issues/:id | Issue details |
| GET /api/v1/board | Board view |
| GET /api/v1/graph | Dependency graph |
| GET /api/v1/events | SSE stream |

See `000-docs/003-API-gastown-viewer-intent.md` for full spec.

## Session Recovery

After context compaction, run:
```bash
bd ready              # What's unblocked?
bd show <id>          # Details of top item
```

Then continue from where you left off based on the Beads state.
