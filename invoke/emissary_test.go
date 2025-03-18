package invoke

import (
	"github.com/behavioral-ai/collective/content"
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
	agent := newAgent(common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}, messaging.Activity, messaging.Notify, messaging.NewTraceDispatcher())

	go func() {
		go emissaryAttend(agent, content.Resolver, s)
		agent.Message(messaging.NewMessage(messaging.Emissary, messaging.DataChangeEvent))
		time.Sleep(testDuration * 2)
		agent.Message(messaging.NewMessage(messaging.Emissary, messaging.PauseEvent))
		time.Sleep(testDuration * 2)
		agent.Message(messaging.NewMessage(messaging.Emissary, messaging.ResumeEvent))
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

func ExampleEmissary_Observation() {
	ch := make(chan struct{})
	s := messagingtest.NewTestSpanner(testDuration, testDuration)
	origin := common.Origin{Region: common.WestRegion, Zone: common.WestZoneB}
	agent := newAgent(origin, messaging.Activity, messaging.Notify, messaging.NewTraceDispatcher())

	go func() {
		go emissaryAttend(agent, content.Resolver, s)
		time.Sleep(testDuration * 2)

		// Receive observation message
		/*
			msg := <-agent.master.C
			o, status := getObservation(msg)
			status.AgentUri = agent.Uri()
			status.Msg = o.String()
			agent.notify(status)

		*/

		agent.Shutdown()
		time.Sleep(testDuration * 3)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail
}
