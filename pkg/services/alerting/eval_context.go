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

