package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redhoundinfosec/scopecheck/internal/audit"
)

func runAudit(gf globalFlags) int {
	last := 0 // 0 means all

	// Parse audit-specific flags
	for i := 0; i < len(gf.args); i++ {
		switch gf.args[i] {
		case "--last", "-n":
			if i+1 < len(gf.args) {
				n, err := strconv.Atoi(gf.args[i+1])
				if err != nil || n < 1 {
					fmt.Fprintf(os.Stderr, "Error: --last requires a positive integer\n")
					return ExitError
				}
				last = n
				i++
			}
		case "--clear":
			if err := os.Remove(".scopecheck-audit.jsonl"); err != nil && !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error clearing audit log: %v\n", err)
				return ExitError
			}
			fmt.Fprintln(os.Stdout, "Audit log cleared.")
			return ExitOK
		}
	}

	entries, err := audit.ReadEntries("", last)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	if len(entries) == 0 {
		fmt.Fprintln(os.Stdout, "No audit entries found.")
		return ExitOK
	}

	if gf.format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return exitIf(enc.Encode(entries))
	}

	// Text format
	for _, e := range entries {
		status := e.Status
		matchInfo := ""
		if e.MatchedBy != "" {
			matchInfo = fmt.Sprintf(" (matched: %s)", e.MatchedBy)
		}
		fmt.Fprintf(os.Stdout, "  %s  %-14s %-40s %s%s\n",
			e.Timestamp, status, e.Target, e.Engagement, matchInfo)
	}

	fmt.Fprintf(os.Stdout, "\n  %d entries shown\n", len(entries))
	return ExitOK
}

func exitIf(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}
	return ExitOK
}
