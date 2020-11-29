package rabbit

import (
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRabbitMQExecutor(t *testing.T) {
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
			Name: "second_rabbit_receiver",
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
					Consumer: "second_rabbit_receiver",
					AutoAck:  true,
					NoLocal:  false,
				},
				Exchange: "testing_topic",
			},
		},
		"testing.rabbit.router.two",
		mockReceiverHandler,
	)
	assert.Nil(t, err, "error while registering a receiver handler")

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
		mockRPCHandler,
	)
	assert.Nil(t, err, "while registering RPC handler :", err)

	err = factory.AddRPCRoute(
		router,
		&messaging.RouteOption{
			Name: "second_rabbit_RPC",
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
					Consumer: "second_rabbit_RPC",
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
		"testing.rabbit.rpc.two.listen",
		"testing.rabbit.rpc.two.send",
		mockRPCHandler,
	)
	assert.Nil(t, err, "while registering RPC handler :", err)

	executor := NewRabbitMQExecutor(factory.Client)
	executor.RegisterRouter(router)
	go func() {
		err = executor.Run()
		assert.Nil(t, err, "while starting executor :", err)
	}()
	time.Sleep(5 * time.Second)
	err = executor.SendMessageOnRoute(
		"testing.rabbit.router.first",
		&messaging.RouteOption{
			Name: "",
			SpecOption: SendingRouteOption{&PublishOption{
				Exchange:       "testing_topic",
				Mandatory:      false,
				Immediate:      false,
				PublishingType: "topic",
			}},
		},
		"toto",
	)
	assert.Nil(t, err, "while sending data on some route :", err)

	err = executor.SendMessageOnRoute(
		"testing.rabbit.router.two",
		&messaging.RouteOption{
			Name: "",
			SpecOption: SendingRouteOption{&PublishOption{
				Exchange:       "testing_topic",
				Mandatory:      false,
				Immediate:      false,
				PublishingType: "topic",
			}},
		},
		"tata",
	)
	assert.Nil(t, err, "while sending data on some route :", err)

	err = executor.SendMessageOnRoute(
		"testing.rabbit.rpc.first.listen",
		&messaging.RouteOption{
			Name: "",
			SpecOption: SendingRouteOption{&PublishOption{
				Exchange:       "testing_topic",
				Mandatory:      false,
				Immediate:      false,
				PublishingType: "topic",
			}},
		},
		"titi",
	)
	assert.Nil(t, err, "while sending data on some route :", err)

	err = executor.SendMessageOnRoute(
		"testing.rabbit.rpc.two.listen",
		&messaging.RouteOption{
			Name: "",
			SpecOption: SendingRouteOption{&PublishOption{
				Exchange:       "testing_topic",
				Mandatory:      false,
				Immediate:      false,
				PublishingType: "topic",
			}},
		},
		"tutu",
	)
	assert.Nil(t, err, "while sending data on some route :", err)

	time.Sleep(2 * time.Second)
	err = executor.Stop()
	assert.Nil(t, err, "while stopping broker router :", err)
}
