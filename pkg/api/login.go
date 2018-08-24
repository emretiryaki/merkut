package api

import (
	"strings"
	"github.com/emretiryaki/merkut/pkg/setting"
	m "github.com/emretiryaki/merkut/pkg/model"
		"github.com/emretiryaki/merkut/pkg/api/dto"

)
const (
	ViewIndex = "index"
)


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
		Locale:                     locale,
		AppUrl:                  appURL,
		AppSubUrl:               appSubURL,
		AppName:                 setting.ApplicationName,
	}


	themeURLParam := c.Query("theme")
	if themeURLParam == "light" {
		data.Theme = "light"
	}


	return &data, nil
}