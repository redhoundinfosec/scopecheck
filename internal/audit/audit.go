// Package audit provides append-only audit logging for scope checks.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultLogFile = ".scopecheck-audit.jsonl"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp  string `json:"timestamp"`
	Engagement string `json:"engagement"`
	Target     string `json:"target"`
	Status     string `json:"status"`
	MatchedBy  string `json:"matched_by,omitempty"`
	ScopeFile  string `json:"scope_file"`
}

// Logger writes audit entries to an append-only JSONL file.
type Logger struct {
	file    *os.File
	enabled bool
}

// NewLogger creates a new audit logger. If dir is empty, uses current directory.
func NewLogger(dir string, enabled bool) (*Logger, error) {
	if !enabled {
		return &Logger{enabled: false}, nil
	}

	if dir == "" {
		dir = "."
	}
	path := filepath.Join(dir, defaultLogFile)

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("opening audit log: %w", err)
	}
	return &Logger{file: f, enabled: true}, nil
}

// Log writes an audit entry.
func (l *Logger) Log(engagement, target, status, matchedBy, scopeFile string) error {
	if !l.enabled {
		return nil
	}

	entry := Entry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		Engagement: engagement,
		Target:     target,
		Status:     status,
		MatchedBy:  matchedBy,
		ScopeFile:  scopeFile,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshaling audit entry: %w", err)
	}

	if _, err := l.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("writing audit entry: %w", err)
	}
	return nil
}

// Close closes the audit log file.
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// ReadEntries reads the last N entries from the audit log.
// If n <= 0, returns all entries.
func ReadEntries(dir string, n int) ([]Entry, error) {
	if dir == "" {
		dir = "."
	}
	path := filepath.Join(dir, defaultLogFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No audit log yet
		}
		return nil, fmt.Errorf("reading audit log: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			continue // Skip malformed lines
		}
		entries = append(entries, e)
	}

	if n > 0 && len(entries) > n {
		entries = entries[len(entries)-n:]
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
