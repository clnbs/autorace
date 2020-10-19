package logger

import "testing"

func TestSetFluentLogger(t *testing.T) {
	fl, err := SetFluentLogger("localhost", "trace", "test", 24224)
	if err != nil {
		t.Fatal("error while creating fluent logger :", err)
	}
	fl.Trace("this is a test on trace ...")
	fl.Debug("this is a test on debug ...")
	fl.Warning("this is a test on warning ...")
	fl.Error("this is a test on error ...")
}