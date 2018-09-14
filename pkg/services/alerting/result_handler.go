package alerting

import (
	"github.com/emretiryaki/merkut/pkg/log"
	"github.com/emretiryaki/merkut/pkg/services/rendering"
)

type ResultHandler interface {
	Handle(evalContext *EvalContext) error
}

type  DefaultResultHandler struct {
	notifier NotificationService
	log log.Logger
}

func NewResultHandler(renderService rendering.Service) *DefaultResultHandler{

	return &DefaultResultHandler{
		log:      log.New("alerting.resultHandler"),
		notifier: NewNotificationService(renderService),
	}
}


func (handler *DefaultResultHandler) Handle(evalContext *EvalContext) error {

	handler.notifier.SendIfNeeded(evalContext)

	return nil
}