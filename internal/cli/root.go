// Package cli implements the command-line interface.
package cli

import (
	"fmt"
	"os"

	"github.com/redhoundinfosec/scopecheck/internal/output"
)

const (
	Version   = "0.1.0"
	ExitOK    = 0
	ExitFail  = 1
	ExitError = 2
)

// Run is the main entry point for the CLI.
func Run(args []string) int {
	if len(args) < 1 {
		printUsage()
		return ExitError
	}

	// Parse global flags
	gf := parseGlobalFlags(args)

	switch gf.command {
	case "init":
		return runInit(gf)
	case "validate":
		return runValidate(gf)
	case "audit":
		return runAudit(gf)
	case "version":
		fmt.Fprintf(os.Stdout, "scopecheck v%s\n", Version)
		return ExitOK
	case "help", "--help", "-h":
		printUsage()
		return ExitOK
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", gf.command)
		printUsage()
		return ExitError
	}
}

type globalFlags struct {
	command   string
	scopeFile string
	format    string
	quiet     bool
	verbose   bool
	noColor   bool
	noAudit   bool
	args      []string // remaining args after flag parsing
}

func parseGlobalFlags(args []string) globalFlags {
	gf := globalFlags{
		scopeFile: "scope.yaml",
		format:    "text",
	}

	i := 0
	for i < len(args) {
		switch args[i] {
		case "-s", "--scope":
			if i+1 < len(args) {
				gf.scopeFile = args[i+1]
				i += 2
				continue
			}
			i++
		case "-f", "--format":
			if i+1 < len(args) {
				gf.format = args[i+1]
				i += 2
				continue
			}
			i++
		case "-q", "--quiet":
			gf.quiet = true
			i++
		case "-v", "--verbose":
			gf.verbose = true
			i++
		case "--no-color":
			gf.noColor = true
			i++
		case "--no-audit":
			gf.noAudit = true
			i++
		default:
			if gf.command == "" {
				gf.command = args[i]
				i++
			} else {
				gf.args = args[i:]
				return gf
			}
		}
	}
	return gf
}

func printUsage() {
	banner := output.FormatBanner()
	fmt.Fprintf(os.Stderr, `%s
  v%s — Validate targets against authorized engagement scope

USAGE
  scopecheck <command> [flags] [arguments]

COMMANDS
  init          Generate a template scope file
  validate      Check targets against scope definition
  audit         View scope check audit log
  version       Print version information
  help          Show this help

GLOBAL FLAGS
  -s, --scope <file>    Path to scope file (default: scope.yaml)
  -f, --format <fmt>    Output format: text, json, csv (default: text)
  -q, --quiet           Suppress non-essential output
  -v, --verbose         Show detailed match information
      --no-color        Disable colored output
      --no-audit        Disable audit logging
  -h, --help            Show help

EXAMPLES
  scopecheck init
  scopecheck validate 192.168.1.50
  scopecheck validate --file targets.txt
  scopecheck validate --file targets.txt --format json
  cat targets.txt | scopecheck validate --stdin
  scopecheck audit --last 20

EXIT CODES
  0  All targets are in scope
  1  One or more targets are out of scope or excluded
  2  Error (invalid input, missing files, etc.)

AUTHORIZED USE ONLY
  This tool is designed to support legitimate, authorized security testing.
  Always ensure you have written authorization before conducting any assessment.

`, banner, Version)
}
