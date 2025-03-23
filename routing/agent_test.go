package routing

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"time"
)

func ExampleNew() {
	a := newAgent(nil, "", 0)

	fmt.Printf("test: newAgent() -> %v\n", a.Uri())

	m := make(map[string]string)
	m[common.AppHostKey] = "google.com"
	a.Message(messaging.NewConfigMessage(m))
	time.Sleep(time.Second * 2)
	fmt.Printf("test: Message() -> %v\n", a.hostName)

	//Output:
	//test: newAgent() -> resiliency:agent/behavioral-ai/resiliency/routing
	//test: Message() -> google.com

}

/*
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
	dispatcher := event.NewFilteredTraceDispatcher([]string{messaging.ResumeEvent, messaging.PauseEvent}, "")
	agent := newAgent("",eventtest.New(dispatcher)) //content.NewEphemeralResolver(), messaging.NewTraceDispatcher())
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


*/
