package std

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Filters the stack trace to only print line references originating in the local source code
func PrintLocalStackTrace() {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	var filteredStack strings.Builder
	for _, line := range lines {
		if !strings.Contains(line, ".go:") {
			continue
		}
		if strings.Contains(line, "/cmd") || strings.Contains(line, "/internal") {
			filteredStack.WriteString(line + "\n")
		}
	}
	fmt.Println(filteredStack.String())
}
