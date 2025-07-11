package orson

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func GetViolations(root string) {
	ignoreDirs := map[string]bool{
		".git":              true,
		".terraform":        true,
		".terragrunt-cache": true,
		".venv":             true,
	}

	// Common dependency management files across different languages and tools
	targetFiles := map[string]string{
		// Go
		"go.mod": "Go modules",
		"go.sum": "Go modules checksum",

		// Python
		"requirements.txt": "Python pip",
		"Pipfile":          "Python pipenv",
		"pyproject.toml":   "Python poetry/modern Python packaging",
		"Pipfile.lock":     "Python pipenv lock",
		"poetry.lock":      "Python poetry lock",

		// Node.js
		"package.json":      "Node.js npm/yarn",
		"package-lock.json": "Node.js npm lock",
		"yarn.lock":         "Yarn lock",
		"pnpm-lock.yaml":    "pnpm lock",

		// Ruby
		"Gemfile":      "Ruby bundler",
		"Gemfile.lock": "Ruby bundler lock",

		// Java/Kotlin
		"pom.xml":          "Maven",
		"build.gradle":     "Gradle",
		"build.gradle.kts": "Gradle Kotlin",

		// .NET
		"packages.config": ".NET NuGet (legacy)",
		//"*.csproj":        ".NET project",
		//"*.fsproj":        "F# project",

		// Rust
		"Cargo.toml": "Rust cargo",
		"Cargo.lock": "Rust cargo lock",

		// PHP
		"composer.json": "PHP Composer",
		"composer.lock": "PHP Composer lock",

		//// Terraform
		//"terraform.tfstate": "Terraform state",
		//"*.tf":             "Terraform configuration",

		// Other
		"deps.edn":     "Clojure deps",
		"project.clj":  "Clojure Leiningen",
		"sbt":          "Scala build tool",
		"mix.exs":      "Elixir mix",
		"rebar.config": "Erlang rebar",
	}

	Findings, err := FindDependencies(root, ignoreDirs, targetFiles)

	if err != nil {
		log.Fatal().Err(err).Msg("Error finding dependencies")
	}

	var violations []Violation

	for _, finding := range Findings {
		switch finding.Dependency {
		case "requirements.txt", "Pipfile", "pyproject.toml", "Pipfile.lock", "poetry.lock":
			pythonViolations, err := ExamPython(finding)
			if err != nil {
				log.Fatal().Err(err).Msg("Error finding Python dependencies")
			}
			violations = append(violations, pythonViolations...)
		case "package.json", "yarn.lock", "pnpm-lock.yaml", "package-lock.json":
			packageViolations, err := ExamJS(finding)
			if err != nil {
				log.Fatal().Err(err).Msg("Error finding JS dependencies")
			}
			violations = append(violations, packageViolations...)
		case "go.mod", "go.sum":
			goViolations, err := ExamGo(finding)
			if err != nil {
				log.Fatal().Err(err).Msg("Error finding Go dependencies")
			}
			violations = append(violations, goViolations...)
		case "Gemfile":
			rubyViolations, err := examRuby(finding)
			if err != nil {
				log.Fatal().Err(err).Msg("Error finding Go dependencies")
			}
			violations = append(violations, rubyViolations...)

		case "pom.xml", "build.gradle", "build.gradle.kts":
			fmt.Printf("Java/Kotlin: %s\n", finding.Path)
		case "deps.edn", "project.clj", "sbt", "mix.exs", "rebar.config":
			fmt.Printf("Clojure: %s\n", finding.Path)
		case "Cargo.toml", "Cargo.lock":
			fmt.Printf("Rust: %s\n", finding.Path)
		case "composer.json", "composer.lock":
			fmt.Printf("PHP: %s\n", finding.Path)
		case "packages.config":
			fmt.Printf(".NET: %s\n", finding.Path)
		default:
			fmt.Printf("Unknown dependency: %s\n", finding.Dependency)
		}
	}

	writeViolations(violations)
}
