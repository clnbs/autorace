package messaging

import "context"

type BrokerContext struct {
	context       context.Context
	Data          []byte
	Header        []byte //??
	ResponseTopic string
	Topic         string
	Option        *RouteOption
	Send          func(string, interface{}) error
}

func NewBrokerContext(ctx context.Context) *BrokerContext {
	brokerCtx := new(BrokerContext)
	brokerCtx.context = ctx
	return brokerCtx
}

func (brokerContext BrokerContext) Context() context.Context {
	return brokerContext.context
}
