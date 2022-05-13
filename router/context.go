package router

import (
	"net/http"
	"strings"
)

type (
	Context struct {
		Response http.ResponseWriter
		Request  *http.Request
		Params   map[string]string
	}
)

func (ctx Context) BearerToken() string {
	return strings.TrimPrefix(ctx.Request.Header.Get("authorization"), "Bearer ")
}

func (ctx Context) CheckClientError(err error, status int) {
	if err != nil {
		panic(clientError{Cause: err.Error(), Status: status})
	}
}
