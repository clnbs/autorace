package messaging

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type FactoryMock struct {
	mutex sync.Mutex
}

func (f *FactoryMock) SetupRouter(client *BrokerClientOption) (*Router, error) {
	if client.Host == "" {
		return nil, errors.New("not host set up")
	}
	router := new(Router)
	router.ReceivingRoutes = make(map[string]*ReceivingRoute)
	router.RPCRoutes = make(map[string]*RPCRoute)
	return router, nil
}

func (f *FactoryMock) AddReceiveOnlyRoute(router *Router, routeOption *RouteOption, route string, handler Handler) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	router.ReceivingRoutes[routeOption.Name] = &ReceivingRoute{
		Handler:        handler,
		ListeningRoute: route,
		Option:         routeOption,
	}
	return nil
}

func (f *FactoryMock) AddRPCRoute(router *Router, routeOption *RouteOption, listeningRoute string, sendingRoute string, handler Handler) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	router.RPCRoutes[routeOption.Name] = &RPCRoute{
		Handler:        handler,
		SendingRoute:   sendingRoute,
		ListeningRoute: listeningRoute,
		Option:         routeOption,
	}
	return nil
}

type ExecutorMock struct {
	router       *Router
	orchestrator *Orchestrator
}

func NewExecutorMock() *ExecutorMock {
	e := &ExecutorMock{
		router:       nil,
		orchestrator: NewOrchestrator(),
	}
	return e
}

func (e *ExecutorMock) RegisterRouter(router *Router) {
	e.router = router
}

func (e *ExecutorMock) Run() error {
	ready := make(chan bool)
	for routeName, route := range e.router.RPCRoutes {
		go func() {
			log.Println("about to start", routeName)
			err := e.startRPCRoute(route, ready)
			if err != nil {
				// TODO something
			}
		}()
		<-ready
	}
	for routeName, route := range e.router.ReceivingRoutes {
		go func() {
			log.Println("about to start", routeName)
			err := e.startReceivingRoute(route, ready)
			if err != nil {
				//TODO error handling
			}
		}()
		<-ready
	}
	return nil
}

func (e *ExecutorMock) Stop() error {
	e.orchestrator.SendMessage("")
	return nil
}

func (e *ExecutorMock) SendMessageOnRoute(route string, routeOption *RouteOption, data interface{}) error {
	log.Println("mock sending fake message to", route, "with option", routeOption, " ... :")
	log.Println(data)
	return nil
}

func (e *ExecutorMock) createSendingMethodWithinContextMock(option *RouteOption) func(string, interface{}) error {
	return func(route string, data interface{}) error {
		return e.SendMessageOnRoute(route, option, data)
	}
}

func (e *ExecutorMock) startReceivingRoute(route *ReceivingRoute, ready chan bool) error {
	end := make(chan interface{})
	e.orchestrator.RegisterListener(&end)
	go func() {
		err := route.Handler(&BrokerContext{
			context:       context.TODO(),
			Data:          []byte("testing receiving route..."),
			Header:        nil,
			ResponseTopic: "",
			Topic:         route.ListeningRoute,
			Send:          nil,
			Option:        route.Option,
		})
		if err != nil {
			// TODO error handling
		}
	}()
	ready <- true
	<-end
	return nil
}

func (e *ExecutorMock) startRPCRoute(route *RPCRoute, ready chan bool) error {
	end := make(chan interface{})
	e.orchestrator.RegisterListener(&end)
	go func() {
		err := route.Handler(&BrokerContext{
			context:       context.TODO(),
			Data:          []byte("testing RPC route"),
			Header:        nil,
			ResponseTopic: route.SendingRoute,
			Topic:         route.ListeningRoute,
			Option:        route.Option,
			Send:          e.createSendingMethodWithinContextMock(route.Option),
		})
		if err != nil {
			// TODO error handling
		}
	}()
	ready <- true
	<-end
	return nil
}

func fakeHandlingFunctionWithCallBack(ctx *BrokerContext) error {
	log.Println("mock handler named \""+ctx.Option.Name+"\" received a message on topic :", ctx.Topic)
	log.Println("mock handler named \""+ctx.Option.Name+"\" received :", string(ctx.Data))

	log.Println(ctx.Option.Name, "is about to send a message on :", ctx.ResponseTopic)
	return ctx.Send(ctx.ResponseTopic, "some message from a mock handler")
}

func fakeHandlingFunction(ctx *BrokerContext) error {
	log.Println("mock handler named \""+ctx.Option.Name+"\" received a message on topic :", ctx.Topic)
	log.Println("mock handler named \""+ctx.Option.Name+"\" received :", string(ctx.Data))
	return nil
}

func TestNewBroker(t *testing.T) {
	factoryMock := &FactoryMock{}
	executorMock := NewExecutorMock()
	broker, err := NewBroker(&BrokerClientOption{
		Host:     "localhost",
		Port:     "",
		User:     "",
		Password: "",
	},
		factoryMock,
		executorMock,
	)
	assert.Nil(t, err, "error while creating new broker with mock :", err)
	err = broker.AddHandler(
		"mock.listener.first",
		&RouteOption{
			Name:       "first_listener",
			SpecOption: nil,
		},
		fakeHandlingFunction,
	)
	assert.Nil(t, err, "error while registering first listener :", err)
	err = broker.AddHandler(
		"mock.listener.two",
		&RouteOption{
			Name:       "second_listener",
			SpecOption: nil,
		},
		fakeHandlingFunction,
	)
	assert.Nil(t, err, "error while registering second listener :", err)
	err = broker.AddRPCHandler(
		"mock.rpc.listen.first",
		"mock.rpc.response.first",
		&RouteOption{
			Name:       "first_RPC_handler",
			SpecOption: nil,
		},
		fakeHandlingFunctionWithCallBack,
	)
	assert.Nil(t, err, "error while registering first RPC handler :", err)
	err = broker.AddRPCHandler(
		"mock.rpc.listen.two",
		"mock.rpc.response.two",
		&RouteOption{
			Name:       "second_RPC_handler",
			SpecOption: nil,
		},
		fakeHandlingFunctionWithCallBack,
	)
	assert.Nil(t, err, "error while registering second RPC handler :", err)
	broker.Run()
	time.Sleep(5 * time.Second)
	broker.Stop()
}
