package model

type Alert struct {

	Id             int64
	Name           string
	State          string
	Comment        string
	LastFired	   string
	LastTriggered  string
	Schedule 	   string
	When 		   string

}

type GetAllAlertsQuery struct {
	Result []*Alert
}