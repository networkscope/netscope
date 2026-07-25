<p align="center">
<br/>
<img src="https://github.com/networkscope/netscope/blob/main/branding/wordmark.png?raw=true" alt="NetScope Logo" width=250>
<h3 align="center">NetScope</h3>
<p align="center">
NetScope is an open-source cybersecurity reconnaissance and analysis platform. It discovers assets, models services, identifies security findings, and builds relationship graphs in authorized environments.
</p>
<p align="center">
<img alt="Static Badge" src="https://img.shields.io/badge/NetScope-v0.1.0-6d7592?logo=github&labelColor=45455e&link=https%3A%2F%2Fgithub.com%2Fnetworkscope%2Fnetscope%2Freleases">
</p>

> [!NOTE]\
> NetScope is designed for authorized security assessments only. Always ensure you have permission before scanning.

## Features

- **Asset Discovery:** Discover hosts, IPs, domains, and applications from target inputs
- **Service Modeling:** Map services to assets with port, protocol, and software identification
- **Security Findings:** Rule-based analysis produces severity-ranked security observations
- **Relationship Graphs:** Build and visualize asset-service relationships
- **SQLite Persistence:** Save and load assessments with full history tracking
- **Change Tracking:** Compare snapshots to identify added, removed, and modified assets
- **Multiple Output Formats:** Text, JSON, and CSV output for flexibility

## Install

```bash
# Clone the repository
git clone https://github.com/networkscope/netscope.git
cd netscope

# Build
go build ./...

# Or install directly
go install ./cmd/netscope
```

## Requirements

- Go 1.22 or later
- SQLite is embedded via `modernc.org/sqlite` (pure Go, no CGO)

## Quick Start

```bash
# Scan a target
netscope scan 192.168.1.1

# Scan and save to SQLite
netscope scan example.com --save assessment.db

# View changes since last snapshot
netscope changes assessment.db

# JSON output
netscope scan 10.0.0.1 -o json

# Start API server
netscope serve :8080
```

## CLI Reference

### `netscope scan <target>`

Runs asset discovery on a target. Produces asset, service, finding, and graph output.

Flags:
- `-s, --save <path>`: save assessment to SQLite database
- `-o, --output <format>`: `text`, `json`, or `csv` (default: text)
- `-v, --verbose`: enable verbose output
- `-q, --quiet`: suppress non-essential output

### `netscope changes <path>`

Compares current state against the latest snapshot in the database.

### `netscope serve <address>`

Starts the HTTP API server, serving the web frontend at `/` and the API at `/api/v1`.

## Architecture

```
Assessment
├── Assets       — discovered hosts, IPs, domains, applications
├── Services     — ports, protocols, identified software
├── Findings     — rule-generated security observations
└── Graph        — relationships between entities
```

The engine coordinates discovery, analysis, and persistence. Results are stored in-memory and optionally persisted to SQLite via the storage subsystem.

## Contribute

Please follow the instructions in [CONTRIBUTING.md](https://github.com/networkscope/netscope/blob/main/CONTRIBUTING.md) for guidelines on contributing to NetScope.

## Contributors ❤️

Thank you to:

<a href="https://github.com/networkscope/netscope/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=networkscope/netscope" />
</a>

Made with [contrib.rocks](https://contrib.rocks).

## License

NetScope is released under GNU General Public License v3.0. See [LICENSE](https://github.com/networkscope/netscope/blob/main/LICENSE) for details.

&copy; Copyright NetScope contributors 2026
All Rights Reserved!