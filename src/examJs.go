package orson

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// PackageJSON represents the structure of package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Type            string            `json:"type"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type Violation struct {
	Finding
	Dependency string `json:"dependency"`
}

// this list could come directly from npmjs as projects tagged with mcp and/or modelcontextprotocol
var matchListJS = []string{
	"modelcontextprotocol",
	"mcp-server",
	"npm-helper-mcp",
	"mcp-utils",
	"mcp-framework",
	"mcp-proxy",
	"playwright/mcp",
	"composion/mcp",
	"context7-mcp",
	"n8n-nodes-mcp",
	"browsermcp/mcp",
	"-mcp", //might be too loose
}

func ExamJS(finding Finding) ([]Violation, error) {

	var Violations []Violation

	if finding.Dependency == "package.json" {
		// Read the package.json file
		data, err := os.ReadFile(finding.Path)
		if err != nil {
			return nil, fmt.Errorf("failed to read package.json: %s %w", finding.Path, err)
		}

		// Parse the JSON
		var pkg PackageJSON
		if err := json.Unmarshal(data, &pkg); err != nil {
			return nil, fmt.Errorf("failed to parse package.json: %s %w", finding.Path, err)
		}

		//match on name
		for _, match := range matchListJS {
			if strings.Contains(pkg.Name, match) {
				Violations = append(Violations, Violation{
					Finding:    finding,
					Dependency: pkg.Name,
				})
				break
			}
		}

		// Check dependencies
		for dep := range pkg.Dependencies {
			for _, match := range matchListJS {
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
			for _, match := range matchListJS {
				if strings.Contains(dep, match) {
					Violations = append(Violations, Violation{
						Finding:    finding,
						Dependency: dep,
					})
				}
			}
		}
	}
	return Violations, nil
}
