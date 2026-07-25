# CLI usage

## Global flags

- `-o, --output <format>`: `text` or `json`
- `-v, --verbose`: enable verbose output
- `-q, --quiet`: suppress non-essential output
- `-h, --help`: show help

## Commands

### scan <target>

Scan a target to discover assets and services.

```bash
netscope scan 192.168.1.1
netscope scan example.com -o json
netscope scan web01 --save assessment.db
```

### changes <path>

Show added, removed, or modified entities since the latest snapshot in the database.

```bash
netscope changes assessment.db
netscope changes assessment.db -o json
```

## Exit codes

- `0`: success
- `1`: command execution error
- `2`: invalid arguments

## Output format stability

JSON output is intended to be stable. If the schema changes, the change will be documented in release notes. Do not rely on undocumented fields.
