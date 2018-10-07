package model

type ConditionType string


const (

	ConditionTypeCompare ConditionType = "compare"
	ConditionTypeNever ConditionType = "never"
	ConditionTypeAlways ConditionType = "always"

)
type	Condition struct {
	Id				int64
	Alarm_id		int64
	Type 			ConditionType
	Operator 		string

}
