package operations

import (
	"github.com/behavioral-ai/core/messaging"
	"time"
)

func _ExampleAgent_NotFound() {
	ch := make(chan struct{})
	agent := newAgent(messaging.Activity, messaging.Notify, messaging.NewTraceDispatcher())

	go func() {
		agent.Run()
		time.Sleep(testDuration * 20)

		agent.Shutdown()
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
	dispatcher := messaging.NewFilteredTraceDispatcher([]string{messaging.ResumeEvent, messaging.PauseEvent}, "")
	agent := newAgent(messaging.Activity, messaging.Notify, dispatcher) //content.NewEphemeralResolver(), messaging.NewTraceDispatcher())
	//test.Startup()

	go func() {
		agent.Run()
		time.Sleep(testDuration * 6)
		agent.Message(messaging.Pause)
		time.Sleep(testDuration * 6)
		agent.Message(messaging.Resume)
		time.Sleep(testDuration * 6)
		agent.Shutdown()
		time.Sleep(testDuration * 4)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
