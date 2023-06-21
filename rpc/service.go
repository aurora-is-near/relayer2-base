package rpc

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/aurora-is-near/relayer2-base/log"
)

var (
	contextType = reflect.TypeOf((*context.Context)(nil)).Elem()
	stringType  = reflect.TypeOf("")
)

// ServiceMap holds all the services and subscriptions
type ServiceMap struct {
	mu            sync.Mutex
	services      map[string]service
	subscriptions map[string]subscription
}

// subscription is a registered rpc object storing the event subscription handlers
type subscription struct {
	name    string
	handler *handler
}

// service is a registered rpc object storing the service handlers
type service struct {
	name    string
	handler *handler
}

// handler is the object holding the json rpc method callbacks
type handler struct {
	name     string
	fn       reflect.Value
	rcvr     reflect.Value
	argTypes []reflect.Type
	hasCtx   bool
}

// register registers the methods and subscriptions provided by the service
func (s *ServiceMap) register(namespace string, rcvr interface{}, isSubscription bool) error {
	rcvrVal := reflect.ValueOf(rcvr)
	if namespace == "" {
		return fmt.Errorf("no service name for type %s", rcvrVal.Type().String())
	}
	handlers := getHandlers(namespace, rcvrVal)
	if len(handlers) == 0 {
		return fmt.Errorf("service %T doesn't have any suitable methods/subscriptions to expose", rcvr)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if isSubscription {
		if s.subscriptions == nil {
			s.subscriptions = make(map[string]subscription)
		}
		tmpSubscription := subscription{}
		for name, h := range handlers {
			tmpSubscription.handler = h
			tmpSubscription.name = name
			s.subscriptions[name] = tmpSubscription
		}

	} else {
		if s.services == nil {
			s.services = make(map[string]service)
		}
		tmpService := service{}
		for name, h := range handlers {
			tmpService.handler = h
			tmpService.name = name
			s.services[name] = tmpService
		}
	}
	return nil
}

// getHandlers iterates over the exposed methods of the given receiver type.
func getHandlers(namespace string, receiver reflect.Value) map[string]*handler {
	typ := receiver.Type()
	handlers := make(map[string]*handler)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if method.PkgPath != "" {
			continue // method not exported
		}
		h := newHandler(receiver, method.Func)
		if h == nil {
			continue // function invalid
		}
		name := namespace + "_" + strings.ToLower(method.Name)
		handlers[name] = h
	}
	return handlers
}

// newHandler turns the provided method into a handler and returns the created object
func newHandler(receiver, fn reflect.Value) *handler {
	fntype := fn.Type()
	h := &handler{fn: fn, rcvr: receiver}
	// Determine parameter types. They must all be exported or builtin types.
	h.makeArgTypes()

	// Verify return types. The function must return at most one error
	// and/or one other non-error value.
	outs := make([]reflect.Type, fntype.NumOut())
	for i := 0; i < fntype.NumOut(); i++ {
		outs[i] = fntype.Out(i)
	}
	if len(outs) > 2 {
		return nil
	}
	return h
}

// makeArgTypes composes the argTypes list.
func (h *handler) makeArgTypes() {
	fntype := h.fn.Type()
	// Skip receiver and context.Context parameter (if present).
	firstArg := 0
	if h.rcvr.IsValid() {
		firstArg++
	}
	if fntype.NumIn() > firstArg && fntype.In(firstArg) == contextType {
		h.hasCtx = true
		firstArg++
	}
	// Add all remaining parameters.
	h.argTypes = make([]reflect.Type, fntype.NumIn()-firstArg)
	for i := firstArg; i < fntype.NumIn(); i++ {
		h.argTypes[i-firstArg] = fntype.In(i)
	}
}

// call invokes the handler method
func (h *handler) call(ctx *context.Context, args []reflect.Value) (res interface{}, errRes error) {
	// Create the argument slice.
	fullargs := make([]reflect.Value, 0, 2+len(args))
	if h.rcvr.IsValid() {
		fullargs = append(fullargs, h.rcvr)
	}
	fullargs = append(fullargs, reflect.ValueOf(*ctx))
	fullargs = append(fullargs, args...)

	// Catch panic while running the callback.
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			log.Log().Error().Msgf("RPC method " + h.name + " crashed: " + fmt.Sprintf("%v\n%s", err, buf))
			errRes = errors.New("method handler crashed")
		}
	}()
	// Run the callback.
	results := h.fn.Call(fullargs)
	if len(results) == 0 {
		return nil, nil
	}
	if !results[1].IsNil() {
		// Method has returned non-nil error value.
		err := results[1].Interface().(error)
		return reflect.Value{}, err
	}
	return results[0].Interface(), nil
}
