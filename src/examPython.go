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
	"-mcp", //might be too lose
	"sse",  //might be too lose
	"anthrophic",
	"mcp[cli]",
}

type pyproject struct {
	Project map[string]project
	//[build-system]
	//	requires = ["setuptools>=61.0"]
	//	build-backend = "setuptools.build_meta"
	//
	//[project]
	//	name='aisbtokenizers'
	//	version='0.0.1'
	//	description='Tokenizers library'
	//	authors= [{ name="James Woolfenden", email='jwoolfenden@paloaltonetworks.com'},]
	//	readme = "README.md"
	//	requires-python = ">=3.10"
}

type project struct {
	Name string
	//version='0.0.1'
	//description='Tokenizers library'
	//authors= [{ name="James Woolfenden", email='jwoolfenden@paloaltonetworks.com'},]
	//readme = "README.md"
	//requires-python = ">=3.10"
}

type PythonRequirement struct {
	Name    string
	Version string
}

func ExamPython(finding Finding) ([]Violation, error) {

	var violations []Violation

	switch finding.Dependency {
	case "requirements.txt":
		violations, err := parseRequirementsTxt(finding)
		if err != nil {
			return violations, err
		}
	case "pyproject.toml":
		violations, err := parsePyProject(finding)
		if err != nil {
			return violations, err
		}
	default:

	}
	if finding.Dependency != "requirements.txt" {
		return violations, nil
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
		// Join the rest as version specification
		version = strings.TrimSpace(strings.Join(parts[1:], ""))
	}

	return PythonRequirement{
		Name:    name,
		Version: version,
	}
}

func parsePyProject(finding Finding) ([]Violation, error) {
	var violations []Violation
	fmt.Println(finding.Path)

	var config pyproject
	Pyproject, err := toml.DecodeFile(finding.Path, &config)

	if err != nil {
		return violations, err
	}

	fmt.Println(Pyproject)

	//file, err := os.Open(finding.Path)
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()
	//
	return violations, nil
}
