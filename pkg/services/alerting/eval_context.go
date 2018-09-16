package alerting

import (
	"context"
	"time"

	"github.com/emretiryaki/merkut/pkg/log"
	m "github.com/emretiryaki/merkut/pkg/model"
)

type EvalContext struct {
	Firing         bool
	IsTestRun      bool
	EvalMatches    []*EvalMatch
	Logs           []*ResultLogEntry
	Error          error
	ConditionEvals string
	StartTime      time.Time
	EndTime        time.Time
	Rule           *Rule
	log            log.Logger

	NoDataFound     bool
	PrevAlertState  m.AlertStateType
	Ctx context.Context


}

func NewEvalContext(alertCtx context.Context, rule *Rule) *EvalContext {
	return &EvalContext{
		Ctx:            alertCtx,
		StartTime:      time.Now(),
		Rule:           rule,
		Logs:           make([]*ResultLogEntry, 0),
		EvalMatches:    make([]*EvalMatch, 0),
		log:            log.New("alerting.evalContext"),
		PrevAlertState: rule.State,
	}
}


func (c *EvalContext) GetNewState() m.AlertStateType {

	if c.Error != nil {
		c.log.Error("Alert Rule Result Error",
			"ruleId", c.Rule.Id,
			"name", c.Rule.Name,
			"error", c.Error,
			"changing state to", c.Rule.ExecutionErrorState.ToAlertState())

		if c.Rule.ExecutionErrorState == m.ExecutionErrorKeepState {
			return c.PrevAlertState
		}
		return c.Rule.ExecutionErrorState.ToAlertState()

	} else if c.Firing {
		return m.AlertStateAlerting

	} else if c.NoDataFound {
		c.log.Info("Alert Rule returned no data",
			"ruleId", c.Rule.Id,
			"name", c.Rule.Name,
			"changing state to", c.Rule.NoDataState.ToAlertState())

		if c.Rule.NoDataState == m.NoDataKeepState {
			return c.PrevAlertState
		}
		return c.Rule.NoDataState.ToAlertState()
	}

	return m.AlertStateOK
}

func (a *EvalContext) GetDurationMs() float64 {
	return float64(a.EndTime.Nanosecond()-a.StartTime.Nanosecond()) / float64(1000000)
}
