package operations

import (
	"errors"
	"fmt"
	cops "github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
)

var (
	Agent = New()
)

/*
func ConfigureEventing(notifier eventing.NotifyFunc, activity eventing.ActivityFunc) {
	if notifier != nil {
		eventing.Handler.Message(eventing.NewNotifyConfigMessage(notifier))
	}
	if activity != nil {
		eventing.Handler.Message(eventing.NewActivityConfigMessage(activity))
	}
}


*/

func ConfigureLogging(operatorsPath, originPath string) error {
	if originPath != "" {
		m, err := iox.ReadMap(originPath)
		if err != nil {
			return err
		}
		//o, err1 := originFromMap(m)
		//if err1 != nil {
		//	return err1
		//}
		access.SetOrigin(m[cops.RegionKey], m[cops.ZoneKey], m[cops.SubZoneKey], m[cops.HostKey], m[cops.InstanceIdKey])
	}
	if operatorsPath != "" {
		err := access.LoadOperators(func() ([]byte, error) {
			return iox.ReadFile(operatorsPath)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// ConfigureAgents - configure all agents
// TODO : add configuration for caching profile
func ConfigureAgents(mapPath, profilePath string) error {
	if mapPath != "" {
		m, err := iox.ReadMap(mapPath)
		if err != nil {
			return err
		}
		msg := messaging.NewConfigMapMessage(m)
		repository.Broadcast(msg)
	}
	if profilePath != "" {
	}
	return nil

}

// Message - operations agent messaging
func Message(event string) error {
	switch event {
	case messaging.StartupEvent:
		Agent.Message(messaging.StartupMessage)
	case messaging.ShutdownEvent:
		Agent.Message(messaging.ShutdownMessage)
	case messaging.PauseEvent:
		Agent.Message(messaging.PauseMessage)
	case messaging.ResumeEvent:
		Agent.Message(messaging.ResumeMessage)
	default:
		return errors.New(fmt.Sprintf("operations.Message() -> [%v] [%v]", "error: invalid event", event))
	}
	return nil
}
