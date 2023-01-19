package httputil

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

type (
	CorsConfig struct {
		Origins []string
		Methods []string
		Headers []string
	}
)

// Access-Control-Allow-Origin defaults to "*" if no origins are configured.
func CorsHandler(config CorsConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(config.Origins) < 1 {
			w.Header().Add("access-control-allow-origin", "*")
		} else {
			for _, origin := range config.Origins {
				if origin == r.Header.Get("origin") {
					w.Header().Add("access-control-allow-origin", origin)
					break
				}
			}
		}

		for _, method := range config.Methods {
			w.Header().Add("access-control-allow-methods", method)
		}

		for _, header := range config.Headers {
			w.Header().Add("access-control-allow-headers", header)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func PanicHandler(w http.ResponseWriter, r *http.Request, err any) {
	isLocalCode := regexp.MustCompile(`((cmd)|(internal))/.+\.go`)
	var stackTrace strings.Builder
	for _, line := range strings.Split(string(debug.Stack()), "\n") {
		if isLocalCode.MatchString(line) {
			stackTrace.WriteString(line)
			stackTrace.WriteRune('\n')
		}
	}

	msg := fmt.Sprintf("%#v", err)
	if e, ok := err.(error); ok {
		msg = e.Error()
	}

	log.Printf("500 Internal Server Error: %s\n%s", msg, stackTrace.String())
	w.WriteHeader(http.StatusInternalServerError)
}
