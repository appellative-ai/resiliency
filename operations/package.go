package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/iox"
	"github.com/behavioral-ai/core/messaging"
)

var (
	Agent = New()
)

func ConfigureOrigin(path string, m map[string]string) error {
	var m2 = make(map[string]string)

	if path == "" {
		return errors.New("origin path is empty")
	}
	// Read the origin JSON
	buf, err := iox.ReadFile(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &m2)
	if err != nil {
		return err
	}
	// Add region, zone, sub-zone, domain, collective, service-name from map m
	for k, v := range m {
		m2[k] = v
	}
	status := messaging.SetOrigin(m2)
	if !status.OK() {
		return status.Err
	}
	access.SetOrigin(m2[messaging.RegionKey], m2[messaging.ZoneKey], m2[messaging.SubZoneKey], m2[messaging.HostKey], m2[messaging.InstanceIdKey])
	return nil
}

func ConfigureLogging(path string) error {
	if path == "" {
		return errors.New("logging operator path is empty")
	}
	return access.LoadOperators(func() ([]byte, error) {
		return iox.ReadFile(path)
	})
}

// ConfigureAgents - configure all agents
// TODO : add configuration for caching profile
func ConfigureAgents(mapPath, profilePath string) error {
	if mapPath != "" {
		m, err := iox.ReadMap(mapPath)
		if err != nil {
			return err
		}
		msg := messaging.NewMapMessage(m)
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
