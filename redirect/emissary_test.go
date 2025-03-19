package redirect

import (
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/collective/event/eventtest"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	"github.com/behavioral-ai/resiliency/common"
	"time"
)

const (
	testDuration = time.Second * 5
)

func ExampleEmissary() {
	ch := make(chan struct{})
	s := messagingtest.NewTestSpanner(time.Second*2, testDuration)
	agent := newAgent(common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}, eventtest.New(event.NewTraceDispatcher()))

	go func() {
		go emissaryAttend(agent, content.Resolver, s)
		agent.Message(messaging.NewMessage(messaging.Emissary, messaging.DataChangeEvent))
		time.Sleep(testDuration * 2)
		agent.Message(messaging.NewMessage(messaging.Emissary, messaging.PauseEvent))
		time.Sleep(testDuration * 2)
		agent.Message(messaging.NewMessage(messaging.Emissary, messaging.ResumeEvent))
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

func ExampleEmissary_Observation() {
	ch := make(chan struct{})
	s := messagingtest.NewTestSpanner(testDuration, testDuration)
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneB}
	agent := newAgent(origin, eventtest.New(event.NewTraceDispatcher()))

	go func() {
		go emissaryAttend(agent, content.Resolver, s)
		time.Sleep(testDuration * 2)

		agent.Message(messaging.ShutdownMessage)
		time.Sleep(testDuration * 3)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}
