package main

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
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
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	readyToReceive := make(chan bool)
	srvr, err := server.NewStaticServer("rabbit")
	if err != nil {
		logger.Error("error while creation server :", err)
		return
	}
	logger.Trace("server created")
	go func() {
		err := srvr.ReceivePlayerCreation(readyToReceive)
		if err != nil {
			logger.Error("while listening to player creation :", err)
			return
		}
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	go func() {
		err := srvr.PartyListRequest(readyToReceive)
		if err != nil {
			logger.Error("while listening to party list request :", err)
			return
		}
		logger.Trace("end of ReceivePartyCreation")
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	go func() {
		err := srvr.ReceivePartyCreation(readyToReceive)
		if err != nil {
			logger.Error("while listening to party creation :", err)
			return
		}
		logger.Trace("end of PartyListRequest")
	}()
	if !<-readyToReceive {
		logger.Error("server could not listen continuously")
		return
	}
	logger.Trace("static server started ...")
	<-stop
	err = srvr.Close()
	if err != nil {
		logger.Error("while closing ongoing connection :", err)
	}
}
