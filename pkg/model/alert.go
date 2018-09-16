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

const (
	NoDataSetOK       NoDataOption = "ok"
	NoDataSetNoData   NoDataOption = "no_data"
	NoDataKeepState   NoDataOption = "keep_state"
	NoDataSetAlerting NoDataOption = "alerting"
)

const   (
	ExecutionErrorSetAlerting ExecutionErrorOption = "alerting"
	ExecutionErrorKeepState   ExecutionErrorOption = "keep_state"
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
	Indice 		   string
}


type GetAllAlertsQuery struct {
	Result []*Alert
}


func (s AlertStateType) IsValid() bool {
	return s == AlertStateOK || s == AlertStateNoData || s == AlertStatePaused || s == AlertStatePending
}

func (s NoDataOption) IsValid() bool {
	return s == NoDataSetNoData || s == NoDataSetAlerting || s == NoDataKeepState || s == NoDataSetOK
}

func (s NoDataOption) ToAlertState() AlertStateType {
	return AlertStateType(s)
}

func (s ExecutionErrorOption) IsValid() bool {
	return s == ExecutionErrorSetAlerting || s == ExecutionErrorKeepState
}

func (s ExecutionErrorOption) ToAlertState() AlertStateType {
	return AlertStateType(s)
}