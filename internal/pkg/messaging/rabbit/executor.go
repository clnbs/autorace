package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/streadway/amqp"

	"github.com/clnbs/autorace/internal/pkg/messaging"
)

type Executor struct {
	*Client
	router       *messaging.Router
	orchestrator *messaging.Orchestrator
}

func NewRabbitMQExecutor(client *Client) *Executor {
	executor := new(Executor)
	executor.Client = new(Client)
	executor.Client = client
	executor.orchestrator = messaging.NewOrchestrator()
	return executor
}

func (executor *Executor) RegisterRouter(router *messaging.Router) {
	executor.router = router
}

func (executor *Executor) Run() error {
	ready := make(chan bool)
	for _, route := range executor.router.RPCRoutes {
		go func() {
			err := executor.startRPCRoute(route, ready)
			if err != nil {
				// TODO error handling
			}
		}()
		<-ready
	}

	for _, route := range executor.router.ReceivingRoutes {
		go func() {
			err := executor.startReceivingRoute(route, ready)
			if err != nil {
				// TODO error handling
			}
		}()
		<-ready
	}
	return nil
}

func (executor *Executor) Stop() error {
	executor.orchestrator.SendMessage("")
	return nil
}

// TODO decoder as a dep injection
func (executor *Executor) SendMessageOnRoute(route string, routeOption *messaging.RouteOption, data interface{}) error {
	if reflect.TypeOf(routeOption.SpecOption) != reflect.TypeOf(SendingRouteOption{}) {
		return errors.New("bad router option type, got " + reflect.TypeOf(routeOption.SpecOption).String() + " should be " + reflect.TypeOf(SendingRouteOption{}).String())
	}
	switch data.(type) {
	case error:
		return data.(error)
	}
	bitifyMessage, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = executor.Channel.Publish(
		routeOption.SpecOption.(SendingRouteOption).Exchange,
		route,
		routeOption.SpecOption.(SendingRouteOption).Mandatory,
		routeOption.SpecOption.(SendingRouteOption).Immediate,
		amqp.Publishing{
			ContentType: routeOption.SpecOption.(SendingRouteOption).PublishingType,
			Timestamp:   time.Now(),
			Body:        bitifyMessage,
		},
	)
	return err
}

func (executor *Executor) createSendingMethodWithinContext(option *messaging.RouteOption) func(string, interface{}) error {
	return func(route string, data interface{}) error {
		return executor.SendMessageOnRoute(route, option, data)
	}
}

func (executor *Executor) startReceivingRoute(route *messaging.ReceivingRoute, ready chan bool) error {
	if reflect.TypeOf(route.Option.SpecOption) != reflect.TypeOf(ReceiveRouteOption{}) {
		ready <- true
		return errors.New("bad router option type, got " + reflect.TypeOf(route.Option.SpecOption).String() + " should be " + reflect.TypeOf(ReceiveRouteOption{}).String())
	}
	queue, err := executor.Channel.QueueDeclare(
		route.Option.SpecOption.(ReceiveRouteOption).Name,
		route.Option.SpecOption.(ReceiveRouteOption).Durable,
		route.Option.SpecOption.(ReceiveRouteOption).AutoDelete,
		route.Option.SpecOption.(ReceiveRouteOption).Exclusive,
		route.Option.SpecOption.(ReceiveRouteOption).NoWait,
		route.Option.SpecOption.(ReceiveRouteOption).Args,
	)
	if err != nil {
		ready <- true
		return err
	}
	err = executor.Channel.QueueBind(
		queue.Name,
		route.ListeningRoute,
		route.Option.SpecOption.(ReceiveRouteOption).Exchange,
		route.Option.SpecOption.(ReceiveRouteOption).NoWait,
		route.Option.SpecOption.(ReceiveRouteOption).Args,
	)
	if err != nil {
		ready <- true
		return err
	}

	msgs, err := executor.Channel.Consume(
		queue.Name,
		route.Option.SpecOption.(ReceiveRouteOption).Consumer,
		route.Option.SpecOption.(ReceiveRouteOption).AutoAck,
		route.Option.SpecOption.(ReceiveRouteOption).Exclusive,
		route.Option.SpecOption.(ReceiveRouteOption).NoLocal,
		route.Option.SpecOption.(ReceiveRouteOption).NoWait,
		route.Option.SpecOption.(ReceiveRouteOption).Args,
	)
	if err != nil {
		ready <- true
		return err
	}
	end := make(chan interface{})
	executor.orchestrator.RegisterListener(&end)
	go func() {
		for msg := range msgs {
			ctx := messaging.NewBrokerContext(context.TODO())
			ctx.Data = msg.Body
			ctx.Header = nil // TODO find a good header type
			ctx.ResponseTopic = ""
			ctx.Topic = route.ListeningRoute
			ctx.Send = nil
			ctx.Option = route.Option
			err = route.Handler(ctx)
			if err != nil {
				// TODO something with the error
			}
		}
	}()
	ready <- true
	<-end
	return nil
}

func (executor *Executor) startRPCRoute(route *messaging.RPCRoute, ready chan bool) error {
	if reflect.TypeOf(route.Option.SpecOption) != reflect.TypeOf(RPCRouteOption{}) {
		ready <- true
		return errors.New("bad router option type, got " + reflect.TypeOf(route.Option.SpecOption).String() + " should be " + reflect.TypeOf(RPCRouteOption{}).String())
	}
	if executor.Channel == nil {
		log.Println("oooops, channel is nil ....")
	}
	queue, err := executor.Channel.QueueDeclare(
		route.Option.SpecOption.(RPCRouteOption).Name,
		route.Option.SpecOption.(RPCRouteOption).Durable,
		route.Option.SpecOption.(RPCRouteOption).AutoDelete,
		route.Option.SpecOption.(RPCRouteOption).Exclusive,
		route.Option.SpecOption.(RPCRouteOption).NoWait,
		route.Option.SpecOption.(RPCRouteOption).Args,
	)
	if err != nil {
		ready <- true
		return err
	}
	err = executor.Channel.QueueBind(
		queue.Name,
		route.ListeningRoute,
		route.Option.SpecOption.(RPCRouteOption).Exchange,
		route.Option.SpecOption.(RPCRouteOption).NoWait,
		route.Option.SpecOption.(RPCRouteOption).Args,
	)
	if err != nil {
		ready <- true
		return err
	}

	msgs, err := executor.Channel.Consume(
		queue.Name,
		route.Option.SpecOption.(RPCRouteOption).Consumer,
		route.Option.SpecOption.(RPCRouteOption).AutoAck,
		route.Option.SpecOption.(RPCRouteOption).Exclusive,
		route.Option.SpecOption.(RPCRouteOption).NoLocal,
		route.Option.SpecOption.(RPCRouteOption).NoWait,
		route.Option.SpecOption.(RPCRouteOption).Args,
	)
	if err != nil {
		ready <- true
		return err
	}
	end := make(chan interface{})
	executor.orchestrator.RegisterListener(&end)
	go func() {
		for msg := range msgs {
			ctx := messaging.NewBrokerContext(context.TODO())
			ctx.Data = msg.Body
			ctx.Header = nil // TODO find a good header type
			ctx.Topic = route.ListeningRoute
			ctx.ResponseTopic = route.SendingRoute
			ctx.Send = executor.createSendingMethodWithinContext(route.Option)
			ctx.Option = route.Option
			err = route.Handler(ctx)
			if err != nil {
				//TODO something
			}
		}
	}()
	ready <- true
	<-end
	return nil
}
