package logger

import (
	"log"
	"os"
	"strings"
)

//LogOutput express the Std logs should be printed on.
// It can either be StdOut or StdErr
type LogOutput int

const (
	stdout = iota
	stderr
)

func (logOut LogOutput) String() string {
	states := []string{"stdout", "stderr"}
	if logOut < 0 || logOut > 1 {
		return "unknown log output"
	}
	return states[logOut]
}

//StdLogger struct is a implement of ILog interface. It contains printer of the four log printer,
// the log output and the log level desired.
type StdLogger struct {
	trace     *log.Logger
	debug     *log.Logger
	warning   *log.Logger
	error     *log.Logger
	logLvl    LogLevel
	logOutput LogOutput
	output    *os.File
}

//NewLogOutput is LogOutput's constructor. It transforms a string in a log output on Std.
// If string in argument does not correspond to any of the two output, it return a StdOut output
func NewLogOutput(logOutInput string) LogOutput {
	lowerCaseLogOut := strings.ToLower(logOutInput)
	logOutputStates := make(map[string]LogOutput)
	logOutputStates["stdout"] = 0
	logOutputStates["stderr"] = 1
	if _, ok := logOutputStates[lowerCaseLogOut]; !ok {
		return stdout
	}
	return logOutputStates[lowerCaseLogOut]
}

//SetStdLogger is use to initialize a StdLogger. It accept a log level and an output as strings
// to simplify usage.
func SetStdLogger(logLvl string, logOutput string) ILog {
	stdLog := StdLogger{}
	stdLog.logLvl = NewLogLevel(logLvl)
	stdLog.logOutput = NewLogOutput(logOutput)

	if stdLog.logOutput == stdout {
		stdLog.output = os.Stdout
	}
	if stdLog.logOutput == stderr {
		stdLog.output = os.Stderr
	}

	stdLog.error = log.New(stdLog.output, "[ERROR]\t", log.LstdFlags)
	stdLog.warning = log.New(stdLog.output, "[WARN]\t", log.LstdFlags)
	stdLog.debug = log.New(stdLog.output, "[DEBUG]\t", log.LstdFlags)
	stdLog.trace = log.New(stdLog.output, "[TRACE]\t", log.LstdFlags)
	logs["std_logger"] = stdLog
	return stdLog
}

//SetLogLevel can be use to change log level at run time.
func (sl StdLogger) SetLogLevel(logLvl string) {
	newSl := sl
	newSl.logLvl = NewLogLevel(logLvl)
	logs["std_logger"] = newSl
}

//GetLogLevel returns currently used log level
func (sl StdLogger) GetLogLevel() LogLevel {
	return sl.logLvl
}

//Trace is use to print trace logs
func (sl StdLogger) Trace(msgs ...interface{}) {
	if sl.logLvl > trace {
		return
	}
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	sl.trace.Println(extractText(msgs))
}

//Debug is use to print debug logs
func (sl StdLogger) Debug(msgs ...interface{}) {
	if sl.logLvl > debug {
		return
	}
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	sl.debug.Println(extractText(msgs))
}

//Warning is use to print warning logs
func (sl StdLogger) Warning(msgs ...interface{}) {
	if sl.logLvl > warning {
		return
	}
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	sl.warning.Println(extractText(msgs))
}

//Error is use to print error logs
func (sl StdLogger) Error(msgs ...interface{}) {
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	sl.error.Println(extractText(msgs))
}
