package dto


type IndexViewData struct {
	Settings                map[string]interface{}
	AppUrl                  string
	AppSubUrl               string
	GoogleAnalyticsId       string
	GoogleTagManagerId      string
	BuildVersion            string
	BuildCommit             string
	Theme                   string
	AppName                 string
	Locale                  string
}
