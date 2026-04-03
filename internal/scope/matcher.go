package scope

import (
	"net"
	"strings"
)

// MatchStatus represents the result of a scope check.
type MatchStatus int

const (
	StatusInScope    MatchStatus = iota // Target is within authorized scope
	StatusExcluded                      // Target matches an exclusion rule
	StatusOutOfScope                    // Target is not in scope
)

func (s MatchStatus) String() string {
	switch s {
	case StatusInScope:
		return "in_scope"
	case StatusExcluded:
		return "excluded"
	case StatusOutOfScope:
		return "out_of_scope"
	default:
		return "unknown"
	}
}

// MatchResult represents the outcome of checking a single target.
type MatchResult struct {
	Target    string      `json:"target"`
	Status    MatchStatus `json:"-"`
	StatusStr string      `json:"status"`
	MatchedBy string      `json:"matched_by"`
}

// Matcher performs scope validation checks.
type Matcher struct {
	includeNets    []*net.IPNet
	includeIPs     []net.IP
	includeDomains []string // exact domains
	includeWilds   []string // wildcard patterns like "*.example.com"

	excludeNets    []*net.IPNet
	excludeIPs     []net.IP
	excludeDomains []string
	excludeWilds   []string
}

// NewMatcher creates a Matcher from a Scope definition.
func NewMatcher(s *Scope) (*Matcher, error) {
	m := &Matcher{}

	for _, entry := range s.ScopeRules.Include {
		if err := m.addEntry(entry, true); err != nil {
			return nil, err
		}
	}
	for _, entry := range s.ScopeRules.Exclude {
		if err := m.addEntry(entry, false); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func (m *Matcher) addEntry(entry string, include bool) error {
	entry = strings.TrimSpace(entry)
	if entry == "" {
		return nil
	}

	// Try CIDR first
	_, ipNet, err := net.ParseCIDR(entry)
	if err == nil {
		if include {
			m.includeNets = append(m.includeNets, ipNet)
		} else {
			m.excludeNets = append(m.excludeNets, ipNet)
		}
		return nil
	}

	// Try plain IP
	ip := net.ParseIP(entry)
	if ip != nil {
		if include {
			m.includeIPs = append(m.includeIPs, ip)
		} else {
			m.excludeIPs = append(m.excludeIPs, ip)
		}
		return nil
	}

	// Must be a domain or wildcard
	lower := strings.ToLower(entry)
	if strings.HasPrefix(lower, "*.") {
		if include {
			m.includeWilds = append(m.includeWilds, lower)
		} else {
			m.excludeWilds = append(m.excludeWilds, lower)
		}
	} else {
		if include {
			m.includeDomains = append(m.includeDomains, lower)
		} else {
			m.excludeDomains = append(m.excludeDomains, lower)
		}
	}
	return nil
}

// Check validates a single target against the scope.
// Exclusions always take precedence over inclusions.
func (m *Matcher) Check(target string) MatchResult {
	target = strings.TrimSpace(target)
	result := MatchResult{Target: target}

	// Determine if target is IP or domain
	ip := net.ParseIP(target)

	if ip != nil {
		result = m.checkIP(ip, target)
	} else {
		result = m.checkDomain(strings.ToLower(target))
	}

	result.StatusStr = result.Status.String()
	return result
}

func (m *Matcher) checkIP(ip net.IP, raw string) MatchResult {
	// Check exclusions first (they take precedence)
	for _, exIP := range m.excludeIPs {
		if exIP.Equal(ip) {
			return MatchResult{Target: raw, Status: StatusExcluded, MatchedBy: exIP.String()}
		}
	}
	for _, exNet := range m.excludeNets {
		if exNet.Contains(ip) {
			return MatchResult{Target: raw, Status: StatusExcluded, MatchedBy: exNet.String()}
		}
	}

	// Check inclusions
	for _, inIP := range m.includeIPs {
		if inIP.Equal(ip) {
			return MatchResult{Target: raw, Status: StatusInScope, MatchedBy: inIP.String()}
		}
	}
	for _, inNet := range m.includeNets {
		if inNet.Contains(ip) {
			return MatchResult{Target: raw, Status: StatusInScope, MatchedBy: inNet.String()}
		}
	}

	return MatchResult{Target: raw, Status: StatusOutOfScope}
}

func (m *Matcher) checkDomain(domain string) MatchResult {
	// Check exclusions first
	for _, exDom := range m.excludeDomains {
		if domain == exDom {
			return MatchResult{Target: domain, Status: StatusExcluded, MatchedBy: exDom}
		}
	}
	for _, exWild := range m.excludeWilds {
		if matchWildcard(domain, exWild) {
			return MatchResult{Target: domain, Status: StatusExcluded, MatchedBy: exWild}
		}
	}

	// Check inclusions
	for _, inDom := range m.includeDomains {
		if domain == inDom {
			return MatchResult{Target: domain, Status: StatusInScope, MatchedBy: inDom}
		}
	}
	for _, inWild := range m.includeWilds {
		if matchWildcard(domain, inWild) {
			return MatchResult{Target: domain, Status: StatusInScope, MatchedBy: inWild}
		}
	}

	return MatchResult{Target: domain, Status: StatusOutOfScope}
}

// matchWildcard checks if a domain matches a wildcard pattern like "*.example.com".
// "*.example.com" matches "sub.example.com" but NOT "example.com" itself.
func matchWildcard(domain, pattern string) bool {
	// pattern is "*.example.com", suffix is ".example.com"
	suffix := pattern[1:] // strip the "*"
	return strings.HasSuffix(domain, suffix) && domain != suffix[1:]
}

// CheckAll validates multiple targets and returns all results.
func (m *Matcher) CheckAll(targets []string) []MatchResult {
	results := make([]MatchResult, 0, len(targets))
	for _, t := range targets {
		t = strings.TrimSpace(t)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		results = append(results, m.Check(t))
	}
	return results
}

// Summary counts the results by status.
type Summary struct {
	InScope    int `json:"in_scope"`
	Excluded   int `json:"excluded"`
	OutOfScope int `json:"out_of_scope"`
	Total      int `json:"total"`
}

// Summarize produces a Summary from a slice of results.
func Summarize(results []MatchResult) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusInScope:
			s.InScope++
		case StatusExcluded:
			s.Excluded++
		case StatusOutOfScope:
			s.OutOfScope++
		}
	}
	return s
}

// HasOutOfScope returns true if any result is out of scope or excluded.
func HasOutOfScope(results []MatchResult) bool {
	for _, r := range results {
		if r.Status == StatusOutOfScope || r.Status == StatusExcluded {
			return true
		}
	}
	return false
}
