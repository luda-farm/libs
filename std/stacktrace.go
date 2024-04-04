package std

import (
	"fmt"
	"runtime/debug"
	"strings"
)

// Prints the lines from the stack trace that originates in 'cmd' or 'internal'
func PrintLocalStackTrace() {
	stack := string(debug.Stack())
	lines := strings.Split(stack, "\n")
	var filteredStack strings.Builder
	for _, line := range lines {
		if strings.Contains(line, "/cmd") || strings.Contains(line, "/internal") {
			filteredStack.WriteString(line + "\n")
		}
	}
	fmt.Println(filteredStack.String())
}
