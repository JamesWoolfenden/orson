package orson

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/mod/modfile"
)

func ExamGo(finding Finding) error {
	// Read the go.mod file
	if finding.Dependency == "go.mod" {

		data, err := os.ReadFile(finding.Path)
		if err != nil {
			return fmt.Errorf("failed to read go.mod: %w", err)
		}

		// Parse the file
		f, err := modfile.Parse(finding.Path, data, nil)
		if err != nil {
			return fmt.Errorf("failed to parse go.mod: %w", err)
		}

		if strings.Contains(f.Module.Mod.Path, "mcp-go") {
			log.Info().Msgf("Go: %s\n", f.Module.Mod.Path)
		}

		for _, r := range f.Require {
			if strings.Contains(r.Mod.Path, "mcp-go") {
				log.Info().Msgf("Go: %s\n", r.Mod.Path)
			}
		}
	}

	return nil
}
