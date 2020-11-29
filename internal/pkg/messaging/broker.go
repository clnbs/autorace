package messaging

//Handler is the basic handler function prototype. User handlers has to be prototyped like so
type Handler func(brokerContext *BrokerContext) error

// BrokerFactory handle setup, router and broker management
type BrokerFactory interface {
	//SetupRouter create a Router from client options
	SetupRouter(client *BrokerClientOption) (*Router, error)

	// AddReceiveOnlyRoute is use to register new listening route (like data digestion, etc ...)
	AddReceiveOnlyRoute(router *Router, routeOption *RouteOption, route string, handler Handler) error
	// AddRPCRoute is use to register new RPC route (basically, a listening function with a callback)
	AddRPCRoute(router *Router, routeOption *RouteOption, listeningRoute string, sendingRoute string, handler Handler) error
}

// BrokerExecutor handle broker life cycle and useful functions when a broker is setup
type BrokerExecutor interface {
	// Run start broker's router set up by the broker factory
	Run() error
	// Stop stop broker's router set up by the broker factory
	Stop() error

	// RegisterRouter is use to register to Executor the router set up by the factory
	RegisterRouter(router *Router)
	// SendMessageOnRoute is use to send a blob message to some route by using router set up by the factory
	SendMessageOnRoute(route string, routeOption *RouteOption, data interface{}) error
}

// BrokerClientOption discribe how router are connected to the real broker (RabbitMQ, MQTT, Kafka, etc ...)
type BrokerClientOption struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Broker handle router with a bigger picture, meant to be used by end user
type Broker struct {
	Option   *BrokerClientOption
	factory  BrokerFactory
	executor BrokerExecutor
	router   *Router
}

// NewBroker create a Broker with client option, a broker factory and a broker executor
func NewBroker(option *BrokerClientOption, factory BrokerFactory, executor BrokerExecutor) (*Broker, error) {
	var err error
	broker := new(Broker)
	broker.executor = executor
	broker.factory = factory
	broker.router, err = broker.factory.SetupRouter(option)
	if err != nil {
		return nil, err
	}
	broker.executor.RegisterRouter(broker.router)
	return broker, nil
}

// AddRPCHandler wrap up factory's function
func (broker *Broker) AddRPCHandler(listeningRoute, responseRoute string, option *RouteOption, handler Handler) error {
	return broker.factory.AddRPCRoute(broker.router, option, listeningRoute, responseRoute, handler)
}

// AddHandler wrap up factory's function
func (broker *Broker) AddHandler(route string, option *RouteOption, handler Handler) error {
	return broker.factory.AddReceiveOnlyRoute(broker.router, option, route, handler)
}

// Send wrap up executor function
func (broker *Broker) Send(route string, option *RouteOption, data interface{}) error {
	return broker.executor.SendMessageOnRoute(route, option, data)
}

// Run wrap up executor function
func (broker *Broker) Run() {
	err := broker.executor.Run()
	if err != nil {
		// TODO something
	}
}

// Stop wrap up executor function
func (broker *Broker) Stop() {
	err := broker.executor.Stop()
	if err != nil {
		// TODO something
	}
}
