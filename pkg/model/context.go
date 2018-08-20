package model

import (
	session "github.com/emretiryaki/merkut/pkg/services/session"
	"gopkg.in/macaron.v1"
	"github.com/emretiryaki/merkut/pkg/log"
	"strings"
	"github.com/emretiryaki/merkut/pkg/setting"
)

type PermissionType int

const (
	PERMISSION_VIEW PermissionType = 1 << iota
	PERMISSION_EDIT
	PERMISSION_ADMIN
)
type ReqContext struct {
	*macaron.Context
	*SignedInUser

	Session session.SessionStore

	IsSignedIn     bool
	IsRenderCall   bool
	AllowAnonymous bool
	Logger         log.Logger
}


func (ctx *ReqContext) IsApiRequest() bool {
	return strings.HasPrefix(ctx.Req.URL.Path, "/api")
}

func (ctx *ReqContext) JsonApiErr(status int, message string, err error) {

	resp := make(map[string]interface{})
	if err != nil {
		ctx.Logger.Error(message, "error", err)
		if setting.Env != setting.PROD {
			resp["error"] = err.Error()
		}
	}
	switch status {
	case 404:
		resp["message"] = "Not Found"
	case 500:
		resp["message"] = "Internal Server Error"
	}

	if message != "" {
		resp["message"] = message
	}

	ctx.JSON(status, resp)
}

