# API

NetScope exposes a minimal HTTP API for dashboard integration.

Base path: `/api/v1`

## OpenAPI Specification

A full OpenAPI 3.0 specification is available in `docs/openapi.yaml`.

## Endpoints

### GET /api/v1/health

Health check.

Response 200:
```json
{
  "status": "ok",
  "timestamp": "2026-07-25T06:00:00Z"
}
```

### POST /api/v1/scan?target=<target>

Run discovery on a target and return structured results.

Request:
- `target` (query): IP, domain, or hostname

Response 200:
```json
{
  "target": "192.168.1.1",
  "assets": [],
  "services": [],
  "findings": [],
  "graph": {}
}
```

### GET /api/v1/assets

List current assets.

Response 200:
```json
[]
```

### GET /api/v1/services

List current services.

Response 200:
```json
[]
```

### GET /api/v1/findings

List current findings.

Response 200:
```json
[]
```

### GET /api/v1/graph

List graph nodes and edges.

Response 200:
```json
{}
```

### GET /api/v1/changes?path=<db>

Compare current state against latest snapshot.

Response 200:
```json
{
  "changes": [],
  "summary": { "added": 0, "removed": 0, "modified": 0 }
}
```

## Stability

These endpoints are intended to be stable for dashboard integration. Future versions may extend schemas but should not break existing fields without a version bump.
