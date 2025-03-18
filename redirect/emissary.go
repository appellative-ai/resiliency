package redirect

import (
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/timeseries1"
)

// emissary attention
func emissaryAttend(agent *agentT, observe *timeseries1.Observer, resolver *content.Resolution, s messaging.Spanner) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false
	agent.reviseTicker(resolver, s)

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(agent.ticker, messaging.ObservationEvent)
			if !paused {
				e, status := observe.Timeseries(agent.origin)
				if status.OK() {
					m := messaging.NewMessage(messaging.Master, messaging.ObservationEvent)
					m.SetContent(contentTypeObservation, observation{origin: e.Origin, latency: e.Latency, gradient: e.Gradient})
					agent.Message(m)
				} else {
					status.AgentUri = agent.Uri()
					agent.notify(status)
				}
				agent.reviseTicker(resolver, s)
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.dispatch(agent.emissary, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				agent.emissaryFinalize()
				return
			default:
			}
		default:
		}
	}
}
