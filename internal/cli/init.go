package cli

import (
	"fmt"
	"os"

	"github.com/redhoundinfosec/scopecheck/internal/scope"
)

func runInit(gf globalFlags) int {
	path := gf.scopeFile

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		fmt.Fprintf(os.Stderr, "Error: %s already exists. Remove it first or use a different path with --scope.\n", path)
		return ExitError
	}

	if err := os.WriteFile(path, []byte(scope.TemplateYAML()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing scope file: %v\n", err)
		return ExitError
	}

	if !gf.quiet {
		fmt.Fprintf(os.Stdout, "Created scope template: %s\n", path)
		fmt.Fprintf(os.Stdout, "Edit this file with your engagement details and authorized targets.\n")
	}
	return ExitOK
}
