package main

import (
	"errors"
	"os"
	"time"

	"github.com/clnbs/autorace/internal/app/server"
	"github.com/clnbs/autorace/pkg/logger"
)

var hitCounter = 5

func init() {
	logger.SetStdLogger("trace", "stdout")
	var err error
	err = errors.New("dummy")
	index := 0
	for err != nil && index < hitCounter {
		_, err = logger.SetFluentLogger("fluentd", "trace", "dynamic", 24224)
		if err != nil {
			time.Sleep(5 * time.Second)
		}
		index++
	}
	if err != nil {
		panic(err)
	}
}

func main() {
	readyToReceive := make(chan bool)
	srvr, err := server.NewDynamicPartyServer(os.Args[1], "rabbit")
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
