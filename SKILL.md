---
name: scopecheck
description: >
  Build, extend, and operate scopecheck — a Go CLI tool that validates pentest targets
  against authorized engagement scope. Use when working on the redhoundinfosec/scopecheck
  repository, when the user asks about scope validation tooling, or when the user needs to
  check whether targets are within an authorized engagement boundary. Covers architecture,
  CLI usage, extension patterns, testing, and safety constraints.
license: MIT
metadata:
  author: Red Hound Information Security LLC
  version: '0.1.0'
  repo: https://github.com/redhoundinfosec/scopecheck
  language: Go
---

# scopecheck Agent Skill

## When to Use This Skill

Use this skill when:
- Working on the `redhoundinfosec/scopecheck` repository
- The user asks about validating pentest scope or engagement boundaries
- The user wants to build or extend a scope checking tool
- The user needs to create scope YAML files for engagements
- The user asks about preventing out-of-scope testing accidents

## What scopecheck Does

scopecheck is a cross-platform Go CLI that validates targets (IPs, CIDRs, domains) against a YAML scope definition. It is a pure safety tool — it never touches the network. It reads a scope file, checks targets against it, and reports which targets are in scope, excluded, or out of scope.

## Core Concepts

### Scope File (YAML)

```yaml
version: 1
engagement:
  name: "Engagement Name"
  id: "ENG-001"
  start: "2026-01-15"
  end: "2026-02-15"
scope:
  include:
    - "192.168.1.0/24"
    - "*.example.com"
    - "10.0.0.50"
  exclude:
    - "192.168.1.1"
    - "mail.example.com"
notes: "Testing window notes"
```

### Matching Rules

1. Exclusions ALWAYS take precedence over inclusions (safety invariant)
2. Supported types: IPv4, IPv6, CIDR ranges, exact domains, wildcard domains (`*.example.com`)
3. Wildcard `*.example.com` matches `sub.example.com` but NOT `example.com` itself
4. Domain matching is case-insensitive

### Exit Codes

- `0` — All targets are in scope
- `1` — One or more targets are out of scope or excluded
- `2` — Error (invalid input, missing files)

## CLI Reference

```bash
# Initialize scope file
scopecheck init
scopecheck init --scope engagement.yaml

# Validate targets
scopecheck validate 192.168.1.50                         # Single target
scopecheck validate 192.168.1.50 10.0.0.1 app.acme.com  # Multiple targets
scopecheck validate --file targets.txt                    # From file
cat targets.txt | scopecheck validate --stdin             # From stdin pipe

# Output formats
scopecheck validate --file targets.txt -f json            # JSON
scopecheck validate --file targets.txt -f csv             # CSV

# Options
scopecheck validate -v --file targets.txt                 # Verbose (show match details)
scopecheck validate -q --file targets.txt                 # Quiet (exit code only)
scopecheck validate --no-color --file targets.txt         # No ANSI colors
scopecheck validate --no-audit --file targets.txt         # Skip audit logging

# Audit trail
scopecheck audit                                          # All entries
scopecheck audit --last 20                                # Last N entries
scopecheck audit -f json                                  # JSON format
scopecheck audit --clear                                  # Clear log
```

## Architecture (for development)

```
internal/scope/scope.go       — Scope, Engagement, ScopeRules structs
internal/scope/parser.go      — YAML parsing and validation
internal/scope/matcher.go     — Matcher engine (CIDR, domain, wildcard matching)
internal/audit/audit.go       — Append-only JSONL audit log
internal/output/output.go     — Text/JSON/CSV renderers
internal/cli/*.go             — CLI commands (init, validate, audit)
```

Single external dependency: `gopkg.in/yaml.v3`.

## Extending scopecheck

### Adding a new target type

1. Add detection in `matcher.go` → `addEntry()`
2. Add fields to `Matcher` struct
3. Implement `checkNewType()` — exclusions first, then inclusions
4. Wire into `Check()`
5. Add tests, update docs and template

### Adding a new output format

1. Add render function in `output/output.go`
2. Wire into `Render()` switch
3. Update CLI help text

## Safety Constraints

- NEVER add network connectivity
- NEVER weaken exclusion precedence
- NEVER add features that help circumvent scope
- ALWAYS maintain audit log capability
- ALWAYS update tests when changing matcher logic

## Build and Test

```bash
go build -o scopecheck ./cmd/scopecheck/
go test ./... -v -count=1
go vet ./...
make release   # Cross-compile for linux/windows/darwin × amd64/arm64
```

## Common Workflows

### Pre-scan gate

```bash
nmap -sL -n 10.0.0.0/24 | awk '/report/{print $5}' | scopecheck validate --stdin
if [ $? -ne 0 ]; then echo "ABORT: Out-of-scope targets"; exit 1; fi
```

### CI/CD pipeline integration

```bash
scopecheck validate -q --file targets.txt -f json > scope-result.json
[ $? -eq 0 ] || { echo "Scope check failed"; exit 1; }
```

### Post-engagement audit export

```bash
scopecheck audit -f json > audit-$(date +%Y%m%d).json
```
