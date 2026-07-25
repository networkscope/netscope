# Project Manual

## Stack & Commands
- Stack: Go 1.22+, SQLite (embedded via modernc.org/sqlite), Cobra CLI
- Install: `go mod tidy && go build ./...`
- Test: `go test ./... && go vet ./...`
- Build: `go build ./...`
- Benchmarks: `go test -bench=. -benchmem ./internal/core`

## Strict Guardrails
- DO NOT refactor unrequested files
- DO NOT remove comments or documentation
- Maintain existing architecture patterns
- All code changes must pass `go test ./...` and `go vet ./...`

## Quick Reference
- `netscope scan <target>` - Scan a target for assets and services
- `netscope scan <target> --save <path>` - Scan and save to SQLite
- `netscope changes <path>` - Compare current state against latest snapshot
- `netscope serve <address>` - Start the API server (serves frontend at `/` and API at `/api/v1`)

Output formats: `text`, `json`, `csv` (use `-o <format>`)

## Project Structure
```
cmd/netscope/          - Main entry point
internal/
  cli/                 - CLI commands and output formatting
  core/                - Engine and discovery logic
  api/                 - HTTP API server (includes embedded web/)
  storage/             - SQLite persistence
  changes/             - Snapshot diff logic
  assets/              - Asset discovery
  findings/            - Rule-based analysis
pkg/models/            - Data models
web/                   - Embedded frontend (moved to internal/api/web/)
docs/                  - Documentation
  openapi.yaml         - OpenAPI 3.0 spec
```