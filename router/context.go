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
		Request  *http.Request
		Response http.ResponseWriter
		Params   map[string]string
	}
)

func (ctx Context) BearerToken() string {
	return strings.TrimPrefix(ctx.Request.Header.Get("authorization"), "Bearer ")
}

func (ctx Context) ReadBytes() []byte {
	return assert.Must(ioutil.ReadAll(ctx.Request.Body))
}

// Writes error to context on failure
func (ctx Context) ReadJson(body any) bool {
	if !strings.HasPrefix(ctx.Request.Header.Get("content-type"), "application/json") {
		ctx.WriteError(http.StatusBadRequest, "expected content-type application/json")
		return false
	}

	if json.Unmarshal(ctx.ReadBytes(), body) != nil {
		ctx.WriteError(http.StatusBadRequest, "failed to parse json")
		return false
	}
	return true
}

func (ctx Context) WriteCsv(body []byte) {
	ctx.Response.Header().Set("content-type", "text/csv; charset=UTF-8")
	assert.Must(ctx.Response.Write(body))
}

func (ctx Context) WriteError(status int, msg string) {
	if status < 400 || status > 599 {
		panic("call to 'Context.WriteError' with non-error status code")
	}
	ctx.Response.WriteHeader(status)
	ctx.WriteJson(struct {
		Status  int
		Message string
	}{
		Status:  status,
		Message: msg,
	})
}

func (ctx Context) WriteJson(body any) {
	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	data := assert.Must(json.Marshal(body))
	assert.Must(ctx.Response.Write(data))
}

func (ctx Context) WriteZip(body []byte) {
	ctx.Response.Header().Set("Content-Type", "application/zip")
	assert.Must(ctx.Response.Write(body))
}
