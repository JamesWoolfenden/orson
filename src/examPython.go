package orson

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var matchListPython = []string{
	"modelcontextprotocol",
	"mcp-server",
	"-mcp", //might be too loose
	"sse",  //might be too loose
	"anthrophic",
	"mcp[cli]",
	"fastmcp",
	"mcp", //might be too loose
}

// PyProject Updated struct to properly parse pyproject.toml
type PyProject struct {
	Project     Project     `toml:"project"`
	BuildSystem BuildSystem `toml:"build-system"`
}

type Project struct {
	Name           string              `toml:"name"`
	Version        string              `toml:"version"`
	Description    string              `toml:"description"`
	Authors        []Author            `toml:"authors"`
	Readme         string              `toml:"readme"`
	RequiresPython string              `toml:"requires-python"`
	Dependencies   []string            `toml:"dependencies"`
	OptionalDeps   map[string][]string `toml:"optional-dependencies"`
}

type Author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type BuildSystem struct {
	Requires     []string `toml:"requires"`
	BuildBackend string   `toml:"build-backend"`
}

type PythonRequirement struct {
	Name    string
	Version string
}

func ExamPython(finding Finding) ([]Violation, error) {

	var violations []Violation
	var err error

	switch finding.Dependency {
	case "requirements.txt":
		violations, err = parseRequirementsTxt(finding)
		if err != nil {
			return violations, err
		}
	case "pyproject.toml":
		violations, err = parsePyProject(finding)
		if err != nil {
			return violations, err
		}
	default:

	}
	return violations, nil
}

func parseRequirementsTxt(finding Finding) ([]Violation, error) {
	var violations []Violation

	requirements, err := parseRequirements(finding.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse requirements.txt: %w", err)
	}

	for _, req := range requirements {
		for _, match := range matchListPython {
			if strings.Contains(req.Name, match) {
				violations = append(violations, Violation{
					Finding:    finding,
					Dependency: req.Name,
				})
				break
			}
		}
	}
	return violations, nil
}

func parseRequirements(filepath string) ([]PythonRequirement, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var requirements []PythonRequirement
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the requirement line
		req := parseLine(line)
		requirements = append(requirements, req)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return requirements, nil
}

func parseLine(line string) PythonRequirement {
	// Split on common version specifiers
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == '=' || r == '>' || r == '<' || r == '~' || r == '!'
	})

	name := strings.TrimSpace(parts[0])
	version := ""

	if len(parts) > 1 {
		// Join the rest as a version specification
		version = strings.TrimSpace(strings.Join(parts[1:], ""))
	}

	return PythonRequirement{
		Name:    name,
		Version: version,
	}
}

func parsePyProject(finding Finding) ([]Violation, error) {
	var violations []Violation
	var config PyProject

	_, err := toml.DecodeFile(finding.Path, &config)
	if err != nil {
		return violations, fmt.Errorf("failed to parse pyproject.toml: %w", err)
	}

	//if it's the actual project, that's an MCP server, and we only need this once
	for _, match := range matchListPython {
		if strings.Contains(strings.ToLower(config.Project.Name), strings.ToLower(match)) {
			violations = append(violations, Violation{
				Finding:    finding,
				Dependency: config.Project.Name,
			})
			break
		}
	}

	// Check dependencies in project.dependencies
	for _, dep := range config.Project.Dependencies {
		req := parseLine(dep)
		for _, match := range matchListPython {
			if strings.Contains(strings.ToLower(req.Name), strings.ToLower(match)) {
				violations = append(violations, Violation{
					Finding:    finding,
					Dependency: req.Name,
				})
				break
			}
		}
	}

	// Check optional dependencies
	for group, deps := range config.Project.OptionalDeps {
		for _, dep := range deps {
			req := parseLine(dep)
			for _, match := range matchListPython {
				if strings.Contains(strings.ToLower(req.Name), strings.ToLower(match)) {
					violations = append(violations, Violation{
						Finding:    finding,
						Dependency: fmt.Sprintf("%s (optional: %s)", req.Name, group),
					})
					break
				}
			}
		}
	}

	// Check build system requirements
	for _, dep := range config.BuildSystem.Requires {
		req := parseLine(dep)
		for _, match := range matchListPython {
			if strings.Contains(strings.ToLower(req.Name), strings.ToLower(match)) {
				violations = append(violations, Violation{
					Finding:    finding,
					Dependency: fmt.Sprintf("%s (build)", req.Name),
				})
				break
			}
		}
	}

	return violations, nil
}
