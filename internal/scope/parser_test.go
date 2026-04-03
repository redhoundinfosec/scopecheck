package scope

import (
	"testing"
)

func TestParseValid(t *testing.T) {
	yaml := `
version: 1
engagement:
  name: "Test"
  id: "T-001"
scope:
  include:
    - "192.168.1.0/24"
    - "*.test.com"
  exclude:
    - "192.168.1.1"
`
	s, err := Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if s.Engagement.Name != "Test" {
		t.Errorf("Name = %q, want %q", s.Engagement.Name, "Test")
	}
	if len(s.ScopeRules.Include) != 2 {
		t.Errorf("Include count = %d, want 2", len(s.ScopeRules.Include))
	}
	if len(s.ScopeRules.Exclude) != 1 {
		t.Errorf("Exclude count = %d, want 1", len(s.ScopeRules.Exclude))
	}
}

func TestParseInvalidVersion(t *testing.T) {
	yaml := `
version: 99
engagement:
  name: "Test"
scope:
  include:
    - "192.168.1.0/24"
`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("Expected error for invalid version")
	}
}

func TestParseMissingName(t *testing.T) {
	yaml := `
version: 1
engagement:
  id: "T-001"
scope:
  include:
    - "192.168.1.0/24"
`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("Expected error for missing name")
	}
}

func TestParseMissingIncludes(t *testing.T) {
	yaml := `
version: 1
engagement:
  name: "Test"
scope:
  exclude:
    - "192.168.1.1"
`
	_, err := Parse([]byte(yaml))
	if err == nil {
		t.Fatal("Expected error for empty includes")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	_, err := Parse([]byte(`{{{invalid`))
	if err == nil {
		t.Fatal("Expected error for invalid YAML")
	}
}

func TestTemplateYAML(t *testing.T) {
	tmpl := TemplateYAML()
	if tmpl == "" {
		t.Fatal("TemplateYAML returned empty string")
	}
	// Verify the template is valid YAML
	s, err := Parse([]byte(tmpl))
	if err != nil {
		t.Fatalf("TemplateYAML produces invalid scope: %v", err)
	}
	if s.Version != 1 {
		t.Errorf("Template version = %d, want 1", s.Version)
	}
}
