package router

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

type (
	clientError struct {
		Cause  string
		Status int
	}
)

func errorHandler(res http.ResponseWriter) {
	switch err := recover().(type) {
	case nil:
		return
	case clientError:
		res.WriteHeader(err.Status)
		data, _ := json.Marshal(err)
		res.Write(data)
	case error:
		log.Printf("500 Internal Server Error: %s\n%s", err.Error(), filteredStackTrace())
		res.WriteHeader(http.StatusInternalServerError)
	default:
		log.Printf("500 Internal Server Error: %v\n%s", err, filteredStackTrace())
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func filteredStackTrace() string {
	srcCodeLine := regexp.MustCompile(`((cmd)|(internal))/.+\.go`)
	buffer := strings.Builder{}
	for _, line := range strings.Split(string(debug.Stack()), "\n") {
		if srcCodeLine.MatchString(line) {
			buffer.WriteString(line)
			buffer.WriteRune('\n')
		}
	}
	return buffer.String()
}
