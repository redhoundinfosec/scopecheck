package scope

import (
	"testing"
	"time"
)

func TestEngagementDateRange(t *testing.T) {
	e := Engagement{Start: "2026-01-15", End: "2026-02-15"}
	start, end, err := e.DateRange()
	if err != nil {
		t.Fatalf("DateRange: %v", err)
	}
	if start.Year() != 2026 || start.Month() != 1 || start.Day() != 15 {
		t.Errorf("Start = %v, want 2026-01-15", start)
	}
	if end.Year() != 2026 || end.Month() != 2 || end.Day() != 15 {
		t.Errorf("End = %v, want 2026-02-15", end)
	}
}

func TestEngagementDateRangeEmpty(t *testing.T) {
	e := Engagement{}
	start, end, err := e.DateRange()
	if err != nil {
		t.Fatalf("DateRange: %v", err)
	}
	if !start.IsZero() || !end.IsZero() {
		t.Error("Expected zero times for empty dates")
	}
}

func TestEngagementDateRangeInvalid(t *testing.T) {
	e := Engagement{Start: "not-a-date"}
	_, _, err := e.DateRange()
	if err == nil {
		t.Fatal("Expected error for invalid date")
	}
}

func TestIsWithinWindow(t *testing.T) {
	e := Engagement{Start: "2026-01-15", End: "2026-02-15"}

	// Within window
	within := time.Date(2026, 1, 20, 12, 0, 0, 0, time.UTC)
	if !e.IsWithinWindow(within) {
		t.Error("IsWithinWindow(within) = false, want true")
	}

	// Before window
	before := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	if e.IsWithinWindow(before) {
		t.Error("IsWithinWindow(before) = true, want false")
	}

	// After window
	after := time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
	if e.IsWithinWindow(after) {
		t.Error("IsWithinWindow(after) = true, want false")
	}

	// On start date
	onStart := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	if !e.IsWithinWindow(onStart) {
		t.Error("IsWithinWindow(onStart) = false, want true")
	}

	// On end date
	onEnd := time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)
	if !e.IsWithinWindow(onEnd) {
		t.Error("IsWithinWindow(onEnd) = false, want true")
	}

	// No dates set = always within window
	noWindow := Engagement{}
	if !noWindow.IsWithinWindow(before) {
		t.Error("IsWithinWindow(no dates) = false, want true")
	}
}
