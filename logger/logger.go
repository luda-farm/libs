package logger

import (
	"fmt"
	"net/http"
)

// Writes a generic 500 response and logs the error with fmt.Errorf formatting.
func InternalServerError(w http.ResponseWriter, r *http.Request, format string, e error) {
	http.Error(w, "internal server error", http.StatusInternalServerError)
	fmt.Println(fmt.Errorf("[ERROR] %s: %w", r.URL.Path, fmt.Errorf(format, e)))
}
