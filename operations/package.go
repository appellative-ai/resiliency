package operations

import (
	"errors"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/rest"
	"net/http"
)

const (
	RoleKey = "role"
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

type EndpointT struct {
	Service http.Handler
	Health  http.Handler
	Primary *rest.Endpoint
}

var (
	Endpoint = EndpointT{
		Service: newServiceEndpoint("/operations"),
		Health:  newHealthEndpoint("/health"),
		Primary: nil,
	}
)
