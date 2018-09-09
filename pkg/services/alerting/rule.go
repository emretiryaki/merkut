package alerting

import (
	"fmt"
	"regexp"
	"github.com/emretiryaki/merkut/pkg/components/simplejson"
	m "github.com/emretiryaki/merkut/pkg/model"
	"strconv"
)
type Rule struct {
	Id                  int64
	OrgId               int64
	DashboardId         int64
	PanelId             int64
	Frequency           int64
	Name                string
	Message             string
	NoDataState         m.NoDataOption
	ExecutionErrorState m.ExecutionErrorOption
	State               m.AlertStateType
	Conditions          []Condition
	Notifications       []int64
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

type ConditionFactory func(model *simplejson.Json, index int) (Condition, error)

var conditionFactories = make(map[string]ConditionFactory)

func RegisterCondition(typeName string, factory ConditionFactory) {
	conditionFactories[typeName] = factory
}
