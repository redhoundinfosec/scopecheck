package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/redhoundinfosec/scopecheck/internal/audit"
	"github.com/redhoundinfosec/scopecheck/internal/output"
	"github.com/redhoundinfosec/scopecheck/internal/scope"
)

func runValidate(gf globalFlags) int {
	// Load scope file
	s, err := scope.LoadFromFile(gf.scopeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	// Warn if outside engagement window
	if !s.Engagement.IsWithinWindow(time.Now()) {
		fmt.Fprintf(os.Stderr, "WARNING: Current date is outside the engagement window (%s to %s)\n",
			s.Engagement.Start, s.Engagement.End)
	}

	// Build matcher
	matcher, err := scope.NewMatcher(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building matcher: %v\n", err)
		return ExitError
	}

	// Collect targets
	targets, err := collectTargets(gf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitError
	}

	if len(targets) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no targets provided\n")
		fmt.Fprintf(os.Stderr, "Usage: scopecheck validate <target> | --file <file> | --stdin\n")
		return ExitError
	}

	// Validate
	results := matcher.CheckAll(targets)
	summary := scope.Summarize(results)

	// Output
	cfg := output.Config{
		Format:  gf.format,
		NoColor: gf.noColor,
		Quiet:   gf.quiet,
		Verbose: gf.verbose,
		Writer:  os.Stdout,
	}
	if err := output.Render(cfg, s.Engagement.Name, results, summary); err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering output: %v\n", err)
		return ExitError
	}

	// Audit log
	if !gf.noAudit {
		logger, err := audit.NewLogger("", true)
		if err != nil {
			// Non-fatal: warn but continue
			fmt.Fprintf(os.Stderr, "Warning: could not open audit log: %v\n", err)
		} else {
			defer logger.Close()
			for _, r := range results {
				if err := logger.Log(s.Engagement.Name, r.Target, r.StatusStr, r.MatchedBy, gf.scopeFile); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not write audit log entry for %q: %v\n", r.Target, err)
				}
			}
		}
	}

	// Exit code
	if scope.HasOutOfScope(results) {
		return ExitFail
	}
	return ExitOK
}

func collectTargets(gf globalFlags) ([]string, error) {
	var targets []string
	fromFile := ""
	fromStdin := false

	// Parse validate-specific flags
	remaining := gf.args
	for i := 0; i < len(remaining); i++ {
		switch remaining[i] {
		case "--file":
			if i+1 < len(remaining) {
				fromFile = remaining[i+1]
				i++
			} else {
				return nil, fmt.Errorf("--file requires a path argument")
			}
		case "--stdin":
			fromStdin = true
		default:
			// Bare argument = inline target
			targets = append(targets, remaining[i])
		}
	}

	// Read from file
	if fromFile != "" {
		lines, err := readTargetFile(fromFile)
		if err != nil {
			return nil, fmt.Errorf("reading targets from file: %w", err)
		}
		targets = append(targets, lines...)
	}

	// Read from stdin
	if fromStdin {
		lines, err := readStdin()
		if err != nil {
			return nil, fmt.Errorf("reading targets from stdin: %w", err)
		}
		targets = append(targets, lines...)
	}

	// Dedupe while preserving order (e.g., when combining --file + inline targets)
	targets = dedupeTargets(targets)

	// Limit target count to reduce accidental misuse and resource exhaustion
	if err := validateTargetCount(targets); err != nil {
		return nil, err
	}

	return targets, nil
}

func readTargetFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening target file: %w", err)
	}
	defer f.Close()

	lines, err := scanLines(f)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func readStdin() ([]string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, fmt.Errorf("reading stdin metadata: %w", err)
	}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("--stdin specified but no data piped to stdin")
	}
	lines, err := scanLines(os.Stdin)
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func scanLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	return lines, nil
}

func dedupeTargets(targets []string) []string {
	out := make([]string, 0, len(targets))
	seen := make(map[string]struct{}, len(targets))
	for _, t := range targets {
		if _, ok := seen[t]; ok {
			continue
		}
		seen[t] = struct{}{}
		out = append(out, t)
	}
	return out
}

func validateTargetCount(targets []string) error {
	// Keep this intentionally conservative to avoid accidental large-batch validation from huge stdin input.
	const maxTargets = 10000
	if len(targets) > maxTargets {
		return fmt.Errorf("too many targets (%d); split input into smaller batches", len(targets))
	}
	return nil
}
