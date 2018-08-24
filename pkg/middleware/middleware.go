package middleware

import( "gopkg.in/macaron.v1"
m "github.com/emretiryaki/merkut/pkg/model"
		"github.com/emretiryaki/merkut/pkg/services/session"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/bus"
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
			SignedInUser:   &m.SignedInUser{},
			Session:        session.GetSession(),
			IsSignedIn:     false,
			AllowAnonymous: false,
			Logger:         log.New("context"),
		}

		ctx.Logger = log.New("context", "userId", ctx.UserId, "orgId", ctx.OrgId, "uname", ctx.Login)
		ctx.Data["ctx"] = ctx

		c.Map(ctx)

		// update last seen every 5min
		if ctx.ShouldUpdateLastSeenAt() {
			ctx.Logger.Debug("Updating last user_seen_at", "user_id", ctx.UserId)
			if err := bus.Dispatch(&m.UpdateUserLastSeenAtCommand{UserId: ctx.UserId}); err != nil {
				ctx.Logger.Error("Failed to update last_seen_at", "error", err)
			}
		}
	}
}

