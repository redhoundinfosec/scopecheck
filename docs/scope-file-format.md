# Scope File Format

scopecheck uses a YAML file to define the authorized engagement scope.

## Format Version

Current version: `1`

The `version` field is required and must match the tool's supported format version.

## Structure

```yaml
version: 1

engagement:
  name: "Engagement Name"        # Required
  id: "ENG-ID"                   # Optional
  start: "2026-01-15"            # Optional (YYYY-MM-DD)
  end: "2026-02-15"              # Optional (YYYY-MM-DD)

scope:
  include:                       # Required (at least one entry)
    - "192.168.1.0/24"
    - "*.example.com"
  exclude:                       # Optional
    - "192.168.1.1"
    - "mail.example.com"

notes: "Free-form notes"         # Optional
```

## Fields

### engagement.name (required)
Human-readable name for the engagement. Displayed in output headers.

### engagement.id (optional)
Internal identifier for tracking. Not used in matching logic.

### engagement.start / engagement.end (optional)
Date boundaries in `YYYY-MM-DD` format. If set, the tool warns when validating outside this window.

### scope.include (required)
List of targets that are authorized. At least one entry is required. Entries can be:
- IPv4 addresses (`10.0.0.50`)
- IPv4 CIDR ranges (`192.168.1.0/24`)
- IPv6 addresses (`2001:db8::1`)
- IPv6 CIDR ranges (`2001:db8::/32`)
- Exact domains (`portal.example.com`)
- Wildcard domains (`*.example.com`)

### scope.exclude (optional)
List of targets explicitly excluded from scope. Same format as include entries. **Exclusions always take precedence over inclusions.**

### notes (optional)
Free-form text for operational notes (e.g., authorized testing windows, contact information).

## Matching Rules

1. A target is checked against exclusions first. If it matches any exclusion, the result is `excluded`.
2. If not excluded, the target is checked against inclusions. If it matches, the result is `in_scope`.
3. If neither excluded nor included, the result is `out_of_scope`.

## Wildcard Behavior

- `*.example.com` matches `sub.example.com`, `deep.sub.example.com`
- `*.example.com` does NOT match `example.com` (the bare domain)
- Matching is case-insensitive

## CIDR Behavior

- `192.168.1.0/24` matches all 256 addresses in that range
- An IP excluded individually will be excluded even if its CIDR is included
