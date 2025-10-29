package std

import (
	"fmt"
	"runtime"
	"strings"
)

// Uses runtime.Callers to create a filtered stacktrace that excludes library code.
func GetLocalStackTrace(separator string) string {
	maxFrames := 32

	traceParts := make([]string, 1, maxFrames+1)
	traceParts[0] = "Goroutine stack trace (filtered):"

	// Start at depth 2 (skips runtime.Callers and GetLocalStackTrace)
	pc := make([]uintptr, 32)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()

		// Filter out Go standard library and module cache code
		switch {
		case strings.Contains(frame.File, "/go/src/"):
		case strings.Contains(frame.File, "/pkg/mod/"):
		default:
			// Format: function_name() path/to/file.go:line_number
			frameString := fmt.Sprintf("%s() %s:%d", frame.Function, frame.File, frame.Line)
			traceParts = append(traceParts, frameString)
		}

		if !more {
			break
		}
	}

	return strings.Join(traceParts, separator)
}
