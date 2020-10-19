package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"testing"
	"time"
)

func TestNewRabbitConnection(t *testing.T) {

}

func byteToString(msg []byte) interface{} {
	return string(msg)
}

func TestRabbitConnection_ReceiveMessageOnTopic(t *testing.T) {
	rConn, err := NewRabbitConnection()
	if err != nil {
		t.Fatal("could not get a connection with RabbitMQ :", err)
	}
	done := make(chan bool)
	receive := make(chan interface{})
	rdyToReceive := make(chan bool)
	var receivedMsg string
	go func() {
		err = rConn.ReceiveMessageOnTopic("test.topic.receive", byteToString, receive, rdyToReceive)
	}()
	go func() {
		for receivedMsg != "\"end\"" {
			tmpMsg := <-receive
			receivedMsg = tmpMsg.(string)
			fmt.Println("received :", receivedMsg)
			fmt.Println(receivedMsg == "\"end\"")
		}
		done<-true
	}()
	<-rdyToReceive
	go func() {
		rConn.SendMessageOnTopic("message 1 on test.topic.receive", "test.topic.receive")
		rConn.SendMessageOnTopic("message 2 on test.topic.receive", "test.topic.receive")
		rConn.SendMessageOnTopic("message 3 on test.topic.receive", "test.topic.receive")
		rConn.SendMessageOnTopic("message 4 on test.topic.receive", "test.topic.receive")
		rConn.SendMessageOnTopic("message 5 on test.topic.receive", "test.topic.receive")
		rConn.SendMessageOnTopic("message 6 on test.topic.receive", "test.topic.receive")
		rConn.SendMessageOnTopic("end", "test.topic.receive")
	}()
	<-done
}

func stringHandler(msg amqp.Delivery) {
	fmt.Println(string(msg.Body))
}

func TestRabbitConnection_ReceiveMessageOnTopicWithHandler(t *testing.T) {
	rConn, err := NewRabbitConnection()
	if err != nil {
		t.Fatal("could not get a connection with RabbitMQ :", err)
	}
	rdyToReceive := make(chan bool)
	go func() {
		err := rConn.ReceiveMessageOnTopicWithHandler("test.topic.handler", stringHandler, rdyToReceive)
		if err != nil {
			t.Fatal("could not receive data with handler :", err)
		}
	}()
	<-rdyToReceive
	go func() {
		rConn.SendMessageOnTopic("message 1 on test.topic.handler", "test.topic.handler")
		rConn.SendMessageOnTopic("message 2 on test.topic.handler", "test.topic.handler")
		rConn.SendMessageOnTopic("message 3 on test.topic.handler", "test.topic.handler")
		rConn.SendMessageOnTopic("message 4 on test.topic.handler", "test.topic.handler")
		rConn.SendMessageOnTopic("message 5 on test.topic.handler", "test.topic.handler")
		rConn.SendMessageOnTopic("message 6 on test.topic.handler", "test.topic.handler")
		rConn.SendMessageOnTopic("end on test.topic.handler", "test.topic.handler")
	}()
	time.Sleep(1*time.Second)
}

func stringResponseCreator(msg amqp.Delivery) (interface{}, string) {
	clientID := new(string)
	err := json.Unmarshal(msg.Body, clientID)
	if err != nil {
		return err, ""
	}
	return "Callback received : " + string(msg.Body), "." + *clientID
}

func stringHandlerCallback(msg []byte) {
	fmt.Println("Response from callback : " + string(msg))
}

func computeTestCallbackResponse(msg []byte) interface{} {
	str := new(string)
	err := json.Unmarshal(msg, str)
	if err != nil {
		return err
	}
	return str
}

func TestRabbitConnection_ReceiveMessageOnTopicWithCallback(t *testing.T) {
	rConn, err := NewRabbitConnection()
	if err != nil {
		t.Fatal("could not get a connection with RabbitMQ :", err)
	}
	rdyToReceive := make(chan bool)
	received := make(chan interface{})
	// Mock started to listen to possible response
	// first step, starting callback listener with callback function
	clientID := uuid.New()
	go func() {
		err := rConn.ReceiveMessageOnTopicWithCallback("test.topic.callback.testing", "test.topic.callback", rConn.SendMessageOnTopic, stringResponseCreator, rdyToReceive)
		if err != nil {
			t.Fatal("could not receive data with callback :", err)
		}
	}()
	<-rdyToReceive
	// Starting a mock that fake a asset request
	// second step, start to listen to a possible response
	go func() {
		err := rConn.ReceiveMessageOnTopic("test.topic.callback." + clientID.String(), computeTestCallbackResponse, received ,rdyToReceive)
		if err != nil {
			t.Fatal("could not receive data with handler :", err)
		}
	}()
	<-rdyToReceive
	// Callback ready to receive
	// third step, sending from "request"
	go func() {
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
		rConn.SendMessageOnTopic(clientID.String(), "test.topic.callback.testing")
	}()
	for i := 0; i <  7; i++ {
		str := <-received
		fmt.Println(*str.(*string))
	}

	time.Sleep(2*time.Second)
}
