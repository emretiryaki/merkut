package middleware

import (
	"gopkg.in/macaron.v1"
	m "github.com/emretiryaki/merkut/pkg/model"
	"net/url"
	"github.com/emretiryaki/merkut/pkg/setting"
)

type AuthOptions struct {
	ReqGrafanaAdmin bool
	ReqSignedIn     bool
}


func Auth(options *AuthOptions) macaron.Handler {
	return func(c *m.ReqContext) {
		if !c.IsSignedIn && options.ReqSignedIn && !c.AllowAnonymous {
			notAuthorized(c)
			return
		}

		if !c.IsGrafanaAdmin && options.ReqGrafanaAdmin {
			accessForbidden(c)
			return
		}
	}
}

func notAuthorized(c *m.ReqContext){
	if c.IsApiRequest() {
		c.JsonApiErr(401, "Unauthorized", nil)
		return
	}
	c.SetCookie("redirect_to", url.QueryEscape(setting.AppSubUrl+c.Req.RequestURI), 0, setting.AppSubUrl+"/", nil, false, true)

	c.Redirect(setting.AppSubUrl + "/login")
}

func accessForbidden(c *m.ReqContext) {
	if c.IsApiRequest() {
		c.JsonApiErr(403, "Permission denied", nil)
		return
	}

	c.Redirect(setting.AppSubUrl + "/")
}