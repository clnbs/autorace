package container

import "testing"

func TestExample(t *testing.T) {
	err := CreateDynamicServer("testingID")
	if err != nil {
		t.Fatal("error while creating a new dynamic server")
	}
}