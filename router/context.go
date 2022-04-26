package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/luda-farm/libs/assert"
)

type (
	Context struct {
		request  *http.Request
		response http.ResponseWriter
		body     []byte
		Params   map[string]string
	}
)

func (ctx Context) BearerToken() string {
	return strings.TrimPrefix(ctx.request.Header.Get("authorization"), "Bearer ")
}

func (ctx Context) ReadHeader(key string) string {
	return ctx.request.Header.Get(key)
}

func (ctx Context) ReadBytes() []byte {
	return assert.Must(ioutil.ReadAll(ctx.request.Body))
}

// Writes error message to response body
func (ctx Context) ReadJson(body any) bool {
	if !strings.HasPrefix(ctx.ReadHeader("content-type"), "application/json") {
		ctx.Error(http.StatusBadRequest, "expected content-type application/json")
		return false
	}

	if json.Unmarshal(ctx.ReadBytes(), body) != nil {
		ctx.Error(http.StatusBadRequest, "failed to parse json")
		return false
	}
	return true
}

func (ctx *Context) WriteCsv(body []byte) {
	ctx.WriteHeader("content-type", "text/csv; charset=UTF-8")
	ctx.body = body
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

func (ctx Context) WriteHeader(key, value string) {
	ctx.response.Header().Set(key, value)
}

func (ctx Context) WriteJson(body any) {
	ctx.WriteHeader("Content-Type", "application/json; charset=UTF-8")
	data := assert.Must(json.Marshal(body))
	ctx.body = data
}

func (ctx Context) WriteZip(body []byte) {
	ctx.WriteHeader("Content-Type", "application/zip")
	ctx.body = body
}
