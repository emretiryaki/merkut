package dto


type CurrentUser struct {
	IsSignedIn                 bool         `json:"isSignedIn"`
	Id                         int64        `json:"id"`
	Login                      string       `json:"login"`
	Email                      string       `json:"email"`
	Name                       string       `json:"name"`
	LightTheme                 bool         `json:"lightTheme"`
	OrgCount                   int          `json:"orgCount"`
	OrgId                      int64        `json:"orgId"`
	OrgName                    string       `json:"orgName"`
	Locale                     string       `json:"locale"`
}

