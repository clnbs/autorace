package rabbit

import (
	"errors"
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"github.com/streadway/amqp"
	"log"
	"reflect"
)

var ErrorNoHostSetup = errors.New("no host is setup for broker RabbitMQ")

type Factory struct {
	*Client
	url          string
	RouterOption ExchangeDeclareOption
}

func NewRabbitMQFactory(routerOption interface{}) (*Factory, error) {
	if reflect.TypeOf(routerOption) != reflect.TypeOf(ExchangeDeclareOption{}) {
		return nil, errors.New("bad router option type, got " + reflect.TypeOf(routerOption).String() + " should be " + reflect.TypeOf(ExchangeDeclareOption{}).String())
	}
	rabbitFactory := new(Factory)
	rabbitFactory.RouterOption = routerOption.(ExchangeDeclareOption)
	rabbitFactory.Client = new(Client)
	return rabbitFactory, nil
}

func (factory *Factory) SetupRouter(client *messaging.BrokerClientOption) (*messaging.Router, error) {
	var err error
	if client.Password == "" {
		client.Password = DefaultPassword
	}
	if client.User == "" {
		client.User = DefaultUser
	}
	if client.Port == "" {
		client.Port = DefaultPort
	}
	if client.Host == "" {
		return nil, ErrorNoHostSetup
	}
	factory.url = "amqp://" + client.User + ":" + client.Password + "@" + client.Host + ":" + client.Port + "/"
	factory.Connection, err = amqp.Dial(factory.url)
	if err != nil {
		return nil, err
	}
	factory.Channel, err = factory.Connection.Channel()
	if err != nil {
		return nil, err
	}
	err = factory.Channel.ExchangeDeclare(
		factory.RouterOption.Name,
		factory.RouterOption.Kind,
		factory.RouterOption.Durable,
		factory.RouterOption.AutoDelete,
		factory.RouterOption.Internal,
		factory.RouterOption.NoWait,
		factory.RouterOption.Args,
	)
	if err != nil {
		return nil, err
	}
	router := new(messaging.Router)
	router.ReceivingRoutes = make(map[string]*messaging.ReceivingRoute)
	router.RPCRoutes = make(map[string]*messaging.RPCRoute)
	if factory.Channel == nil {
		log.Println("channel is nil here too ...")
	}
	return router, nil
}

func (factory *Factory) GetRouterOption() interface{} {
	return factory.RouterOption
}

func (factory *Factory) AddReceiveOnlyRoute(router *messaging.Router, routeOption *messaging.RouteOption, route string, handler messaging.Handler) error {
	if reflect.TypeOf(routeOption.SpecOption) != reflect.TypeOf(ReceiveRouteOption{}) {
		return errors.New("bad option type for listener")
	}
	router.ReceivingRoutes[routeOption.Name] = &messaging.ReceivingRoute{
		Handler:        handler,
		ListeningRoute: route,
		Option:         routeOption,
	}
	return nil
}

func (factory *Factory) AddRPCRoute(router *messaging.Router, routeOption *messaging.RouteOption, listeningRoute string, sendingRoute string, handler messaging.Handler) error {
	if reflect.TypeOf(routeOption.SpecOption) != reflect.TypeOf(RPCRouteOption{}) {
		return errors.New("bad option type for RPC")
	}
	router.RPCRoutes[routeOption.Name] = &messaging.RPCRoute{
		Handler:        handler,
		SendingRoute:   sendingRoute,
		ListeningRoute: listeningRoute,
		Option:         routeOption,
	}
	return nil
}
