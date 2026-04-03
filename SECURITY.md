# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| 0.1.x   | Yes       |

## Reporting a Vulnerability

If you discover a security vulnerability in scopecheck, please report it responsibly:

1. **Email:** security@redhoundinfosec.com
2. **Do NOT** open a public GitHub issue for security vulnerabilities.
3. Include a clear description, steps to reproduce, and potential impact.
4. We will acknowledge receipt within 48 hours.
5. We will provide a fix timeline within 7 days.

## Threat Model

scopecheck is a local validation tool that:
- Never makes network connections
- Reads only local files (scope YAML, target lists, audit log)
- Writes only to the local audit log file
- Does not execute external commands
- Does not process untrusted binary data

### Potential risks:
- **Malicious scope file:** A crafted YAML file could potentially cause excessive memory use. Mitigated by using a well-tested YAML parser.
- **Audit log tampering:** The audit log is append-only but not cryptographically signed. It should be treated as a convenience record, not a tamper-proof evidence chain.
- **Path traversal:** Scope file and target file paths are user-provided. Standard OS path handling applies.

## Design Safety

scopecheck is intentionally designed with safety constraints:
- It cannot be used to scan or attack targets
- It has no network capabilities
- Exclusions always override inclusions (fail-safe)
- The engagement date window produces warnings when outside authorized periods
