package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	Context struct {
		Request  *http.Request
		Response http.ResponseWriter
		Params   map[string]string
	}
)

func (ctx Context) RawRequestBody() []byte {
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}
	return data
}

func (ctx Context) JsonRequestBody(body interface{}) bool {
	err := json.Unmarshal(ctx.RawRequestBody(), body)
	return err == nil
}

func (ctx Context) BearerToken() string {
	return strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
}

func (ctx Context) WriteCsv(body []byte) {
	ctx.Response.Header().Set("Content-Type", "text/csv; charset=UTF-8")
	ctx.Response.Write(body)
}

func (ctx Context) WriteJson(body interface{}) {
	ctx.Response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	ctx.Response.Write(data)
}

func (ctx Context) WriteZip(body []byte) {
	ctx.Response.Header().Set("Content-Type", "application/zip")
	ctx.Response.Write(body)
}
