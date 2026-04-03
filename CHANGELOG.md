# Changelog

All notable changes to scopecheck will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [0.1.0] — 2026-04-03

### Added
- Initial release
- YAML scope file format (version 1)
- `init` command to generate template scope files
- `validate` command with inline, file, and stdin target input
- `audit` command to view scope check history
- IPv4 address and CIDR range matching
- IPv6 address and CIDR range matching
- Domain exact matching (case-insensitive)
- Domain wildcard matching (`*.example.com`)
- Exclusion rules with guaranteed precedence over inclusions
- Text output with colored status indicators
- JSON structured output
- CSV output
- Append-only audit log (JSONL format, 0600 permissions)
- Engagement date window warnings
- Quiet, normal, and verbose output modes
- Exit codes: 0 (all in scope), 1 (out of scope found), 2 (error)
- Cross-platform support: Linux, Windows, macOS (amd64, arm64)
