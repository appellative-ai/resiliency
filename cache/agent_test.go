package cache

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"time"
)

func ExampleNew() {
	url := "https://www.google.com/search"
	a := newAgent(nil)

	fmt.Printf("test: newAgent() -> %v\n", a.Uri())

	a.Message(messaging.NewConfigMessage(url))
	time.Sleep(time.Second * 2)
	fmt.Printf("test: Message() -> %v\n", a.url)

	//Output:
	//test: newAgent() -> resiliency:agent/behavioral-ai/resiliency/cache
	//test: Message() -> https://www.google.com/search

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
