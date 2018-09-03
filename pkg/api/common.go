package api

import (
	"encoding/json"
	m "github.com/emretiryaki/merkut/pkg/model"
	"github.com/emretiryaki/merkut/pkg/setting"
	"gopkg.in/macaron.v1"
	"net/http"
)
var (

	ServerError = func(err error) Response {
		return Error(500, "Server error", err)
	}
)
func Wrap(action interface{}) macaron.Handler {

	return func(c *m.ReqContext) {
		var res Response
		val, err := c.Invoke(action)
		if err == nil && val != nil && len(val) > 0 {
			res = val[0].Interface().(Response)
		} else {
			res = ServerError(err)
		}

		res.WriteTo(c)
	}
}


type Response interface {
	WriteTo(ctx *m.ReqContext)
}

type NormalResponse struct {
	status     int
	body       []byte
	header     http.Header
	errMessage string
	err        error
}
// Error create a erroneous response
func Error(status int, message string, err error) *NormalResponse {
	data := make(map[string]interface{})

	switch status {
	case 404:
		data["message"] = "Not Found"
	case 500:
		data["message"] = "Internal Server Error"
	}

	if message != "" {
		data["message"] = message
	}

	if err != nil {
		if setting.Env != setting.PROD {
			data["error"] = err.Error()
		}
	}

	resp := JSON(status, data)

	if err != nil {
		resp.errMessage = message
		resp.err = err
	}

	return resp
}
// JSON create a JSON response
func JSON(status int, body interface{}) *NormalResponse {
	return Respond(status, body).Header("Content-Type", "application/json")
}

// Respond create a response
func Respond(status int, body interface{}) *NormalResponse {
	var b []byte
	var err error
	switch t := body.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		if b, err = json.Marshal(body); err != nil {
			return Error(500, "body json marshal", err)
		}
	}
	return &NormalResponse{
		body:   b,
		status: status,
		header: make(http.Header),
	}
}

func (r *NormalResponse) Header(key, value string) *NormalResponse {
	r.header.Set(key, value)
	return r
}

func (r *NormalResponse) WriteTo(ctx *m.ReqContext) {
	if r.err != nil {
		ctx.Logger.Error(r.errMessage, "error", r.err)
	}

	header := ctx.Resp.Header()
	for k, v := range r.header {
		header[k] = v
	}
	ctx.Resp.WriteHeader(r.status)
	ctx.Resp.Write(r.body)
}