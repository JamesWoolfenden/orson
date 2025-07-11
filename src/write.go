package orson

import (
	"fmt"
)

func writeViolations(violations []Violation) {
	for _, violation := range violations {
		fmt.Printf("%s %s: %s\n", violation.Path, violation.Dependency, violation.Style)
	}
}
