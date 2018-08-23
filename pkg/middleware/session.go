package middleware

import (
	"gopkg.in/macaron.v1"
	"github.com/emretiryaki/merkut/pkg/services/session"
	ms "github.com/go-macaron/session"
	m "github.com/emretiryaki/merkut/pkg/model"
)

func Sessioner(options *ms.Options, sessionConnMaxLifetime int64) macaron.Handler {
	session.Init(options, sessionConnMaxLifetime)

	return func(ctx *m.ReqContext) {
		ctx.Next()

		if err := ctx.Session.Release(); err != nil {
			panic("session(release): " + err.Error())
		}
	}
}

