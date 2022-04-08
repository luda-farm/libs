package router

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

type (
	Handler func(ctx Context)

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

func (r *Router) AllowOrigin(origin string) {
	r.allowedOrigins = append(r.allowedOrigins, origin)
}

func (r *Router) Delete(pattern string, handler Handler) {
	if handlers, ok := r.handlers[pattern]; ok {
		handlers.delete = handler
	} else {
		r.handlers[pattern] = &handlerGroup{delete: handler}
	}
}

func (r *Router) Get(pattern string, handler Handler) {
	if group, ok := r.handlers[pattern]; ok {
		group.get = handler
	} else {
		r.handlers[pattern] = &handlerGroup{get: handler}
	}
}

func (r *Router) Post(pattern string, handler Handler) {
	if handlers, ok := r.handlers[pattern]; ok {
		handlers.post = handler
	} else {
		r.handlers[pattern] = &handlerGroup{post: handler}
	}
}

func (r *Router) Put(pattern string, handler Handler) {
	if handlers, ok := r.handlers[pattern]; ok {
		handlers.put = handler
	} else {
		r.handlers[pattern] = &handlerGroup{put: handler}
	}
}

func (r Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer panicRecovery(res)
	for pattern, group := range r.handlers {
		matches, params := matchPattern(pattern, req.URL.Path)
		if !matches {
			continue
		}

		res.Header().Add("vary", "origin")
		res.Header().Add("access-control-allow-methods", group.allowedMethods())
		res.Header().Add("access-control-allow-headers", "authorization")
		res.Header().Add("access-control-allow-headers", "content-type")
		for _, origin := range r.allowedOrigins {
			if origin == "*" || origin == req.Header.Get("origin") {
				res.Header().Add("access-control-allow-origin", origin)
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

		ctx := Context{
			Request:  req,
			Response: res,
			Params:   params,
		}
		handler(ctx)
		return
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
