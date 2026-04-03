# FAQ

## Does scopecheck make any network connections?

No. scopecheck is a pure validation tool. It reads local files and performs string/CIDR matching. It never touches the network.

## Can I use scopecheck for bug bounty programs?

Yes. Define the bug bounty scope in your scope YAML file (authorized domains, IP ranges, and any exclusions) and validate targets before testing.

## Why doesn't wildcard *.example.com match example.com?

By design. `*.example.com` matches subdomains of `example.com`, not the bare domain itself. If you need to include the bare domain, add it as a separate entry:

```yaml
scope:
  include:
    - "*.example.com"
    - "example.com"
```

## Can I use multiple scope files?

Yes. Use the `--scope` flag to specify which file to use:

```bash
scopecheck validate --scope phase1.yaml --file targets.txt
scopecheck validate --scope phase2.yaml --file targets.txt
```

## Is the audit log tamper-proof?

No. The audit log is an append-only JSONL file with restricted permissions (0600), but it is not cryptographically signed. Treat it as a convenience record, not forensic evidence.

## How do I use scopecheck with nmap?

```bash
# List scan to generate IPs, then validate
nmap -sL -n 10.0.0.0/24 | awk '/report/{print $5}' | scopecheck validate --stdin
```

## What happens if I run scopecheck outside the engagement window?

You get a warning on stderr, but validation still proceeds. This is a safety warning, not a hard block.

## Can scopecheck resolve DNS names?

Not in v0.1. DNS resolution is planned for v0.2 with an explicit `--resolve` flag.
