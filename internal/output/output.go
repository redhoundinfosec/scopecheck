// Package output handles formatting of scope check results.
package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/redhoundinfosec/scopecheck/internal/scope"
)

const (
	FormatText = "text"
	FormatJSON = "json"
	FormatCSV  = "csv"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

// Config controls output behavior.
type Config struct {
	Format  string
	NoColor bool
	Quiet   bool
	Verbose bool
	Writer  io.Writer
}

// Render outputs results in the specified format.
func Render(cfg Config, engagement string, results []scope.MatchResult, summary scope.Summary) error {
	switch cfg.Format {
	case FormatJSON:
		return renderJSON(cfg.Writer, engagement, results, summary)
	case FormatCSV:
		return renderCSV(cfg.Writer, results)
	default:
		return renderText(cfg, engagement, results, summary)
	}
}

func renderText(cfg Config, engagement string, results []scope.MatchResult, summary scope.Summary) error {
	w := cfg.Writer

	if !cfg.Quiet {
		header := fmt.Sprintf("scopecheck — %s", engagement)
		fmt.Fprintln(w)
		if !cfg.NoColor {
			fmt.Fprintf(w, "  %s%s%s\n\n", colorBold, header, colorReset)
		} else {
			fmt.Fprintf(w, "  %s\n\n", header)
		}
	}

	for _, r := range results {
		icon, label, color := statusDisplay(r.Status, cfg.NoColor)
		matchInfo := ""
		if r.MatchedBy != "" && cfg.Verbose {
			matchInfo = fmt.Sprintf("  (matches %s)", r.MatchedBy)
		}

		if cfg.NoColor {
			fmt.Fprintf(w, "  %s %-40s %s%s\n", icon, r.Target, label, matchInfo)
		} else {
			fmt.Fprintf(w, "  %s %s%-40s%s %s%s%s%s%s\n",
				icon, colorReset, r.Target, colorReset,
				color, label, colorReset,
				colorDim, matchInfo+colorReset)
		}
	}

	if !cfg.Quiet {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "  Summary: %d in scope, %d excluded, %d out of scope (%d total)\n",
			summary.InScope, summary.Excluded, summary.OutOfScope, summary.Total)
		fmt.Fprintln(w)
	}

	return nil
}

func statusDisplay(status scope.MatchStatus, noColor bool) (icon, label, color string) {
	switch status {
	case scope.StatusInScope:
		if noColor {
			return "[+]", "IN SCOPE", ""
		}
		return colorGreen + "✓" + colorReset, "IN SCOPE", colorGreen
	case scope.StatusExcluded:
		if noColor {
			return "[-]", "EXCLUDED", ""
		}
		return colorYellow + "✗" + colorReset, "EXCLUDED", colorYellow
	case scope.StatusOutOfScope:
		if noColor {
			return "[!]", "OUT OF SCOPE", ""
		}
		return colorRed + "✗" + colorReset, "OUT OF SCOPE", colorRed
	default:
		return "?", "UNKNOWN", ""
	}
}

type jsonOutput struct {
	Engagement string              `json:"engagement"`
	Results    []scope.MatchResult `json:"results"`
	Summary    scope.Summary       `json:"summary"`
}

func renderJSON(w io.Writer, engagement string, results []scope.MatchResult, summary scope.Summary) error {
	out := jsonOutput{
		Engagement: engagement,
		Results:    results,
		Summary:    summary,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func renderCSV(w io.Writer, results []scope.MatchResult) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"target", "status", "matched_by"}); err != nil {
		return err
	}
	for _, r := range results {
		record := []string{r.Target, r.StatusStr, r.MatchedBy}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return nil
}

// RenderAudit formats audit entries for display.
func RenderAudit(w io.Writer, entries []interface{ Format() string }, noColor bool) {
	// This is handled directly in the CLI command
}

// FormatBanner returns the ASCII banner.
func FormatBanner() string {
	return strings.TrimSpace(`
                                 _               _    
  ___  ___ ___  _ __   ___  ___| |__   ___  ___| | __
 / __|/ __/ _ \| '_ \ / _ \/ __| '_ \ / _ \/ __| |/ /
 \__ \ (_| (_) | |_) |  __/ (__| | | |  __/ (__|   < 
 |___/\___\___/| .__/ \___|\___|_| |_|\___|\___|_|\_\
               |_|
`)
}
