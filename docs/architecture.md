# Architecture

NetScope is structured around a scanner engine that coordinates discovery, analysis, and persistence.

## Packages

| Package | Responsibility |
|---------|----------------|
| `cmd/netscope` | Entry point, CSRF |
| `internal/cli` | Cobra command definitions and output formatting |
| `internal/core` | Engine orchestration: scan, save, load, snapshot, diff |
| `internal/assets` | Target parsing and asset registry |
| `internal/services` | Service registry and discovery stubs |
| `internal/findings` | Rule-based finding evaluator and registry |
| `internal/graph` | Relationship graph builder and text serialization |
| `internal/storage` | SQLite-backed persistence layer |
| `internal/changes` | Snapshot model and diff engine |
| `pkg/models` | Standalone domain structs shared across packages |

## Design constraints

- No global mutable state except for CLI flag variables
- Dependencies flow inward: `cli -> core -> subsystems -> models`
- Packages are cohesive; avoid circular dependencies
- Avoid generic utility packages
- Interfaces only where they reduce coupling

## Concurrency

The current implementation is sequential. Future network scanning and fingerprinting should use worker pools with bounded concurrency, deadlines, and context cancellation.

## Persistence

SQLite is used for local-first storage. Data is organized in tables for assets, services, findings, graph nodes, graph edges, and snapshots. Snapshots serialize the full assessment state as JSON payloads so historical data is queryable without custom SQL for each entity type.

## Output

- Text: formatted for terminal readability
- JSON: stable machine-readable format with documented fields
- Future formats may include CSV, HTML, or structured reporting

## Extensibility

Future modules (TLS analysis, DNS analysis, cloud discovery, container analysis, vulnerability intelligence) should be added as internal packages with clean interfaces. Dynamic plugin loading via `go plugin` is not recommended due to known limitations; prefer interface-based composition.
