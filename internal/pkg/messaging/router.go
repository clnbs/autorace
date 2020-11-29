package messaging

import "context"

type RouteOption struct {
	Name       string
	SpecOption interface{}
}

type Router struct {
	RPCRoutes       map[string]*RPCRoute
	ReceivingRoutes map[string]*ReceivingRoute
}

type RPCRoute struct {
	Handler        Handler
	SendingRoute   string
	ListeningRoute string
	Option         *RouteOption
	Terminating    chan interface{}
	Context        context.Context
}

type ReceivingRoute struct {
	Handler        Handler
	ListeningRoute string
	Option         *RouteOption
	Terminating    chan interface{}
	Context        context.Context
}
