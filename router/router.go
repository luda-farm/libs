package router

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
)

const (
	Retry   = -1
	retries = 2
)

type (
	// Return status < 100 to trigger a retry
	Handler func(ctx Context) int

	Router struct {
		allowedOrigins []string
		handlers       map[string]*handlerGroup
	}

	handlerGroup struct {
		delete, get, post, put Handler
	}
)

func New() Router {
	return Router{handlers: make(map[string]*handlerGroup)}
}

func (router *Router) AllowOrigin(origin string) {
	router.allowedOrigins = append(router.allowedOrigins, origin)
}

func (router *Router) Delete(pattern string, handler Handler) {
	if handlers, ok := router.handlers[pattern]; ok {
		handlers.delete = handler
	} else {
		router.handlers[pattern] = &handlerGroup{delete: handler}
	}
}

func (router *Router) Get(pattern string, handler Handler) {
	if group, ok := router.handlers[pattern]; ok {
		group.get = handler
	} else {
		router.handlers[pattern] = &handlerGroup{get: handler}
	}
}

func (router *Router) Post(pattern string, handler Handler) {
	if handlers, ok := router.handlers[pattern]; ok {
		handlers.post = handler
	} else {
		router.handlers[pattern] = &handlerGroup{post: handler}
	}
}

func (router *Router) Put(pattern string, handler Handler) {
	if handlers, ok := router.handlers[pattern]; ok {
		handlers.put = handler
	} else {
		router.handlers[pattern] = &handlerGroup{put: handler}
	}
}

func (router Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer panicRecovery(res)
	for pattern, group := range router.handlers {
		matches, params := matchPattern(pattern, req.URL.Path)
		if !matches {
			continue
		}

		ctx := Context{
			request:  req,
			response: res,
			Params:   params,
		}

		// Handle CORS
		ctx.WriteHeader("vary", "origin")
		ctx.WriteHeader("access-control-allow-methods", group.allowedMethods())
		ctx.WriteHeader("access-control-allow-headers", "authorization, content-type")
		for _, origin := range router.allowedOrigins {
			if origin == "*" || origin == ctx.ReadHeader("origin") {
				ctx.WriteHeader("access-control-allow-origin", origin)
				break
			}
		}

		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		handler, ok := group.handler(req.Method)
		if !ok {
			break
		}

		for i := 0; i < retries; i++ {
			status := handler(ctx)
			if status == Retry {
				// random back off up to 100 ms
				time.Sleep(time.Duration(rand.Intn(1e8)))
				continue
			}

			res.WriteHeader(status)
			if status != http.StatusNoContent {
				res.Write(ctx.body)
			}
			return
		}
	}
	res.WriteHeader(http.StatusNotFound)
}

func panicRecovery(res http.ResponseWriter) {
	err := recover()
	if err == nil {
		return
	}

	stackTrace := make([]string, 0)
	for _, line := range strings.Split(string(debug.Stack()), "\n") {
		if regexp.MustCompile(`\.go:[1-9][0-9]*`).MatchString(line) {
			stackTrace = append(stackTrace, line, "\n")
		}
	}

	switch err := err.(type) {
	case error:
		fmt.Println("[PANIC]", err.Error(), "\n", stackTrace)
	default:
		fmt.Println("[PANIC]", err, "\n", stackTrace)
	}
	res.WriteHeader(http.StatusInternalServerError)
}

func matchPattern(rawPattern, rawRequest string) (bool, map[string]string) {
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

func (group handlerGroup) handler(method string) (Handler, bool) {
	switch method {
	case http.MethodDelete:
		return group.delete, group.delete != nil
	case http.MethodGet:
		return group.get, group.get != nil
	case http.MethodPost:
		return group.post, group.post != nil
	case http.MethodPut:
		return group.put, group.put != nil
	default:
		return nil, false
	}
}

func (group handlerGroup) allowedMethods() string {
	methods := make([]string, 0)
	if group.delete != nil {
		methods = append(methods, http.MethodDelete)
	}
	if group.get != nil {
		methods = append(methods, http.MethodGet)
	}
	if group.post != nil {
		methods = append(methods, http.MethodPost)
	}
	if group.put != nil {
		methods = append(methods, http.MethodPut)
	}
	return strings.Join(methods, ", ")
}
