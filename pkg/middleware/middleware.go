package middleware

import( "gopkg.in/macaron.v1"
m "github.com/emretiryaki/merkut/pkg/model"
			"github.com/emretiryaki/merkut/pkg/log"
	)

func AddDefaultResponseHeaders() macaron.Handler {
	return func(ctx *m.ReqContext) {
		if ctx.IsApiRequest() && ctx.Req.Method == "GET" {
			ctx.Resp.Header().Add("Cache-Control", "no-cache")
			ctx.Resp.Header().Add("Pragma", "no-cache")
			ctx.Resp.Header().Add("Expires", "-1")
		}
	}
}


func GetContextHandler() macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &m.ReqContext{
			Context:        c,
			AllowAnonymous: false,
			Logger:         log.New("context"),
		}
		ctx.Data["ctx"] = ctx
		c.Header().Add("Access-Control-Allow-Origin","*")
		c.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		c.Map(ctx)


	}
}


