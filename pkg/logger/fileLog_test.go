package logger

import (
	"os"
	"testing"
)

func TestFileLogger_SetLogLevel(t *testing.T) {
	fileLog, err := SetFileLogger("error", os.Getenv("GOPATH") + "/src/github.com/clnbs/autorace/log/log_test.txt")
	if err != nil {
		t.Fatal("error while setting up file logger :", err)
	}

	t.Log("log level set at init on file", fileLog.(FileLogger).path ,"at :", fileLog.GetLogLevel(), "level")
	Trace("Test file logger :", "should not work on \"Trace\"")
	Debug("Test file logger :", "should not work on \"Debug\"")
	Warning("Test file logger :", "should not work on \"Warning\"")
	Error("Test file logger :", "\"Error\" should work")
	SetLogLevel("warning")
	t.Log("a new log level on file ", fileLog.(FileLogger).path ," after a change :", fileLog.GetLogLevel())
	Trace("Test file logger :", "should not work on \"Trace\"")
	Debug("Test file logger :", "should not work on \"Debug\"")
	Warning("Test file logger :", "\"Warning\" should work")
	Error("Test file logger :", "\"Error\" should work")
	SetLogLevel("debug")
	t.Log("a new log level on file ", fileLog.(FileLogger).path ," after a change :", fileLog.GetLogLevel())
	Trace("Test file logger :", "should not work on \"Trace\"")
	Debug("Test file logger :", "\"Debug\" should work")
	Warning("Test file logger :", "\"Warning\" should work")
	Error("Test file logger :", "\"Error\" should work")
	SetLogLevel("trace")
	t.Log("a new log level on file ", fileLog.(FileLogger).path ," after a change :", fileLog.GetLogLevel())
	Trace("Test file logger :", "\"Trace\"Should work")
	Debug("Test file logger :", "\"Debug\" should work")
	Warning("Test file logger :", "\"Warning\" should work")
	Error("Test file logger :", "\"Error\" should work")
	SetLogLevel("Anything")
	t.Log("a new log level on file ", fileLog.(FileLogger).path ," after a change :", fileLog.GetLogLevel())
	Trace("Test file logger :", "\"Trace\"Should work")
	Debug("Test file logger :", "\"Debug\" should work")
	Warning("Test file logger :", "\"Warning\" should work")
	Error("Test file logger :", "\"Error\" should work")
}

func ExampleSetFileLogger() {
	_, err := SetFileLogger("trace", "./log.txt")
	if err != nil {
		panic(err)
	}
}

func ExampleFileLogger_Warning() {
	//SetFileLogger should be call in an init function
	_, err := SetFileLogger("trace", "./log.txt")
	if err != nil {
		panic(err)
	}
	Warning("my warning message")
	// Output : in file "log.txt" : [WARN]	2020/08/30 15:25:14 fileLog_test.go:48 : my warning message
}

func ExampleFileLogger_Close() {
	//SetFileLogger should be call in an init function
	fileLog, err := SetFileLogger("trace", "./log.txt")
	if err != nil {
		panic(err)
	}
	Trace("file logger loaded successfully")
	// Output : in file "log.txt" : [TRACE]	2020/08/30 15:25:14 fileLog_test.go:58 : file logger loaded successfully
	err = fileLog.(FileLogger).Close()
	if err != nil {
		panic(err)
	}
}
