package alerting

import (
	"github.com/benbjohnson/clock"
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/registry"
	"github.com/emretiryaki/merkut/pkg/services/rendering"
	"time"
)


type AlertingService struct {
	RenderService rendering.Service `inject:""`
	execQueue chan *Job
	//clock         clock.Clock
	ticker        *Ticker
	scheduler     Scheduler
	evalHandler   EvalHandler
	ruleReader    RuleReader
	log           log.Logger
	resultHandler ResultHandler

}


func init() {
	registry.RegisterService(&AlertingService{})
}

func NewEngine() *AlertingService {
	e := &AlertingService{}
	e.Init()
	return e
}


func (e *AlertingService) Init() error {
	e.ticker = NewTicker(time.Now(), time.Second*0, clock.New())
	e.execQueue = make(chan *Job, 1000)
	e.scheduler = NewScheduler()
	e.evalHandler = NewEvalHandler()
	e.ruleReader = NewRuleReader()
	e.log = log.New("alerting.engine")
	e.resultHandler = NewResultHandler(e.RenderService)
	return nil
}

