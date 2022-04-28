package router

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/luda-farm/libs/assert"
)

type (
	Context struct {
		Request  *http.Request
		response http.ResponseWriter
		body     []byte
		Params   map[string]string
	}
)

func (ctx Context) BearerToken() string {
	return strings.TrimPrefix(ctx.Request.Header.Get("authorization"), "Bearer ")
}

func (ctx Context) RequestBodyAsBytes() []byte {
	return assert.Must(io.ReadAll(ctx.Request.Body))
}

func (ctx Context) RequestBodyAsJson(body any) bool {
	return nil == json.Unmarshal(ctx.RequestBodyAsBytes(), body)
}

func (ctx *Context) Error(status int, msg string) int {
	if status < 400 || status > 599 {
		panic("call to 'Context.Error' with non-error status code")
	}
	ctx.WriteJson(struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: msg,
	})
	return status
}

func (ctx *Context) WriteJson(body any) {
	ctx.body = assert.Must(json.Marshal(body))
}
