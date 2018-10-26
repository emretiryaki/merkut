package alerting

import (
	"github.com/emretiryaki/merkut/pkg/bus"
	"github.com/emretiryaki/merkut/pkg/log"
	m "github.com/emretiryaki/merkut/pkg/model"
	"sync"
	"time"
)

type RuleReader interface {
	Fetch() []*Rule
}

type DefaultRuleReader struct {
	sync.RWMutex
	//serverID       string
	serverPosition int
	clusterSize    int
	log            log.Logger
}

func NewRuleReader()  *DefaultRuleReader{

	ruleReader := &DefaultRuleReader{
		log:log.New("alerting.ruleReader"),
	}

	go ruleReader.initReader()

	return ruleReader
}

func (arr *DefaultRuleReader)  initReader() {

	heartbeat := time.NewTicker(time.Second*10)

	for {
		select {
			case <-heartbeat.C:
				arr.heartbeat()
		}
	}
}



func (arr *DefaultRuleReader) Fetch() []*Rule {

	cmdGetAllAlertsQuery := &m.GetAlertsQuery{}

	if err := bus.Dispatch(cmdGetAllAlertsQuery); err != nil {
		arr.log.Error("Could not load alerts", "error", err)
		return []*Rule{}
	}

	cmdGetActions :=&m.GetAllActionsQuery{}

	if err := bus.Dispatch(cmdGetActions); err != nil {
		arr.log.Error("Could not load actions", "error", err)
		return []*Rule{}
	}

	res := make([]*Rule, 0)
	for _, ruleDef := range cmdGetAllAlertsQuery.Result {


		//for _,actionItem :=range cmdGetActions.Result{
		//
		//}

		if model, err := NewRuleFromDBAlert(ruleDef); err != nil {
			arr.log.Error("Could not build alert model for rule", "ruleId", ruleDef.Id, "error", err)
		} else {
			res = append(res, model)
		}
	}

	//elastic	metrics.M_Alerting_Active_Alerts.Set(float64(len(res)))
	return res
}



func (arr *DefaultRuleReader)  heartbeat() {
	arr.clusterSize = 1
	arr.serverPosition = 1

}