package operations

import (
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/collective/event/eventtest"
	"github.com/behavioral-ai/core/messaging"
	"time"
)

const (
	testDuration = time.Second * 4
)

func ExampleEmissary() {
	ch := make(chan struct{})
	dispatcher := event.NewTraceDispatcher()
	agent := newAgent(eventtest.New(dispatcher))

	go func() {
		go emissaryAttend(agent)
		time.Sleep(testDuration * 2)

		agent.Message(messaging.ShutdownMessage)
		time.Sleep(testDuration * 2)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
