package network

import (
	"encoding/json"
	"errors"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
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
	access.Agent.SetOrigin(m2[messaging.RegionKey], m2[messaging.ZoneKey], m2[messaging.SubZoneKey], m2[messaging.HostKey], m2[messaging.InstanceIdKey])
	return nil
}
