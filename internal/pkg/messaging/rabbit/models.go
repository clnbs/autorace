package rabbit

import "github.com/streadway/amqp"

const (
	DefaultUser     = "guest"
	DefaultPassword = "guest"
	DefaultPort     = "5672"
)

type RPCRouteOption struct {
	*QueueDeclareOption
	*ConsumeOption
	*PublishOption
}

type ReceiveRouteOption struct {
	*QueueDeclareOption
	*ConsumeOption
	Exchange string
}

type SendingRouteOption struct {
	*PublishOption
}

type ExchangeDeclareOption struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

type QueueDeclareOption struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type PublishOption struct {
	Exchange       string
	Mandatory      bool
	Immediate      bool
	PublishingType string
}

type ConsumeOption struct {
	Consumer string
	AutoAck  bool
	NoLocal  bool
}

type Client struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}
