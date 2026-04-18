package scope

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFromFile reads and parses a scope YAML file.
func LoadFromFile(path string) (*Scope, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading scope file: %w", err)
	}
	return Parse(data)
}

// Parse parses scope YAML bytes into a Scope struct.
func Parse(data []byte) (*Scope, error) {
	var s Scope
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing scope YAML: %w", err)
	}
	if err := validate(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// validate checks the scope struct for required fields and consistency.
func validate(s *Scope) error {
	if s.Version != FormatVersion {
		return fmt.Errorf("unsupported scope version %d (expected %d)", s.Version, FormatVersion)
	}
	if s.Engagement.Name == "" {
		return fmt.Errorf("engagement name is required")
	}
	if len(s.ScopeRules.Include) == 0 {
		return fmt.Errorf("scope must have at least one include entry")
	}
	return nil
}

// TemplateYAML returns a template scope file as a string.
func TemplateYAML() string {
	return `# scopecheck scope definition
# Documentation: https://github.com/redhoundinfosec/scopecheck/blob/main/docs/scope-file-format.md
version: 1

engagement:
  name: "Example Engagement"
  id: "ENG-2026-001"
  start: "2026-01-15"
  end: "2026-02-15"

scope:
  include:
    # IPv4 addresses and CIDR ranges
    - "192.168.1.0/24"
    - "10.0.0.50"
    # IPv6 addresses and ranges
    # - "2001:db8::/32"
    # Domains (exact match)
    - "app.example.com"
    # Wildcard domains (matches subdomains)
    - "*.example.com"

  exclude:
    # Targets explicitly excluded from scope
    - "192.168.1.1"
    - "192.168.1.2"
    - "mail.example.com"

notes: "Authorized testing window: weekdays 22:00-06:00 EST only"
`
}
