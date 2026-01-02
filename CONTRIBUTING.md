# Contributing to Gastown Viewer Intent

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork locally
3. Install dependencies:
   - Go 1.22+
   - Node.js 20+
   - Beads (`bd` CLI)

4. Run `make dev` to start development

## Development Workflow

1. Check `bd ready` to see available issues
2. Pick an issue and update its status: `bd update <id> --status in_progress`
3. Make your changes
4. Test: `make test`
5. Commit with clear messages
6. Push and open a PR

## Code Style

- Go: Follow standard Go conventions, run `go fmt`
- TypeScript/React: ESLint + Prettier (configured in web/)
- Commits: Use conventional commits when possible

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if needed
3. Reference the Beads issue ID in your PR
4. Request review from maintainers

## Questions?

Open a discussion or issue on GitHub.
