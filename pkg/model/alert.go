package model

type AlertStateType string
type AlertSeverityType string
type NoDataOption string
type ExecutionErrorOption string

const (
	AlertStateNoData   AlertStateType = "no_data"
	AlertStatePaused   AlertStateType = "paused"
	AlertStateAlerting AlertStateType = "alerting"
	AlertStateOK       AlertStateType = "ok"
	AlertStatePending  AlertStateType = "pending"
)


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

