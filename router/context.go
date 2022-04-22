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

func (ctx Context) RawRequestBody() []byte {
	return assert.Must(ioutil.ReadAll(ctx.Request.Body))
}

func (ctx Context) JsonRequestBody(body any) bool {
	return json.Unmarshal(ctx.RawRequestBody(), body) == nil
}

func (ctx Context) BearerToken() string {
	return strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
}

func (ctx Context) WriteCsv(body []byte) {
	ctx.Response.Header().Set("Content-Type", "text/csv; charset=UTF-8")
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
