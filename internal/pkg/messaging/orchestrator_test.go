package messaging

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrchestrator(t *testing.T) {
	o := NewOrchestrator()
	assert.NotNil(t, o, "orchestrator should not be nil")
}

func TestOrchestrator_SendMessage(t *testing.T) {
	o := NewOrchestrator()
	chan1 := make(chan interface{})
	chan2 := make(chan interface{})
	wg := sync.WaitGroup{}
	o.RegisterListener(&chan1)
	o.RegisterListener(&chan2)
	go func() {
		wg.Add(1)
		fmt.Println("from chan 1 :", <-chan1)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		fmt.Println("from chan 2 :", <-chan2)
		wg.Done()
	}()
	o.SendMessage("toto")
	wg.Wait()
}

func TestOrchestrator_RemoveListener(t *testing.T) {
	o := NewOrchestrator()
	chan1 := make(chan interface{})
	chan2 := make(chan interface{})
	o.RegisterListener(&chan1)
	o.RegisterListener(&chan2)
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		fmt.Println("from chan 1 :", <-chan1)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		fmt.Println("from chan 2 :", <-chan2)
		wg.Done()
	}()
	o.SendMessage("toto should be printed on all chan !")
	wg.Wait()

	o.RemoveListener(&chan1)
	go func() {
		wg.Add(1)
		fmt.Println("from chan 1 :", <-chan1)
		wg.Done()
	}()
	go func() {
		wg.Add(1)
		fmt.Println("from chan 2 :", <-chan2)
		wg.Done()
	}()
	o.SendMessage("tata, this one should be printed on chan 2 only")
	chan1 <- "shall you be free"
	wg.Wait()
}
