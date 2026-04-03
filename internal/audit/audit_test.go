package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoggerAndRead(t *testing.T) {
	dir := t.TempDir()

	// Create logger
	logger, err := NewLogger(dir, true)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}

	// Write entries
	if err := logger.Log("Test Engagement", "192.168.1.1", "in_scope", "192.168.1.0/24", "scope.yaml"); err != nil {
		t.Fatalf("Log: %v", err)
	}
	if err := logger.Log("Test Engagement", "10.0.0.1", "out_of_scope", "", "scope.yaml"); err != nil {
		t.Fatalf("Log: %v", err)
	}
	logger.Close()

	// Read all entries
	entries, err := ReadEntries(dir, 0)
	if err != nil {
		t.Fatalf("ReadEntries: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("ReadEntries count = %d, want 2", len(entries))
	}
	if entries[0].Target != "192.168.1.1" {
		t.Errorf("Entry[0].Target = %q, want %q", entries[0].Target, "192.168.1.1")
	}
	if entries[1].Status != "out_of_scope" {
		t.Errorf("Entry[1].Status = %q, want %q", entries[1].Status, "out_of_scope")
	}

	// Read last 1
	last, err := ReadEntries(dir, 1)
	if err != nil {
		t.Fatalf("ReadEntries(last 1): %v", err)
	}
	if len(last) != 1 {
		t.Fatalf("ReadEntries(last 1) count = %d, want 1", len(last))
	}
	if last[0].Target != "10.0.0.1" {
		t.Errorf("Last entry target = %q, want %q", last[0].Target, "10.0.0.1")
	}
}

func TestLoggerDisabled(t *testing.T) {
	logger, err := NewLogger("", false)
	if err != nil {
		t.Fatalf("NewLogger(disabled): %v", err)
	}
	// Should silently succeed
	if err := logger.Log("Test", "target", "in_scope", "", "scope.yaml"); err != nil {
		t.Fatalf("Log(disabled): %v", err)
	}
	logger.Close()
}

func TestReadEntriesNoFile(t *testing.T) {
	dir := t.TempDir()
	entries, err := ReadEntries(dir, 0)
	if err != nil {
		t.Fatalf("ReadEntries(no file): %v", err)
	}
	if entries != nil {
		t.Errorf("Expected nil entries, got %d", len(entries))
	}
}

func TestAuditFilePermissions(t *testing.T) {
	dir := t.TempDir()
	logger, err := NewLogger(dir, true)
	if err != nil {
		t.Fatalf("NewLogger: %v", err)
	}
	logger.Log("Test", "target", "in_scope", "", "scope.yaml")
	logger.Close()

	path := filepath.Join(dir, defaultLogFile)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("Audit file permissions = %o, want 0600", perm)
	}
}
