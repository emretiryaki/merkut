package alerting

import (
	"fmt"
	"regexp"
	"github.com/emretiryaki/merkut/pkg/components/simplejson"
	m "github.com/emretiryaki/merkut/pkg/model"
	"strconv"
	"time"
)
type Rule struct {
	Id                  int64
	Frequency           int64
	Name                string
	Message             string
	NoDataState         m.NoDataOption
	ExecutionErrorState m.ExecutionErrorOption
	State               m.AlertStateType
	Conditions          []Condition
	Notifications       []int64
	Comment             string
	LastFired           time.Time
	LastTriggered       time.Time
	Schedule            string
	When                string
	Indice              string
}

type ValidationError struct {
	Reason      string
	Err         error
	Alertid     int64
	DashboardId int64
	PanelId     int64
}

func (e ValidationError) Error() string {
	extraInfo := ""
	if e.Alertid != 0 {
		extraInfo = fmt.Sprintf("%s AlertId: %v", extraInfo, e.Alertid)
	}

	if e.PanelId != 0 {
		extraInfo = fmt.Sprintf("%s PanelId: %v ", extraInfo, e.PanelId)
	}

	if e.DashboardId != 0 {
		extraInfo = fmt.Sprintf("%s DashboardId: %v", extraInfo, e.DashboardId)
	}

	if e.Err != nil {
		return fmt.Sprintf("%s %s%s", e.Err.Error(), e.Reason, extraInfo)
	}

	return fmt.Sprintf("Failed to extract alert.Reason: %s %s", e.Reason, extraInfo)
}

var (
	ValueFormatRegex = regexp.MustCompile(`^\d+`)
	UnitFormatRegex  = regexp.MustCompile(`\w{1}$`)
)

var unitMultiplier = map[string]int{
	"s": 1,
	"m": 60,
	"h": 3600,
}

func getTimeDurationStringToSeconds(str string) (int64, error) {
	multiplier := 1

	matches := ValueFormatRegex.FindAllString(str, 1)

	if len(matches) <= 0 {
		return 0, fmt.Errorf("Frequency could not be parsed")
	}

	value, err := strconv.Atoi(matches[0])
	if err != nil {
		return 0, err
	}

	unit := UnitFormatRegex.FindAllString(str, 1)[0]

	if val, ok := unitMultiplier[unit]; ok {
		multiplier = val
	}

	return int64(value * multiplier), nil
}



func NewRuleFromDBAlert(ruleDef *m.Alert) (*Rule, error) {

	model := &Rule{}

	model.Id = ruleDef.Id
	model.Name = ruleDef.Name
	model.State = ruleDef.State
	model.Comment = ruleDef.Comment
	model.LastFired = ruleDef.LastFired
	model.LastTriggered = ruleDef.LastTriggered
	model.Schedule = ruleDef.Schedule
	model.When = ruleDef.When
	model.Indice = ruleDef.Indice

	//model.ExecutionErrorState = m.ExecutionErrorOption(ruleDef.Settings.Get("executionErrorState").MustString("alerting"))

	//for _, v := range ruleDef.Settings.Get("notifications").MustArray() {
	//	jsonModel := simplejson.NewFromAny(v)
	//	id, err := jsonModel.Get("id").Int64()
	//	if err != nil {
	//		return nil, ValidationError{Reason: "Invalid notification schema", DashboardId: model.DashboardId, Alertid: model.Id, PanelId: model.PanelId}
	//	}
	//	model.Notifications = append(model.Notifications, id)
	//}
	//
	//for index, condition := range ruleDef.Settings.Get("conditions").MustArray() {
	//	conditionModel := simplejson.NewFromAny(condition)
	//	conditionType := conditionModel.Get("type").MustString()
	//	factory, exist := conditionFactories[conditionType]
	//	if !exist {
	//		return nil, ValidationError{Reason: "Unknown alert condition: " + conditionType, DashboardId: model.DashboardId, Alertid: model.Id, PanelId: model.PanelId}
	//	}
	//	queryCondition, err := factory(conditionModel, index)
	//	if err != nil {
	//		return nil, ValidationError{Err: err, DashboardId: model.DashboardId, Alertid: model.Id, PanelId: model.PanelId}
	//	}
	//	model.Conditions = append(model.Conditions, queryCondition)
	//}
	//
	//if len(model.Conditions) == 0 {
	//	return nil, fmt.Errorf("Alert is missing conditions")
	//}

	//Id             int64
	//Name           string
	//State          AlertStateType
	//Comment        string
	//LastFired	   string
	//LastTriggered  string
	//Schedule 	   string
	//When 		   string
	//Indice 		   string
	return &Rule{}, nil
}

type ConditionFactory func(model *simplejson.Json, index int) (Condition, error)

var conditionFactories = make(map[string]ConditionFactory)

func RegisterCondition(typeName string, factory ConditionFactory) {
	conditionFactories[typeName] = factory
}
