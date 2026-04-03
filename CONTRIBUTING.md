# Contributing to scopecheck

Thank you for your interest in contributing.

## Principles

1. **Safety first.** scopecheck is a safety tool. Never add features that could be used to circumvent scope boundaries.
2. **Keep it simple.** Resist feature creep. Every addition must justify its complexity.
3. **Cross-platform.** All changes must work on Linux, Windows, and macOS.
4. **Test everything.** Every bug fix must include a regression test.

## Development Setup

```bash
git clone https://github.com/redhoundinfosec/scopecheck.git
cd scopecheck
go build ./cmd/scopecheck/
go test ./...
```

## Submitting Changes

1. Fork the repository.
2. Create a feature branch from `main`.
3. Write tests for your changes.
4. Run `go test ./...` and `go vet ./...`.
5. Submit a pull request with a clear description.

## Code Standards

- Run `gofmt` on all Go files.
- Keep functions short and purposeful.
- Use explicit error handling.
- Add comments only where they improve understanding.

## Scope of Contributions

Contributions welcome:
- Bug fixes
- Test coverage improvements
- Documentation improvements
- New target type support (with tests)
- Output format improvements
- Cross-platform compatibility fixes

Contributions NOT accepted:
- Network scanning or probing features
- DNS resolution without explicit opt-in
- GUI or web interface
- Features that bypass or weaken scope enforcement
- Dependencies without strong justification
