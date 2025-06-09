package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/caseofficer"
	"strings"
	"sync"
)

const (
	NameKey = "name"
	PathKey = "@path"
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
	var result = make([]error, len(appCfg)*2)
	var wg sync.WaitGroup
	var i int
	for k, v := range appCfg {
		if v == "" {
			errs = append(errs, errors.New(fmt.Sprintf("value is empty for case officer: %v", k)))
			continue
		}
		m := parseOfficerConfig(v)
		name, fileName, err := validateOfficerConfig(v, m)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		officer, err1 := validateOfficerType(name)
		if err1 != nil {
			errs = append(errs, err1)
			continue
		}

		if i != 0 {
			i++
		}
		wg.Add(1)
		go func(officer caseofficer.Agent, fileName string, read func(fileName string) ([]byte, error), err *error) {
			defer wg.Done()
			var networkCfg []map[string]string

			buf, err2 := read(fileName)
			if err2 != nil {
				*err = err2
				return
			}
			err2 = json.Unmarshal(buf, networkCfg)
			if err2 != nil {
				*err = err2
				return
			}
		}(officer, fileName, read, &result[i])

	}
	wg.Wait()
	// Need to create
	return packErrors(errs)

}

func parseOfficerConfig(s string) map[string]string {
	var m = make(map[string]string)

	tokens := strings.Split(s, ",")
	for _, t := range tokens {
		pairs := strings.Split(t, "=")
		if len(pairs) < 2 || pairs[1] == "" {
			continue
		}
		m[pairs[0]] = pairs[1]
	}
	return m
}

func validateOfficerConfig(caseOfficer string, m map[string]string) (name string, fileName string, err error) {
	var ok bool
	name, ok = m[NameKey]
	if name == "" || !ok {
		return "", "", errors.New(fmt.Sprintf("name is empty for case officer: %v", caseOfficer))
	}
	fileName, ok = m[PathKey]
	if fileName == "" || !ok {
		return "", "", errors.New(fmt.Sprintf("file name is empty for case officer: %v", caseOfficer))
	}
	return
}

func validateOfficerType(name string) (caseofficer.Agent, error) {
	agent := repository.Agent(name)
	if agent == nil {
		return nil, errors.New(fmt.Sprintf("agent is nil for case officer: %v", name))
	}
	officer, ok := any(agent).(caseofficer.Agent)
	if !ok {
		return nil, errors.New(fmt.Sprintf("agent is no of type caseofficer.Agent case officer: %v", name))
	}
	return officer, nil
}

func packErrors(errs []error) []error {
	var result []error
	for _, err := range errs {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
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
