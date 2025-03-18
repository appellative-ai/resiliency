package limit

import (
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/test"
	"time"
)

func ExampleMaster() {
	ch := make(chan struct{})
	agent := newAgent(common.Origin{Region: common.WestRegion}, messaging.Activity, messaging.Notify, messaging.NewTraceDispatcher())

	go func() {
		go masterAttend(agent, content.Resolver)
		agent.Message(messaging.NewMessage(messaging.Master, messaging.ObservationEvent))

		agent.Message(messaging.NewMessage(messaging.Master, messaging.PauseEvent))
		agent.Message(messaging.NewMessage(messaging.Master, messaging.ObservationEvent))
		agent.Message(messaging.NewMessage(messaging.Master, messaging.ResumeEvent))
		agent.Message(messaging.NewMessage(messaging.Master, messaging.ObservationEvent))

		agent.Shutdown()
		time.Sleep(testDuration)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)
	//Output:
	//fail
}

func ExampleMaster_Observation() {
	ch := make(chan struct{})
	origin := common.Origin{Region: common.WestRegion}
	msg := messaging.NewMessage(messaging.Master, messaging.ObservationEvent)
	msg.SetContent(contentTypeObservation, observation{origin: origin, latency: 2350, gradient: 15})
	test.LoadResiliencyContent()
	//resolver, status := test.NewResiliencyResolver()
	//if !status.OK() {
	//	messaging.Notify(status)
	//}
	agent := newAgent(origin, messaging.Activity, messaging.Notify, messaging.NewTraceDispatcher())

	go func() {
		go masterAttend(agent, content.Resolver)
		agent.Message(msg)
		time.Sleep(testDuration * 2)

		agent.Shutdown()
		time.Sleep(testDuration)

		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
