package invoke

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	"github.com/behavioral-ai/resiliency/common"
	"time"
)

func ExampleNewAgent() {
	a := New(common.Origin{Region: "us-central", Zone: "c-zone-a", SubZone: "sub-zone", Host: "www.host.com"}, nil, nil, nil)

	//agent := agentT{}
	//t := reflect.TypeOf(agent)
	//fmt.Printf("test: agenT -> [%v] [name:%v] [path:%v] [kind:%v]\n", t, t.Name(), t.PkgPath(), t.Kind())

	//t = reflect.TypeOf(New)
	//fmt.Printf("test: New() -> [%v] [name:%v] [path:%v] [kind:%v]\n", t, t.Name(), t.PkgPath(), t.Kind())

	fmt.Printf("test: NewAgent() -> [%v]\n", a)
	fmt.Printf("test: NewAgent() -> [%v]\n", a.Name())

	//Output:
	//test: NewAgent() -> [resiliency:agent/behavioral-ai/operative1#us-central.c-zone-a.sub-zone.www.host.com]
	//test: NewAgent() -> [resiliency:agent/behavioral-ai/operative]

}

func ExampleAgent_LoadContent() {
	ch := make(chan struct{})
	dispatcher := messaging.NewTraceDispatcher()
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	s := messagingtest.NewTestSpanner(time.Second*2, testDuration)
	//test.LoadResiliencyContent()
	agent := newAgent(origin, messaging.Activity, messaging.Notify, dispatcher)

	go func() {
		go masterAttend(agent, content.Resolver)
		go emissaryAttend(agent, content.Resolver, s)
		time.Sleep(testDuration * 5)

		agent.Shutdown()
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}

func ExampleAgent_NotFound() {
	ch := make(chan struct{})
	dispatcher := messaging.NewTraceDispatcher()
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	agent := newAgent(origin, messaging.Activity, messaging.Notify, dispatcher)

	go func() {
		agent.Run()
		time.Sleep(testDuration * 5)
		agent.Shutdown()
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
	dispatcher := messaging.NewTraceDispatcher()
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	agent := newAgent(origin, messaging.Activity, messaging.Notify, dispatcher)
	//test2.Startup()

	go func() {
		agent.Run()
		time.Sleep(testDuration * 5)
		agent.Shutdown()
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}
