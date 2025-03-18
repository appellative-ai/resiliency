package operations

import (
	"github.com/behavioral-ai/core/messaging"
)

// emissary attention
func emissaryAttend(agent *agentT) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false
	if paused {
	}

	for {
		select {
		case msg := <-agent.emissary.C:
			agent.dispatch(agent.emissary, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
				agent.agents.Broadcast(messaging.Pause)
			case messaging.ResumeEvent:
				paused = false
				agent.agents.Broadcast(messaging.Resume)
			case messaging.StopEvent:
			case messaging.StartEvent:
			case messaging.ShutdownEvent:
				agent.finalize()
				return
			default:
			}
		default:
		}
	}
}
