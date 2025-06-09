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

// ConfigureOrigin - map must provide region, zone, sub-zone, domain, collective, and service-name
func ConfigureOrigin(path string, m map[string]string, read func(string) ([]byte, error)) error {
	var m2 = make(map[string]string)

	if path == "" || read == nil {
		return errors.New("origin path is empty or read function is nil")
	}
	// Read the origin JSON
	buf, err := read(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &m2)
	if err != nil {
		return err
	}
	// Add host created items
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

func ConfigureLogging(path string, read func(string) ([]byte, error)) error {
	if path == "" || read == nil {
		return errors.New("logging path is empty or read function is nil")
	}
	return access.LoadOperators(func() ([]byte, error) {
		return read(path)
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
