package redirect

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/collective/event/eventtest"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	"github.com/behavioral-ai/resiliency/common"
	"time"
)

func ExampleNewAgent() {
	origin := common.Origin{Region: "us-central", Zone: "c-zone-a", SubZone: "sub-zone", Host: "www.host.com"}
	a := newAgent(common.Origin{}, nil)
	fmt.Printf("test: newAgent() -> [origin:%v] [uri:%v}\n", a.origin, a.Uri())

	a.Message(messaging.NewConfigMessage(origin))
	time.Sleep(time.Second * 2)
	fmt.Printf("test: Message() -> %v\n", a.origin)

	//Output:
	//test: newAgent() -> [origin:] [uri:resiliency:agent/behavioral-ai/resiliency/redirect1#}
	//test: Message() -> us-central.c-zone-a.sub-zone.www.host.com

}

func _ExampleAgent_LoadContent() {
	ch := make(chan struct{})
	dispatcher := event.NewTraceDispatcher()
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	s := messagingtest.NewTestSpanner(time.Second*2, testDuration)
	//test.LoadResiliencyContent()
	agent := newAgent(origin, eventtest.New(dispatcher))

	go func() {
		go masterAttend(agent, content.Resolver)
		go emissaryAttend(agent, content.Resolver, s)
		time.Sleep(testDuration * 5)

		agent.Message(messaging.ShutdownMessage)
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}

func _ExampleAgent_NotFound() {
	ch := make(chan struct{})
	dispatcher := event.NewTraceDispatcher()
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	agent := newAgent(origin, eventtest.New(dispatcher))

	go func() {
		agent.Run()
		time.Sleep(testDuration * 5)
		agent.Message(messaging.ShutdownMessage)
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}

func _ExampleAgent_Resolver() {
	ch := make(chan struct{})
	dispatcher := event.NewTraceDispatcher()
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	agent := newAgent(origin, eventtest.New(dispatcher))
	//test2.Startup()

	go func() {
		agent.Run()
		time.Sleep(testDuration * 5)
		agent.Message(messaging.ShutdownMessage)
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}
