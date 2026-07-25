# NetScope

<p align="center">
<br/>
<img src="https://github.com/networkscope/netscope/blob/main/branding/wordmark.png?raw=true" alt="NetScope Logo" width=200>
</p>

NetScope is an open-source cybersecurity reconnaissance and analysis platform. It discovers assets, models services, identifies security findings, and builds relationship graphs in authorized environments.

## Features

- **Asset Discovery:** Discover hosts, IPs, domains, and applications from target inputs
- **Service Modeling:** Map services to assets with port, protocol, and software identification
- **Security Findings:** Rule-based analysis produces severity-ranked security observations
- **Relationship Graphs:** Build and visualize asset-service relationships
- **SQLite Persistence:** Save and load assessments with full history tracking
- **Change Tracking:** Compare snapshots to identify added, removed, and modified assets
- **Multiple Output Formats:** Text, JSON, and CSV output for flexibility

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

## Installation

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

## Documentation

- [Architecture](architecture.md) - System design and components
- [CLI Reference](cli.md) - Command-line interface documentation
- [API Reference](api.md) - HTTP API endpoints
- [Development](development.md) - Building and testing
- [Security](security.md) - Security model and best practices

## License

NetScope is released under GNU General Public License v3.0.