package bus

import (
	"context"
	"errors"
	"reflect"
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

func GetBus()  Bus{
	return globalBus
}

func (b *InProcBus) Dispatch(msg Msg) error {
	var msgName = reflect.TypeOf(msg).Elem().Name()

	var handler = b.handlersWithCtx[msgName]
	withCtx := true

	if handler == nil {
		withCtx = false
		handler = b.handlers[msgName]
	}

	if handler == nil {
		return ErrHandlerNotFound
	}

	var params = []reflect.Value{}
	if withCtx {
		params = append(params, reflect.ValueOf(context.Background()))
	}
	params = append(params, reflect.ValueOf(msg))

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}

func (b *InProcBus) DispatchCtx(ctx context.Context, msg Msg) error {

	var msgName = reflect.TypeOf(msg).Elem().Name()

	var handler = b.handlersWithCtx[msgName]
	if handler == nil {
		return ErrHandlerNotFound
	}

	var params = []reflect.Value{}
	params = append(params, reflect.ValueOf(ctx))
	params = append(params, reflect.ValueOf(msg))

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	}
	return err.(error)
}

func (b *InProcBus)  Publish(msg Msg)  error{

	var msgName = reflect.TypeOf(msg).Elem().Name()
	var listeners = b.listeners[msgName]

	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(msg)

	for _, listenerHandler := range listeners{
		ret := reflect.ValueOf(listenerHandler).Call(params)
		err := ret[0].Interface()
		if err != nil {
			return err.(error)
		}
	}

	for _, listenerHandler := range b.wildcardListeners {
		ret := reflect.ValueOf(listenerHandler).Call(params)
		err := ret[0].Interface()
		if err != nil {
			return err.(error)
		}
	}

	return nil
}

func (b *InProcBus) AddWildcardListener(handler HandlerFunc) {
	b.wildcardListeners = append(b.wildcardListeners, handler)
}

func (b *InProcBus) AddHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(0).Elem().Name()
	b.handlers[queryTypeName] = handler
}

func (b *InProcBus) AddHandlerCtx(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(1).Elem().Name()
	b.handlersWithCtx[queryTypeName] = handler
}

func Publish(msg Msg) error {
	return globalBus.Publish(msg)
}

func (b *InProcBus) AddEventListener(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventName := handlerType.In(0).Elem().Name()
	_, exists := b.listeners[eventName]
	if !exists {
		b.listeners[eventName] = make([]HandlerFunc, 0)
	}
	b.listeners[eventName] = append(b.listeners[eventName], handler)
}

func (b *InProcBus) SetTransactionManager(tm TransactionManager) {
	b.txMng = tm
}

func (b *InProcBus) InTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return b.txMng.InTransaction(ctx, fn)
}

func Dispatch(msg Msg) error {
	return globalBus.Dispatch(msg)
}

type noopTransactionManager struct{}

func (*noopTransactionManager) InTransaction(ctx context.Context , fn func(ctx context.Context) error) error{
	return fn(ctx)
}

