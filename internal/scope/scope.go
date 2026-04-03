// Package scope defines the core scope model and loading logic.
package scope

import (
	"fmt"
	"time"
)

// Version of the scope file format.
const FormatVersion = 1

// Scope represents a fully loaded and validated engagement scope.
type Scope struct {
	Version    int        `yaml:"version"`
	Engagement Engagement `yaml:"engagement"`
	ScopeRules ScopeRules `yaml:"scope"`
	Notes      string     `yaml:"notes,omitempty"`
}

// Engagement contains metadata about the authorized engagement.
type Engagement struct {
	Name  string `yaml:"name"`
	ID    string `yaml:"id"`
	Start string `yaml:"start,omitempty"`
	End   string `yaml:"end,omitempty"`
}

// ScopeRules contains the include and exclude target definitions.
type ScopeRules struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude,omitempty"`
}

// DateRange returns parsed start/end times. Returns zero values if not set.
func (e *Engagement) DateRange() (start, end time.Time, err error) {
	if e.Start != "" {
		start, err = time.Parse("2006-01-02", e.Start)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start date %q: %w", e.Start, err)
		}
	}
	if e.End != "" {
		end, err = time.Parse("2006-01-02", e.End)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end date %q: %w", e.End, err)
		}
	}
	return start, end, nil
}

// IsWithinWindow checks if the current time falls within the engagement window.
// Returns true if no dates are set (no restriction).
func (e *Engagement) IsWithinWindow(now time.Time) bool {
	start, end, err := e.DateRange()
	if err != nil {
		return true // Don't block on parse errors; warn separately
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	if !start.IsZero() && today.Before(start) {
		return false
	}
	if !end.IsZero() && today.After(end) {
		return false
	}
	return true
}
