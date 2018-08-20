package api

import (
	"net/url"
	"strings"
	"github.com/emretiryaki/merkut/pkg/setting"
	m "github.com/emretiryaki/merkut/pkg/model"
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/api/dto"
	session "github.com/emretiryaki/merkut/pkg/services/session"

)
const (
	ViewIndex = "index"
)
func Logout(c *m.ReqContext) {
	c.SetCookie(setting.CookieUserName, "", -1, setting.AppSubUrl+"/")
	c.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubUrl+"/")
	c.Session.Destory(c.Context)
	if setting.SignoutRedirectUrl != "" {
		c.Redirect(setting.SignoutRedirectUrl)
	} else {
		c.Redirect(setting.AppSubUrl + "/login")
	}
}


func LoginView(c *m.ReqContext) {
	viewData, _ := setIndexViewData(c)

	enabledOAuths := make(map[string]interface{})
	for key, oauth := range setting.OAuthService.OAuthInfos {
		enabledOAuths[key] = map[string]string{"name": oauth.Name}
	}

	if loginError, ok := c.Session.Get("loginError").(string); ok {
		c.Session.Delete("loginError")
		viewData.Settings["loginError"] = loginError
	}

	if !tryLoginUsingRememberCookie(c) {
		c.HTML(200, ViewIndex, viewData)
		return
	}

	if redirectTo, _ := url.QueryUnescape(c.GetCookie("redirect_to")); len(redirectTo) > 0 {
		c.SetCookie("redirect_to", "", -1, setting.AppSubUrl+"/")
		c.Redirect(redirectTo)
		return
	}

	c.Redirect(setting.AppSubUrl + "/")
}


func tryLoginUsingRememberCookie(c *m.ReqContext) bool {
	// Check auto-login.
	uname := c.GetCookie(setting.CookieUserName)
	if len(uname) == 0 {
		return false
	}

	isSucceed := false
	defer func() {
		if !isSucceed {
			c.SetCookie(setting.CookieUserName, "", -1, setting.AppSubUrl+"/")
			c.SetCookie(setting.CookieRememberName, "", -1, setting.AppSubUrl+"/")
			return
		}
	}()

	userQuery := m.GetUserByLoginQuery{LoginOrEmail: uname}
	if err := bus.Dispatch(&userQuery); err != nil {
		return false
	}

	user := userQuery.Result

	// validate remember me cookie
	if val, _ := c.GetSuperSecureCookie(user.Rands+user.Password, setting.CookieRememberName); val != user.Login {
		return false
	}

	isSucceed = true
	loginUserWithUser(user, c)
	return true
}

func loginUserWithUser(user *m.User, c *m.ReqContext) {

	c.Resp.Header().Del("Set-Cookie")

	days := 86400 * setting.LogInRememberDays
	if days > 0 {
		c.SetCookie(setting.CookieUserName, user.Login, days, setting.AppSubUrl+"/")
		c.SetSuperSecureCookie(user.Rands+user.Password, setting.CookieRememberName, user.Login, days, setting.AppSubUrl+"/")
	}

	c.Session.RegenerateId(c.Context)
	c.Session.Set(session.SESS_KEY_USERID, user.Id)
}

func setIndexViewData(c *m.ReqContext) (*dto.IndexViewData, error) {

	acceptLang := c.Req.Header.Get("Accept-Language")
	locale := "en-US"

	if len(acceptLang) > 0 {
		parts := strings.Split(acceptLang, ",")
		locale = parts[0]
	}

	appURL := setting.AppUrl
	appSubURL := setting.AppSubUrl


	var data = dto.IndexViewData{
		User: &dto.CurrentUser{
			Id:                         c.UserId,
			IsSignedIn:                 c.IsSignedIn,
			Login:                      c.Login,
			Email:                      c.Email,
			Name:                       c.Name,
			Locale:                     locale,
		},
		AppUrl:                  appURL,
		AppSubUrl:               appSubURL,
		AppName:                 setting.ApplicationName,
	}


	if len(data.User.Name) == 0 {
		data.User.Name = data.User.Login
	}

	themeURLParam := c.Query("theme")
	if themeURLParam == "light" {
		data.User.LightTheme = true
		data.Theme = "light"
	}


	return &data, nil
}