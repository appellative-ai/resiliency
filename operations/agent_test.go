package operations

import (
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/collective/event/eventtest"
	"github.com/behavioral-ai/core/messaging"
	"time"
)

func _ExampleAgent_NotFound() {
	ch := make(chan struct{})
	agent := newAgent(nil)

	go func() {
		agent.Message(messaging.StartupMessage)
		time.Sleep(testDuration * 20)

		messaging.Shutdown(agent)
		time.Sleep(testDuration * 5)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}

func ExampleAgent() {
	ch := make(chan struct{})
	dispatcher := event.NewFilteredTraceDispatcher([]string{messaging.ResumeEvent, messaging.PauseEvent}, "")
	agent := newAgent(eventtest.New(dispatcher))

	go func() {
		agent.Message(messaging.StartupMessage)
		time.Sleep(testDuration * 6)
		agent.Message(messaging.PauseMessage)
		time.Sleep(testDuration * 6)
		agent.Message(messaging.ResumeMessage)
		time.Sleep(testDuration * 6)
		messaging.Shutdown(agent)
		time.Sleep(testDuration * 4)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
