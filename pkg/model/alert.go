package model

import "github.com/emretiryaki/merkut/pkg/components/null"

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



type Job struct {
	Offset     int64
	OffsetWait bool
	Delay      bool
	Running    bool
	//Rule       *Rule
}

type ResultLogEntry struct {
	Message string
	Data    interface{}
}

type EvalMatch struct {
	Value  null.Float        `json:"value"`
	Metric string            `json:"metric"`
	Tags   map[string]string `json:"tags"`
}

type Level struct {
	Operator string
	Value    float64
}
