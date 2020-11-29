package rabbit

import (
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func mockReceiverHandler(ctx *messaging.BrokerContext) error {
	handlerName := ctx.Option.Name
	log.Println(handlerName, ": message received on topic", ctx.Topic)
	log.Println(handlerName, ": data :", string(ctx.Data))
	return nil
}

func mockRPCHandler(ctx *messaging.BrokerContext) error {
	handlerName := ctx.Option.Name
	log.Println(handlerName, ": message received on topic", ctx.Topic)
	log.Println(handlerName, ": data :", string(ctx.Data))

	log.Println(handlerName, "is about to send back some data ...")
	return ctx.Send(ctx.ResponseTopic, "from "+handlerName+" : your data was : "+string(ctx.Data))
}

func TestNewRabbitMQFactory(t *testing.T) {
	_, err := NewRabbitMQFactory(ExchangeDeclareOption{
		Name:       "testing_topic",
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	assert.Nil(t, err, "while creating a new rabbit factory :", err)
	_, err = NewRabbitMQFactory("")
	assert.NotNil(t, err, "while creating new rabbit factory, error should be rise, no viable option registered")
}

func TestFactory_SetupRouter(t *testing.T) {
	factory, err := NewRabbitMQFactory(ExchangeDeclareOption{
		Name:       "testing_topic",
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	assert.Nil(t, err, "while creating a new rabbit factory :", err)

	_, err = factory.SetupRouter(&messaging.BrokerClientOption{
		Host:     "localhost",
		Port:     "",
		User:     "",
		Password: "",
	})
	assert.Nil(t, err, "while setting up router :", err)

	_, err = factory.SetupRouter(&messaging.BrokerClientOption{
		Host:     "",
		Port:     "",
		User:     "",
		Password: "",
	})
	assert.NotNil(t, err, "while setting up router, an error should be rise : no viable option")
}

func TestFactory_AddReceiveOnlyRoute(t *testing.T) {
	factory, err := NewRabbitMQFactory(ExchangeDeclareOption{
		Name:       "testing_topic",
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	assert.Nil(t, err, "while creating a new rabbit factory :", err)

	router, err := factory.SetupRouter(&messaging.BrokerClientOption{
		Host:     "localhost",
		Port:     "",
		User:     "",
		Password: "",
	})
	assert.Nil(t, err, "while setting up router :", err)

	err = factory.AddReceiveOnlyRoute(
		router,
		&messaging.RouteOption{
			Name: "first_rabbit_receiver",
			SpecOption: ReceiveRouteOption{
				QueueDeclareOption: &QueueDeclareOption{
					Name:       "",
					Durable:    false,
					AutoDelete: false,
					Exclusive:  true,
					NoWait:     false,
					Args:       nil,
				},
				ConsumeOption: &ConsumeOption{
					Consumer: "first_rabbit_receiver",
					AutoAck:  true,
					NoLocal:  false,
				},
				Exchange: "testing_topic",
			},
		},
		"testing.rabbit.router.first",
		mockReceiverHandler,
	)
	assert.Nil(t, err, "error while registering a receiver handler")

	err = factory.AddReceiveOnlyRoute(
		router,
		&messaging.RouteOption{
			Name:       "first_rabbit_receiver",
			SpecOption: RPCRouteOption{},
		},
		"testing.rabbit.router.first",
		mockReceiverHandler,
	)
	assert.NotNil(t, err, "while registering a receiver handler, an error should be rise")
}

func TestFactory_AddRPCRoute(t *testing.T) {
	factory, err := NewRabbitMQFactory(ExchangeDeclareOption{
		Name:       "testing_topic",
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	assert.Nil(t, err, "while creating a new rabbit factory :", err)

	router, err := factory.SetupRouter(&messaging.BrokerClientOption{
		Host:     "localhost",
		Port:     "",
		User:     "",
		Password: "",
	})
	assert.Nil(t, err, "while setting up router :", err)

	err = factory.AddRPCRoute(
		router,
		&messaging.RouteOption{
			Name: "first_rabbit_RPC",
			SpecOption: RPCRouteOption{
				QueueDeclareOption: &QueueDeclareOption{
					Name:       "",
					Durable:    false,
					AutoDelete: false,
					Exclusive:  true,
					NoWait:     false,
					Args:       nil,
				},
				ConsumeOption: &ConsumeOption{
					Consumer: "first_rabbit_RPC",
					AutoAck:  false,
					NoLocal:  false,
				},
				PublishOption: &PublishOption{
					Exchange:       "testing_topic",
					Mandatory:      false,
					Immediate:      false,
					PublishingType: "topic",
				},
			},
		},
		"testing.rabbit.rpc.first.listen",
		"testing.rabbit.rpc.first.send",
		mockReceiverHandler,
	)
	assert.Nil(t, err, "while registering RPC handler :", err)

	err = factory.AddRPCRoute(
		router,
		&messaging.RouteOption{
			Name:       "first_rabbit_RPC",
			SpecOption: ReceiveRouteOption{},
		},
		"testing.rabbit.rpc.first.listen",
		"testing.rabbit.rpc.first.send",
		mockReceiverHandler,
	)
}

func TestFactory_GetRouterOption(t *testing.T) {
	factory, err := NewRabbitMQFactory(ExchangeDeclareOption{
		Name:       "testing_topic",
		Kind:       "topic",
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	})
	assert.Nil(t, err, "while creating a new rabbit factory :", err)

	_, err = factory.SetupRouter(&messaging.BrokerClientOption{
		Host:     "localhost",
		Port:     "",
		User:     "",
		Password: "",
	})
	assert.Nil(t, err, "while setting up router :", err)

	option := factory.GetRouterOption()
	log.Println("option :", option)
}
