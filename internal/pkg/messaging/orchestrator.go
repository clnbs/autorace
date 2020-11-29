package messaging

type Orchestrator struct {
	sending map[*chan interface{}]interface{}
}

func NewOrchestrator() *Orchestrator {
	orchestrator := new(Orchestrator)
	orchestrator.sending = make(map[*chan interface{}]interface{}, 0)
	return orchestrator
}

func (orchestrator *Orchestrator) RegisterListener(c *chan interface{}) {
	orchestrator.sending[c] = nil
}

func (orchestrator *Orchestrator) SendMessage(message interface{}) {
	for c := range orchestrator.sending {
		*c <- message
	}
}

func (orchestrator *Orchestrator) RemoveListener(c *chan interface{}) {
	delete(orchestrator.sending, c)
}
