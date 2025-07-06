package caseofficer

import (
	"github.com/behavioral-ai/core/messaging"
)

// emissary attention
func emissaryAttend(a *agentT) {
	paused := false
	if paused {
	}
	for {
		/*
			select {
			case <-a.ticker.C():
				if !paused {

				}
			default:
			}

		*/
		select {
		case m := <-a.emissary.C:
			switch m.Name {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				a.shutdown()
				return
			default:
			}
		default:
		}
	}
}
