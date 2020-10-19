package logger

import (
	"fmt"
	"github.com/fluent/fluent-logger-golang/fluent"
)

// FluentLogger struct is a implementation of ILog interface. It contains a Fluentd
// connection, the app's name using it and a log level
//TODO get config from conf file
type FluentLogger struct {
	fluent *fluent.Fluent
	logLvl LogLevel
	app string
}

// SetFluentLogger is use to initialize a FluentLogger. It accept a Fluentd address,
// a log level, an app name and a connection port. It returns an error if Fluentd is
// not reachable
func SetFluentLogger(host, logLevel, app string, port int) (ILog, error) {
	fluentLog := new(FluentLogger)
	var err error
	fluentLog.fluent, err = fluent.New(fluent.Config{
		FluentPort:         port,
		FluentHost:         host,
	})
	if err != nil {
		return nil, err
	}
	fluentLog.logLvl = NewLogLevel(logLevel)
	fluentLog.app = app
	logs["fluent_logger"] = fluentLog
	return fluentLog, nil
}

//SetLogLevel can be use to change log level at run time. It ignore error since file path
// is already tested and registered in constructor
func (fl *FluentLogger) SetLogLevel(logLevel string) {
	newFl := fl
	fl.logLvl = NewLogLevel(logLevel)
	logs["fluent_logger"] = newFl
}

//GetLogLevel returns currently used log level
func (fl *FluentLogger) GetLogLevel() LogLevel {
	return fl.logLvl
}

//Trace is use to send trace log to fluentd
func (fl *FluentLogger) Trace(msgs ...interface{}) {
	if fl.logLvl > trace {
		return
	}
	data := make(map[string]string)
	data["message"] = extractText(msgs)
	data["stack"] = getCallingStack()
	data["level"] = "trace"
	fl.sendLog(fl.app, data)
}

//Debug is use to send debug log to fluentd
func (fl *FluentLogger) Debug(msgs ...interface{}) {
	if fl.logLvl > debug {
		return
	}
	data := make(map[string]string)
	data["message"] = extractText(msgs)
	data["stack"] = getCallingStack()
	data["level"] = "debug"
	fl.sendLog(fl.app, data)
}

//Warning is use to send warning log to fluentd
func (fl *FluentLogger) Warning(msgs ...interface{}) {
	if fl.logLvl > warning {
		return
	}
	data := make(map[string]string)
	data["message"] = extractText(msgs)
	data["stack"] = getCallingStack()
	data["level"] = "warning"
	fl.sendLog(fl.app, data)
}

//Error is use to send error log to fluentd
func (fl *FluentLogger) Error(msgs ...interface{}) {
	data := make(map[string]string)
	data["message"] = extractText(msgs)
	data["stack"] = getCallingStack()
	data["level"] = "error"
	fl.sendLog(fl.app, data)
}

func (fl *FluentLogger) sendLog(tag string, data map[string]string) {
	err := fl.fluent.Post(tag, data)
	if err != nil {
		fmt.Println("error while sending data to Fluentd :", err)
	}
}