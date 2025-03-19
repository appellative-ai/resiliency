package limit

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
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

	//agent := agentT{}
	//t := reflect.TypeOf(agent)
	//fmt.Printf("test: agenT -> [%v] [name:%v] [path:%v] [kind:%v]\n", t, t.Name(), t.PkgPath(), t.Kind())

	//t = reflect.TypeOf(New)
	//fmt.Printf("test: New() -> [%v] [name:%v] [path:%v] [kind:%v]\n", t, t.Name(), t.PkgPath(), t.Kind())

	//Output:
	//test: newAgent() -> [origin:] [uri:resiliency:agent/behavioral-ai/resiliency/rate-limiting1#}
	//test: Message() -> us-central.c-zone-a.sub-zone.www.host.com

}

/*
func ExampleAgent_LoadContent() {
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


*/
