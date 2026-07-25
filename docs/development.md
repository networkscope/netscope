# Development

## Prerequisites

- Go 1.22 or later
- Git

## Build

```bash
go build ./...
```

## Run

```bash
go run ./cmd/netscope scan 192.168.1.1
```

## Test

```bash
# Unit and integration tests
go test ./...

# Benchmarks
go test -bench=. -benchmem ./internal/core

# Coverage
go test -cover ./...

# Static analysis
go vet ./...
```

## Code style

- Small functions, clear names
- Explicit error handling
- Avoid global mutable state
- Comments only when reasoning is non-obvious
- No AI-style commentary in code

## Release checklist

1. Update version in metadata
2. Run `go test ./...`
3. Run `go vet ./...`
4. Tag the release
5. Build binaries for target platforms
6. Generate and publish checksums
7. Update changelog

## Packaging

- Linux binaries are built via `goreleaser` or equivalent
- Debian packaging may be added later if appropriate
- Man page is maintained under `man/man1/netscope.1`
