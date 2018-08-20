package dto


type IndexViewData struct {
	User                    *CurrentUser
	Settings                map[string]interface{}
	AppUrl                  string
	AppSubUrl               string
	GoogleAnalyticsId       string
	GoogleTagManagerId      string
	BuildVersion            string
	BuildCommit             string
	Theme                   string
	NewGrafanaVersionExists bool
	NewGrafanaVersion       string
	AppName                 string
}
