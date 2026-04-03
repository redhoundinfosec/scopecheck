# Roadmap

## v0.1.0 (Current — MVP)
- [x] YAML scope file format
- [x] IPv4/IPv6 address and CIDR matching
- [x] Domain exact and wildcard matching
- [x] Exclusion support with precedence
- [x] Text, JSON, CSV output
- [x] Stdin, file, and inline target input
- [x] Audit log
- [x] Engagement date window warnings
- [x] Cross-platform support (Linux, Windows, macOS)

## v0.2.0 (Planned)
- [ ] DNS resolution mode (`--resolve`) — resolve domains before matching
- [ ] Import scope from nmap XML (`--import-nmap`)
- [ ] Port-level scope restrictions
- [ ] Engagement date enforcement mode (block, not just warn)

## v0.3.0 (Future)
- [ ] Scope diff — compare two scope files and show changes
- [ ] Import from Burp Suite project files
- [ ] Import from Cobalt Strike malleable profiles
- [ ] Team sharing with signed scope files
- [ ] CI/CD integration examples and GitHub Action

## v1.0.0 (Stable)
- [ ] Comprehensive integration test suite
- [ ] Performance benchmarks
- [ ] Stable CLI contract (no breaking changes after v1)
- [ ] Homebrew and APT package support
