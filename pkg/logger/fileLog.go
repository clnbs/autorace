package logger

import (
	"log"
	"os"
)

//FileLogger struct is a implementation of ILog interface. It contains file printer,
// the four log printer, and the log level desired.
type FileLogger struct {
	path    string
	trace   *log.Logger
	debug   *log.Logger
	warning *log.Logger
	error   *log.Logger
	logLvl  LogLevel
	writer  *os.File
}


//SetFileLogger is use to initialize a FileLogger. It accept a log level and a file path as strings
// to simplify usage. It return error if it could not open desired file.
func SetFileLogger(logLvl string, path string) (ILog, error) {
	fileLog := FileLogger{}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fileLog, err
	}
	fileLog.logLvl = NewLogLevel(logLvl)
	fileLog.path = path
	fileLog.writer = f

	fileLog.error = log.New(f, "[ERROR]\t", log.LstdFlags)
	fileLog.warning = log.New(f, "[WARN]\t", log.LstdFlags)
	fileLog.debug = log.New(f, "[DEBUG]\t", log.LstdFlags)
	fileLog.trace = log.New(f, "[TRACE]\t", log.LstdFlags)
	logs["file_logger"] = fileLog
	return fileLog, nil
}

//SetLogLevel can be use to change log level at run time. It ignore error since file path
// is already tested and registered in constructor
func (fl FileLogger) SetLogLevel(logLvl string) {
	newFl := fl
	newFl.logLvl = NewLogLevel(logLvl)
	logs["file_logger"] = newFl
}

//GetLogLevel returns currently used log level
func (fl FileLogger) GetLogLevel() LogLevel {
	return fl.logLvl
}

//Trace is use to print trace log in the log file
func (fl FileLogger) Trace(msgs ...interface{}) {
	if fl.logLvl > trace {
		return
	}
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	fl.trace.Println(extractText(msgs))
}

//Debug is use to print trace log in the log file
func (fl FileLogger) Debug(msgs ...interface{}) {
	if fl.logLvl > debug {
		return
	}
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	fl.debug.Println(extractText(msgs))
}

//Warning is use to print trace log in the log file
func (fl FileLogger) Warning(msgs ...interface{}) {
	if fl.logLvl > warning {
		return
	}
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	fl.warning.Println(extractText(msgs))
}

//Error is use to print trace log in the log file
func (fl FileLogger) Error(msgs ...interface{}) {
	msgs = append([]interface{}{getCallingStack()}, msgs...)
	fl.error.Println(extractText(msgs))
}

//Close is use to close file writer.
func (fl FileLogger) Close() error {
	return fl.writer.Close()
}
