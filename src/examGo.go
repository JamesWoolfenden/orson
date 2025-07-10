package orson

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
)

var matchListGo = []string{
	"modelcontextprotocol",
	"mcp-server",
	"mcp-golang",
	"mcp-go",
	"-mcp", //might be too lose
	"gin-contrib/sse",
	"sse", //might be too lose
}

func ExamGo(finding Finding) ([]Violation, error) {
	var Violations []Violation

	// Read the go.mod file
	if finding.Dependency == "go.mod" {

		data, err := os.ReadFile(finding.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read go.mod: %w", err)
		}

		// Parse the file
		f, err := modfile.Parse(finding.Path, data, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to parse go.mod: %w", err)
		}

		Violations = findInModuleName(finding, f, Violations)

		Violations = findInModuleDependencies(finding, f, Violations)
	}

	return Violations, nil
}

func findInModuleDependencies(finding Finding, f *modfile.File, Violations []Violation) []Violation {
	//look at packages
	for _, r := range f.Require {
		if len(r.Syntax.Token) == 0 {
			log.Info().Msgf("Empty: %s\n", r.Syntax.Token)
			continue
		}

		dep := r.Syntax.Token[0]
		for _, match := range matchListGo {
			if strings.Contains(dep, match) {
				Violations = append(Violations, Violation{
					Finding:    finding,
					Dependency: dep,
				})
				break
			}
		}
	}
	return Violations
}

func findInModuleName(finding Finding, f *modfile.File, Violations []Violation) []Violation {
	//look at module name
	for _, match := range matchListGo {
		if strings.Contains(f.Module.Mod.Path, match) {
			Violations = append(Violations, Violation{
				Finding:    finding,
				Dependency: match,
			})
			//we only need to find this once
			return Violations
		}
	}
	return Violations
}
