package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NameKey = "name"
	PathKey = "@path"
	RoleKey = "role"
)

var (
	Agent = New()
)

// ConfigureOrigin - map must provide region, zone, sub-zone, domain, collective, and service-name
func ConfigureOrigin(m map[string]string, read func() ([]byte, error)) error {
	var m2 = make(map[string]string)

	if read == nil {
		return errors.New("origin read function is nil")
	}
	// Read the origin JSON
	buf, err := read()
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

func ConfigureLogging(read func() ([]byte, error)) error {
	if read == nil {
		return errors.New("logging read function is nil")
	}
	return access.LoadOperators(func() ([]byte, error) {
		return read()
	})
}

// ConfigureNetworks - configure application networks
func ConfigureNetworks(appCfg map[string]string, read func(fileName string) ([]byte, error)) (errs []error) {
	if appCfg == nil {
		return []error{errors.New("application config is nil")}
	}
	//var result = make([]error, len(appCfg)*2)
	//var wg sync.WaitGroup

	var i int
	for k, v := range appCfg {
		if v == "" {
			errs = append(errs, errors.New(fmt.Sprintf("file name is empty for case officer: %v", k)))
			continue
		}
		officer, err1 := validateOfficerType(k)
		if err1 != nil {
			errs = append(errs, err1)
			continue
		}
		if officer == nil {
		}

		if i != 0 {
			i++
		}
		//wg.Add(1)
		//buildEndpoint(officer, v, read, &result[i])
	}
	//wg.Wait()
	// Need to create
	return packErrors(errs)
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
