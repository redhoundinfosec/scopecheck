# Agent Instructions for scopecheck

This document tells AI coding agents how to work with the `scopecheck` codebase — build, test, extend, and contribute safely.

## Project Overview

`scopecheck` validates targets (IPs, CIDRs, domains) against a YAML-defined engagement scope. It is a **safety tool** for pentesters — it never touches the network. Written in Go with zero external dependencies beyond `gopkg.in/yaml.v3`.

## Quick Commands

```bash
# Build
go build -o scopecheck ./cmd/scopecheck/

# Test
go test ./... -v -count=1

# Lint
go vet ./...

# Cross-compile all platforms
make release

# Run
./scopecheck init                                   # Create scope.yaml template
./scopecheck validate 192.168.1.50                  # Check single target
./scopecheck validate --file targets.txt            # Check target list
./scopecheck validate --file targets.txt -f json    # JSON output
./scopecheck audit --last 20                        # View audit trail
```

## Architecture

```
cmd/scopecheck/main.go          Entry point — calls cli.Run(os.Args[1:])
internal/
  scope/
    scope.go                     Core model: Scope, Engagement, ScopeRules, DateRange
    parser.go                    YAML parsing via gopkg.in/yaml.v3, validation, template generation
    matcher.go                   Matching engine: Matcher, MatchResult, MatchStatus enum
                                 Supports: IPv4, IPv6, CIDR, exact domain, wildcard domain
                                 Rule: exclusions ALWAYS take precedence over inclusions
  audit/
    audit.go                     Append-only JSONL audit log (0600 perms), read/write
  output/
    output.go                    Renderers: text (ANSI colored), JSON, CSV
  cli/
    root.go                      CLI dispatcher, global flag parsing, help text
    init.go                      `init` subcommand — writes template scope.yaml
    validate.go                  `validate` subcommand — loads scope, builds matcher,
                                 collects targets (inline/file/stdin), renders output, writes audit
    audit.go                     `audit` subcommand — reads and displays audit log
```

## Key Design Decisions

1. **Exclusions always win.** If a target matches both include and exclude, it is marked `excluded`. This is a safety invariant — never change it.
2. **No DNS resolution.** Matching is pure string/CIDR comparison. DNS resolution is a planned v0.2 feature behind an explicit `--resolve` flag.
3. **No third-party CLI framework.** Flag parsing is hand-rolled to keep the dependency count at 1 (yaml.v3 only).
4. **Exit codes are part of the API.** 0 = all in scope, 1 = violations found, 2 = error. Scripts depend on this.
5. **Audit log is append-only JSONL.** One JSON object per line. File permissions are 0600. Not cryptographically signed.

## How to Add a New Feature

### Adding a new target type (e.g., URL matching)

1. Edit `internal/scope/matcher.go`:
   - Add detection logic in `addEntry()` to identify the new type
   - Add storage fields to the `Matcher` struct
   - Implement matching in a new `checkURL()` method
   - Wire it into `Check()` — check exclusions first, then inclusions
2. Add tests in `internal/scope/matcher_test.go`
3. Update the scope file format docs in `docs/scope-file-format.md`
4. Update the YAML template in `parser.go` → `TemplateYAML()`
5. Update README.md examples and target type table

### Adding a new output format

1. Edit `internal/output/output.go`:
   - Add a new `Format` constant
   - Implement a `renderXxx()` function
   - Wire it into the `Render()` switch
2. Update CLI help text in `internal/cli/root.go`
3. Add an example in README.md

### Adding a new subcommand

1. Create `internal/cli/mycommand.go`
2. Wire it into the switch in `cli/root.go` → `Run()`
3. Add it to the help text in `printUsage()`
4. Add tests if the command has non-trivial logic

## Testing Conventions

- All test files are colocated: `foo.go` → `foo_test.go`
- Use table-driven tests with descriptive subtests
- Test both positive and negative paths
- The `testScope()` helper in `matcher_test.go` provides a reusable fixture
- Edge cases to always test: empty input, overlapping CIDR ranges, case sensitivity, exclusion precedence

## Safety Rules for Agents

1. **Never add network connectivity.** scopecheck must remain a pure local validation tool.
2. **Never weaken exclusion precedence.** Exclusions must always override inclusions.
3. **Never remove the audit log.** It's a safety feature. You may make it optional (it already is via `--no-audit`).
4. **Never add features that help circumvent scope.** If a feature request could help someone bypass scope boundaries, reject it.
5. **Always update tests** when changing matcher logic. Matcher correctness is critical.

## File Format Reference

Scope files are YAML v1:

```yaml
version: 1
engagement:
  name: "string (required)"
  id: "string (optional)"
  start: "YYYY-MM-DD (optional)"
  end: "YYYY-MM-DD (optional)"
scope:
  include:    # required, at least one entry
    - "192.168.1.0/24"
    - "*.example.com"
  exclude:    # optional
    - "192.168.1.1"
notes: "free text (optional)"
```

## Dependencies

- Go 1.22+
- `gopkg.in/yaml.v3` (only external dependency)
- No CGO required
