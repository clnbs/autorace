package rabbit

import (
	"encoding/json"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/clnbs/autorace/pkg/logger"

	"github.com/streadway/amqp"
)

// ErrorResponse is sent back when a error occurs
type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
}

// RabbitConnection hold a representation of a RabbitMQ connection
type RabbitConnection struct {
	RabbitURL       string
	conn            *amqp.Connection
	partiesChannel  *amqp.Channel
	lobbiesChannel  *amqp.Channel
	name            string
	ReceivedMessage map[string]chan []byte
	SendingMessage  map[string]chan []byte
}

// RabbitConnectionConfiguration hold configuration to make a RabbitMQ connection possible
type RabbitConnectionConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
}

// NewRabbitConnection create a RabbitConnection from a RabbitMQ address
func NewRabbitConnection(config RabbitConnectionConfiguration) (*RabbitConnection, error) {
	rConn := new(RabbitConnection)
	rConn.ReceivedMessage = make(map[string]chan []byte)
	rConn.SendingMessage = make(map[string]chan []byte)
	var err error
	rConn.RabbitURL = "amqp://" + config.User + ":" + config.Password + "@" + config.Host + ":" + config.Port + "/"
	rConn.conn, err = amqp.Dial(rConn.RabbitURL)
	if err != nil {
		return nil, err
	}
	// declare channels
	rConn.partiesChannel, err = rConn.conn.Channel()
	if err != nil {
		return nil, err
	}
	rConn.lobbiesChannel, err = rConn.conn.Channel()
	if err != nil {
		return nil, err
	}

	// declare exchange
	err = rConn.partiesChannel.ExchangeDeclare(
		"parties_topic", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}
	err = rConn.lobbiesChannel.ExchangeDeclare(
		"lobbies_topic", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	return rConn, nil
}

// SendMessageOnTopic is used to send a message on a specific topic
func (rConn *RabbitConnection) SendMessageOnTopic(message interface{}, topic string) {
	switch message.(type) {
	case error:
		return
	}
	bitifyMessage, err := json.Marshal(message)
	if err != nil {
		logger.Error("while marshaling \"get\" :", err)
		return
	}
	err = rConn.partiesChannel.Publish(
		"parties_topic", // exchange
		topic,           // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bitifyMessage,
		})
	if err != nil {
		logger.Error("error while sending message :", err)
	}
}

// ReceiveMessageOnTopic is used to receive a specific message on a given topic. This method
// send back the receiving object trough a chan of interface. The receiving object is computed by
// func passed in argument and has to be declared like the following example :
// `func handler(msg []byte) interface{}`
// The slice of byte pass in argument of `handler` is the body of the message received on the topic
func (rConn *RabbitConnection) ReceiveMessageOnTopic(topic string, handler func([]byte) interface{}, communicationChan chan interface{}, readyToReceive chan bool) error {
	queue, err := rConn.partiesChannel.QueueDeclare(
		"",    // name
		false, // type
		false, // durable
		true,  // auto-deleted
		false, // internal
		nil,   // no-wait
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	err = rConn.partiesChannel.QueueBind(
		queue.Name,      //queue name
		topic,           //routing key
		"parties_topic", // exchange
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}

	logger.Trace("waiting message on topic :", topic)
	msgs, err := rConn.partiesChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	loop := make(chan os.Signal)
	signal.Notify(loop, os.Interrupt, syscall.SIGTERM)
	go func() {
		readyToReceive <- true
		for msg := range msgs {
			communicationChan <- handler(msg.Body)
		}
	}()
	<-loop
	return nil
}

//ReceiveMessageOnTopicWithHeader is used to receive a specific message on a given topic. This method
// send back the receiving object trough a chan of interface. The receiving object is computed by
// func passed in argument and has to be declared like the following example :
// `func handler(msg amqp.Delivery) interface{}`
func (rConn *RabbitConnection) ReceiveMessageOnTopicWithHeader(topic string, handler func(amqp.Delivery) interface{}, communicationChan chan interface{}, readyToReceive chan bool) error {
	queue, err := rConn.partiesChannel.QueueDeclare(
		"",    // name
		false, // type
		false, // durable
		true,  // auto-deleted
		false, // internal
		nil,   // no-wait
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	err = rConn.partiesChannel.QueueBind(
		queue.Name,      //queue name
		topic,           //routing key
		"parties_topic", // exchange
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	msgs, err := rConn.partiesChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	loop := make(chan os.Signal)
	signal.Notify(loop, os.Interrupt, syscall.SIGTERM)
	logger.Trace("waiting message on topic :", topic)
	go func() {
		readyToReceive <- true
		for msg := range msgs {
			communicationChan <- handler(msg)
		}
	}()
	<-loop
	return nil
}

// ReceiveMessageOnTopicWithCallback is used to receive a specific message on a given topic and send back a message
// to the sender. The receiving object and the object sent is computed by func passed in argument and has to be
// declared like the following example :
// `func(delivery amqp.Delivery) (interface{}, string)`
// The string value in the returned tuple is added to the response topic.
func (rConn *RabbitConnection) ReceiveMessageOnTopicWithCallback(topic, responseTopic string, callback func(interface{}, string), responseCreator func(delivery amqp.Delivery) (interface{}, string), readyToReceive chan bool) error {
	queue, err := rConn.partiesChannel.QueueDeclare(
		"",    // name
		false, // type
		false, // durable
		true,  // auto-deleted
		false, // internal
		nil,   // no-wait
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	err = rConn.partiesChannel.QueueBind(
		queue.Name,      //queue name
		topic,           //routing key
		"parties_topic", // exchange
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}

	msgs, err := rConn.partiesChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	loop := make(chan os.Signal)
	signal.Notify(loop, os.Interrupt, syscall.SIGTERM)
	logger.Trace("waiting message on topic :", topic)
	go func() {
		readyToReceive <- true
		for msg := range msgs {
			var err error
			responseForged, forgedResponseTopic := responseCreator(msg)
			if forgedResponseTopic == "" {
				forgedResponseTopic = responseTopic
			}
			if []byte(forgedResponseTopic)[0] == byte('.') {
				forgedResponseTopic = responseTopic + forgedResponseTopic
			}
			if reflect.TypeOf(responseForged) == reflect.TypeOf(err) {
				logger.Error("error while calling response creator :", responseForged)
				errMessage := ErrorResponse{ErrorMessage: "error while calling response creator :" + responseForged.(error).Error()}
				callback(errMessage, responseTopic)
				continue
			}
			callback(responseForged, forgedResponseTopic)
		}
	}()
	<-loop
	return nil
}

// ReceiveMessageOnTopicWithHandler is used to receive a specific message on a given topic and handle it in a function.
// The handler function does not send any data. The func passed in argument and has to be declared like
// the following example :
// `handler func(delivery amqp.Delivery)`
func (rConn *RabbitConnection) ReceiveMessageOnTopicWithHandler(topic string, handler func(amqp.Delivery), readyToReceive chan bool) error {
	queue, err := rConn.partiesChannel.QueueDeclare(
		"",    // name
		false, // type
		false, // durable
		true,  // auto-deleted
		false, // internal
		nil,   // no-wait
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	err = rConn.partiesChannel.QueueBind(
		queue.Name,      //queue name
		topic,           //routing key
		"parties_topic", // exchange
		false,           // no wait
		nil,             // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}

	logger.Trace("waiting message on topic :", topic)
	msgs, err := rConn.partiesChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	if err != nil {
		readyToReceive <- false
		return err
	}
	loop := make(chan os.Signal)
	signal.Notify(loop, os.Interrupt, syscall.SIGTERM)
	go func() {
		readyToReceive <- true
		for msg := range msgs {
			handler(msg)
		}
	}()
	<-loop
	return nil
}

// Close terminate RabbitMQ connection
func (rConn *RabbitConnection) Close() error {
	err := rConn.lobbiesChannel.Close()
	if err != nil {
		return err
	}
	return rConn.partiesChannel.Close()
}
