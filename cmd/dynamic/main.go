package main

import (
	"errors"
	"github.com/clnbs/autorace/internal/pkg/messaging"
	"os"
	"strconv"
	"time"

	"github.com/clnbs/autorace/internal/app/server"
	"github.com/clnbs/autorace/pkg/logger"
)

var (
	hitCounter     = 5
	rabbitMQConfig messaging.RabbitConnectionConfiguration
)

func init() {
	var err error
	fluentdAddr := os.Getenv("FLUENTD_HOST")
	stringifyFluentdPort := os.Getenv("FLUENTD_PORT")
	logLevel := os.Getenv("LOG_LEVEL")
	fluentdPort, err := strconv.ParseInt(stringifyFluentdPort, 10, 64)
	if err != nil {
		panic(err)
	}
	logger.SetStdLogger(logLevel, "stdout")
	err = errors.New("dummy")
	index := 0
	for err != nil && index < hitCounter {
		_, err = logger.SetFluentLogger(fluentdAddr, logLevel, "dynamic", int(fluentdPort))
		if err != nil {
			time.Sleep(5 * time.Second)
		}
		index++
	}
	if err != nil {
		panic(err)
	}
	rabbitMQConfig = messaging.RabbitConnectionConfiguration{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASS"),
	}
}

func main() {
	readyToReceive := make(chan bool)
	srvr, err := server.NewDynamicPartyServer(os.Args[1], rabbitMQConfig)
	if err != nil {
		logger.Error("while creating dynamic server :", err)
		return
	}
	go func() {
		err := srvr.ReceiveAddPlayer(readyToReceive)
		if err != nil {
			logger.Error("while listening to adding player in party :", err)
			return
		}
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	go func() {
		err := srvr.ReceiveNewState(readyToReceive)
		if err != nil {
			logger.Error("while listening to new state :", err)
			return
		}
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	go func() {
		err := srvr.ReceiveSyncRequest(readyToReceive)
		if err != nil {
			logger.Error("while listening to sync request :", err)
			return
		}
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	go func() {
		err := srvr.ReceivePlayersInput(readyToReceive)
		if err != nil {
			logger.Error("while listening to player inputs :", err)
			return
		}
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	srvr.Run()
	err = srvr.Close()
	if err != nil {
		logger.Error("while closing ongoing connection :", err)
	}
}
