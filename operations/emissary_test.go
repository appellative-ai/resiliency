package operations

import (
	"github.com/behavioral-ai/core/messaging"
	"time"
)

const (
	testDuration = time.Second * 4
)

func ExampleEmissary() {
	ch := make(chan struct{})
	traceDispatcher := messaging.NewTraceDispatcher()
	agent := newAgent(messaging.Activity, messaging.Notify, traceDispatcher)

	go func() {
		go emissaryAttend(agent)
		time.Sleep(testDuration * 2)

		agent.Shutdown()
		time.Sleep(testDuration * 2)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
