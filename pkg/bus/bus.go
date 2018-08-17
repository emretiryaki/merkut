package bus

import (
	"context"
	"errors"
)


type HandlerFunc interface{}
type CtxHandlerFunc func()
type Msg interface{}

var ErrHandlerNotFound = errors.New("handler not found")

type TransactionManager interface {
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}


type Bus interface {
	Dispatch(msg Msg) error
	DispatchCtx(ctx context.Context, msg Msg) error
	Publish(msg Msg) error
	InTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	AddHandler(handler HandlerFunc)
	AddHandlerCtx(handler HandlerFunc)
	AddEventListener(handler HandlerFunc)
	AddWildcardListener(handler HandlerFunc)
	SetTransactionManager(tm TransactionManager)
}

type InProcBus struct {
	handlers          map[string]HandlerFunc
	handlersWithCtx   map[string]HandlerFunc
	listeners         map[string][]HandlerFunc
	wildcardListeners []HandlerFunc
	txMng             TransactionManager

}


func AddHandler(implName string, handler HandlerFunc) {
	globalBus.AddHandler(handler)
}


var globalBus = New()

func New() Bus {
	bus := &InProcBus{}
	bus.handlers = make(map[string]HandlerFunc)
	bus.handlersWithCtx = make(map[string]HandlerFunc)
	bus.listeners = make(map[string][]HandlerFunc)
	bus.wildcardListeners = make([]HandlerFunc, 0)
	bus.txMng = &noopTransactionManager{}

	return bus
}
type noopTransactionManager struct{}

func (*noopTransactionManager) InTransaction(ctx context.Context , fn func(ctx context.Context) error) error{
	return fn(ctx)
}