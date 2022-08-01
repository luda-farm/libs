package httprouterdefaults

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

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
