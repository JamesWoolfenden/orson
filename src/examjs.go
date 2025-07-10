package orson

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// PackageJSON represents the structure of package.json
type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Violation struct {
	Finding
	Dependency string `json:"dependency"`
}

var matchList = []string{
	"modelcontextprotocol",
}

func ExamJS(finding Finding) ([]Violation, error) {

	var Violations []Violation

	// Read the package.json file
	data, err := os.ReadFile(finding.Path)
	if err != nil {
		return 0, fmt.Errorf("failed to read package.json: %w", err)
	}

	// Parse the JSON
	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return 0, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Check dependencies
	for dep := range pkg.Dependencies {
		for _, match := range matchList {
			if strings.Contains(dep, match) {
				Violations = append(Violations, Violation{
					Finding:    finding,
					Dependency: dep,
				})
			}
		}
	}

	// Check devDependencies
	for dep := range pkg.DevDependencies {
		for _, match := range matchList {
			if strings.Contains(dep, match) {
				Violations = append(Violations, Violation{
					Finding:    finding,
					Dependency: dep,
				})
			}
		}
	}

	return Violations, nil
}
