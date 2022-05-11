package router

import (
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"golang.org/x/exp/maps"
)

type (
	Router struct {
		allowedOrigins []string
		handlers       map[string]map[string]Handler
	}

	Handler func(res http.ResponseWriter, req *http.Request, params map[string]string)
)

func New() Router {
	return Router{handlers: map[string]map[string]Handler{}}
}

func (router *Router) AllowOrigin(origin string) {
	router.allowedOrigins = append(router.allowedOrigins, origin)
}

func (router *Router) Handle(method, path string, handler Handler) {
	_, ok := router.handlers[path]
	if !ok {
		router.handlers[path] = map[string]Handler{
			http.MethodOptions: func(res http.ResponseWriter, req *http.Request, _ map[string]string) {
				res.WriteHeader(http.StatusNoContent)
			},
		}
	}
	router.handlers[path][method] = handler
}

func (router Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer func() {
		switch err := recover().(type) {
		case nil:
			return
		case error:
			log.Printf("500 Internal Server Error: %s\n%s", err.Error(), filteredStackTrace())
		default:
			log.Printf("500 Internal Server Error: %v\n%s", err, filteredStackTrace())
		}
		res.WriteHeader(http.StatusInternalServerError)
	}()

	for path, methods := range router.handlers {
		matches, params := matchPath(path, req.URL.Path)
		if !matches {
			continue
		}

		// Handle CORS
		res.Header().Set("vary", "origin")
		res.Header().Add("access-control-allow-headers", "authorization")
		res.Header().Add("access-control-allow-headers", "content-type")
		for _, method := range maps.Keys(methods) {
			res.Header().Add("access-control-allow-methods", method)
		}
		for _, origin := range router.allowedOrigins {
			if origin == "*" || origin == req.Header.Get("origin") {
				res.Header().Set("access-control-allow-origin", origin)
				break
			}
		}

		handler, ok := methods[req.Method]
		if !ok {
			break
		}

		handler(res, req, params)
		return
	}
}

func matchPath(rawPattern, rawRequest string) (bool, map[string]string) {
	pattern, request := strings.Split(rawPattern, "/"), strings.Split(rawRequest, "/")
	if len(pattern) != len(request) {
		return false, nil
	}

	params := make(map[string]string)
	for i, segment := range pattern {
		if strings.HasPrefix(segment, "$") {
			param := segment[1:]
			params[param] = request[i]
			continue
		}

		if pattern[i] != request[i] {
			return false, nil
		}
	}
	return true, params
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
