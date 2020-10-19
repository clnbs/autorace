package logger

import (
	"fmt"
	"runtime"
	"strings"
)

//ILog is a logger interface to implements new logger
type ILog interface {
	SetLogLevel(string)
	GetLogLevel() LogLevel
	Trace(...interface{})
	Debug(...interface{})
	Warning(...interface{})
	Error(...interface{})
}

//LogLevel express the log level to print. It got four states :
// Trace : print trace of computation and all below
// Debug : print debug output and all below
// Warning : print warning and error below
// Error : print only error
type LogLevel int

const (
	trace = iota
	debug
	warning
	errorl
	notSet
)

func (ll LogLevel) String() string {
	states := []string{"Trace", "Debug", "Warning", "Error", "Not set"}
	if ll < trace || ll > notSet {
		return "unknown log level"
	}
	return states[ll]
}

//NewLogLevel is LogLevel's constructor. It transforms a string in log level output.
//If string in argument does not correspond to any of the four state, it returns a trace level
func NewLogLevel(ll string) LogLevel {
	lowerCaseLL := strings.ToLower(ll)
	logLevelStates := make(map[string]LogLevel)
	logLevelStates["trace"] = 0
	logLevelStates["debug"] = 1
	logLevelStates["warning"] = 2
	logLevelStates["error"] = 3
	logLevelStates["not set"] = 4
	if _, ok := logLevelStates[lowerCaseLL]; !ok {
		return trace
	}
	return logLevelStates[lowerCaseLL]
}

//getCallingStack return a string of the calling stack in order to be printed by
// any logger who call it. It format it like the builtin golang logger
func getCallingStack() string {
	var fileAndLine string
	_, file, line, ok := runtime.Caller(3)
	if ok {
		files := strings.Split(file, "/")
		file = files[len(files)-1]
		fileAndLine = fmt.Sprintf("%s:%d :", file, line)
		return fileAndLine
	}
	return ""
}

var logs map[string]ILog

func init() {
	logs = make(map[string]ILog)
}

func Trace(msgs ...interface{}) {
	for _, l := range logs {
		l.Trace(msgs)
	}
}

func Debug(msgs ...interface{}) {
	for _, l := range logs {
		l.Debug(msgs)
	}
}

func Warning(msgs ...interface{}) {
	for _, l := range logs {
		l.Warning(msgs)
	}
}

func Error(msgs ...interface{}) {
	for _, l := range logs {
		l.Error(msgs)
	}
}

func SetLogLevel(logLvl string) {
	for _, l := range logs {
		l.SetLogLevel(logLvl)
	}
}

func GetLogLevel() LogLevel {
	if len(logs) > 0 {
		for _, l := range logs {
			return l.GetLogLevel()
		}
	}
	return notSet
}

func extractText(msgs ...interface{}) string {
	str := ""
	for _, v := range msgs {
		str = fmt.Sprintf("%v", v)
	}
	//str = strings.Trim(str, "]")
	//fmt.Println(str)
	str = strings.Replace(str, "[", "", -1)
	str = strings.Replace(str, "]", "", -1)
	return str
}