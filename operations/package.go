package operations

import (
	"errors"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
)

// ConfigureOrigin - map must provide region, zone, sub-zone, domain, collective, and service-name
func ConfigureOrigin(m map[string]string, read func() ([]byte, error)) error {
	return configureOrigin(m, read)
}

func ConfigureLogging(read func() ([]byte, error)) error {
	if read == nil {
		return errors.New("logging operators read function is nil")
	}
	return access.Agent.ConfigureOperators(func() ([]byte, error) {
		return read()
	})
}

// ConfigureNetworks - configure application networks
func ConfigureNetworks(appCfg map[string]string, read func(fileName string) ([]byte, error)) (errs []error) {
	return configureNetworks(appCfg, read)
}

// Http endpoints

const (
	ServiceEndpoint = "service"
	HealthEndpoint  = "health"
	PrimaryEndpoint = "primary"
)

var (
	Endpoint = map[string]rest.Endpoint{
		ServiceEndpoint: newServiceEndpoint("/operations"),
		HealthEndpoint:  newHealthEndpoint("/health"),
	}
)

// Startup - application
func Startup() {
	opsAgent.Message(messaging.StartupMessage)
}
