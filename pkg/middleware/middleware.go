package middleware

import( "gopkg.in/macaron.v1"
m "github.com/emretiryaki/merkut/pkg/model")

func AddDefaultResponseHeaders() macaron.Handler {
	return func(ctx *m.ReqContext) {
		if ctx.IsApiRequest() && ctx.Req.Method == "GET" {
			ctx.Resp.Header().Add("Cache-Control", "no-cache")
			ctx.Resp.Header().Add("Pragma", "no-cache")
			ctx.Resp.Header().Add("Expires", "-1")
		}
	}
}
