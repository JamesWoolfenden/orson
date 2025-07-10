package orson

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type Finding struct {
	Dependency string
	Style      string
	Path       string
}

func FindDependencies(root string, ignoreDirs map[string]bool, targetFiles map[string]string) ([]Finding, error) {

	var findings []Finding

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		// Skip ignored directories
		if info.IsDir() && ignoreDirs[info.Name()] {
			return filepath.SkipDir
		}

		var found Finding
		found.Dependency = filepath.Base(path)
		found.Path = filepath.ToSlash(path)

		// Check exact matches
		if description, exists := targetFiles[found.Dependency]; exists {
			found.Style = description
			findings = append(findings, found)
			return nil
		}

		// Check pattern matches (like *.csproj)
		for pattern, description := range targetFiles {
			if strings.HasPrefix(pattern, "*.") {
				ext := strings.TrimPrefix(pattern, "*")
				if strings.HasSuffix(found.Dependency, ext) {
					found.Style = description
					findings = append(findings, found)
					//fmt.Printf("Found %s (%s): %s\n", found.Dependency, found.Style, found.Path)
					return nil
				}
			}
		}

		return nil
	})

	fmt.Println("Found", len(findings), "dependency management files")

	if err != nil {
		log.Err(err).Msgf("Error walking the path: %v\n", err)
		return nil, err
	}

	return findings, nil
}
