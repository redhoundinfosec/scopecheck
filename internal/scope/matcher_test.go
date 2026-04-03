package scope

import (
	"testing"
)

func testScope() *Scope {
	return &Scope{
		Version: 1,
		Engagement: Engagement{
			Name: "Test Engagement",
			ID:   "TEST-001",
		},
		ScopeRules: ScopeRules{
			Include: []string{
				"192.168.1.0/24",
				"10.0.0.50",
				"2001:db8::/32",
				"*.example.com",
				"specific.target.com",
			},
			Exclude: []string{
				"192.168.1.1",
				"192.168.1.2",
				"mail.example.com",
			},
		},
	}
}

func TestMatcherIPv4InScope(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	tests := []struct {
		target string
		want   MatchStatus
	}{
		{"192.168.1.50", StatusInScope},
		{"192.168.1.100", StatusInScope},
		{"192.168.1.254", StatusInScope},
		{"10.0.0.50", StatusInScope},
	}

	for _, tt := range tests {
		result := m.Check(tt.target)
		if result.Status != tt.want {
			t.Errorf("Check(%q) = %s, want %s", tt.target, result.Status, tt.want)
		}
	}
}

func TestMatcherIPv4Excluded(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	tests := []struct {
		target    string
		want      MatchStatus
		matchedBy string
	}{
		{"192.168.1.1", StatusExcluded, "192.168.1.1"},
		{"192.168.1.2", StatusExcluded, "192.168.1.2"},
	}

	for _, tt := range tests {
		result := m.Check(tt.target)
		if result.Status != tt.want {
			t.Errorf("Check(%q) status = %s, want %s", tt.target, result.Status, tt.want)
		}
		if result.MatchedBy != tt.matchedBy {
			t.Errorf("Check(%q) matchedBy = %q, want %q", tt.target, result.MatchedBy, tt.matchedBy)
		}
	}
}

func TestMatcherIPv4OutOfScope(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	tests := []string{
		"10.20.30.40",
		"172.16.0.1",
		"8.8.8.8",
		"10.0.0.51",
	}

	for _, target := range tests {
		result := m.Check(target)
		if result.Status != StatusOutOfScope {
			t.Errorf("Check(%q) = %s, want out_of_scope", target, result.Status)
		}
	}
}

func TestMatcherIPv6(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	tests := []struct {
		target string
		want   MatchStatus
	}{
		{"2001:db8::1", StatusInScope},
		{"2001:db8:abcd::1234", StatusInScope},
		{"2001:db9::1", StatusOutOfScope},
		{"fe80::1", StatusOutOfScope},
	}

	for _, tt := range tests {
		result := m.Check(tt.target)
		if result.Status != tt.want {
			t.Errorf("Check(%q) = %s, want %s", tt.target, result.Status, tt.want)
		}
	}
}

func TestMatcherDomainExact(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	result := m.Check("specific.target.com")
	if result.Status != StatusInScope {
		t.Errorf("Check(specific.target.com) = %s, want in_scope", result.Status)
	}

	result = m.Check("other.domain.com")
	if result.Status != StatusOutOfScope {
		t.Errorf("Check(other.domain.com) = %s, want out_of_scope", result.Status)
	}
}

func TestMatcherDomainWildcard(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	tests := []struct {
		target string
		want   MatchStatus
	}{
		{"sub.example.com", StatusInScope},
		{"deep.sub.example.com", StatusInScope},
		{"app.example.com", StatusInScope},
		{"example.com", StatusOutOfScope},       // wildcard does NOT match bare domain
		{"notexample.com", StatusOutOfScope},     // no match
		{"mail.example.com", StatusExcluded},     // excluded takes precedence
	}

	for _, tt := range tests {
		result := m.Check(tt.target)
		if result.Status != tt.want {
			t.Errorf("Check(%q) = %s, want %s", tt.target, result.Status, tt.want)
		}
	}
}

func TestMatcherDomainCaseInsensitive(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	tests := []struct {
		target string
		want   MatchStatus
	}{
		{"SUB.EXAMPLE.COM", StatusInScope},
		{"Specific.Target.Com", StatusInScope},
		{"MAIL.EXAMPLE.COM", StatusExcluded},
	}

	for _, tt := range tests {
		result := m.Check(tt.target)
		if result.Status != tt.want {
			t.Errorf("Check(%q) = %s, want %s", tt.target, result.Status, tt.want)
		}
	}
}

func TestMatcherExclusionPrecedence(t *testing.T) {
	// Exclusion must always win, even if target matches an include
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	// 192.168.1.1 is in 192.168.1.0/24 (include) but also explicitly excluded
	result := m.Check("192.168.1.1")
	if result.Status != StatusExcluded {
		t.Errorf("Exclusion precedence: Check(192.168.1.1) = %s, want excluded", result.Status)
	}
}

func TestCheckAll(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	targets := []string{
		"192.168.1.50",
		"192.168.1.1",
		"10.20.30.40",
		"# this is a comment",
		"",
		"sub.example.com",
	}

	results := m.CheckAll(targets)
	if len(results) != 4 { // comments and empty lines skipped
		t.Errorf("CheckAll returned %d results, want 4", len(results))
	}
}

func TestSummarize(t *testing.T) {
	results := []MatchResult{
		{Status: StatusInScope},
		{Status: StatusInScope},
		{Status: StatusExcluded},
		{Status: StatusOutOfScope},
	}

	s := Summarize(results)
	if s.InScope != 2 || s.Excluded != 1 || s.OutOfScope != 1 || s.Total != 4 {
		t.Errorf("Summarize = %+v, want 2/1/1/4", s)
	}
}

func TestHasOutOfScope(t *testing.T) {
	allGood := []MatchResult{{Status: StatusInScope}, {Status: StatusInScope}}
	if HasOutOfScope(allGood) {
		t.Error("HasOutOfScope(allGood) = true, want false")
	}

	hasExcluded := []MatchResult{{Status: StatusInScope}, {Status: StatusExcluded}}
	if !HasOutOfScope(hasExcluded) {
		t.Error("HasOutOfScope(hasExcluded) = false, want true")
	}

	hasOOS := []MatchResult{{Status: StatusInScope}, {Status: StatusOutOfScope}}
	if !HasOutOfScope(hasOOS) {
		t.Error("HasOutOfScope(hasOOS) = false, want true")
	}
}

func TestMatcherEmptyTarget(t *testing.T) {
	m, err := NewMatcher(testScope())
	if err != nil {
		t.Fatalf("NewMatcher: %v", err)
	}

	result := m.Check("")
	if result.Status != StatusOutOfScope {
		t.Errorf("Check('') = %s, want out_of_scope", result.Status)
	}
}

func TestMatchStatusString(t *testing.T) {
	tests := []struct {
		s    MatchStatus
		want string
	}{
		{StatusInScope, "in_scope"},
		{StatusExcluded, "excluded"},
		{StatusOutOfScope, "out_of_scope"},
		{MatchStatus(99), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.s.String(); got != tt.want {
			t.Errorf("MatchStatus(%d).String() = %q, want %q", tt.s, got, tt.want)
		}
	}
}
