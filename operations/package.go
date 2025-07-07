package operations

import (
	"errors"
	"fmt"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/resiliency/network"
)

const (
	LoggingRole       = "logging"
	AuthorizationRole = "authorization"
	CacheRole         = "cache"
	RateLimitingRole  = "rate-limiting"
	RoutingRole       = "routing"

	endpointKey = "endpoint"
	patternKey  = "pattern"
	networkKey  = "network"
)

// ConfigureOrigin - map must provide region, zone, sub-zone, domain, collective, and service-name
func ConfigureOrigin(m map[string]string, read func() ([]byte, error)) error {
	return network.ConfigureOrigin(m, read)
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
func ConfigureNetworks(endpointCfg []map[string]string, read func(fileName string) ([]byte, error)) (errs []error) {
	if read == nil {
		return []error{errors.New("network configuration read function is nil")}
	}
	if len(endpointCfg) == 0 {
		return []error{errors.New("endpoint configuration is nil or empty")}
	}
	cfg := network.ShapeConfig(endpointKey, endpointCfg)
	roles := []string{LoggingRole, AuthorizationRole, CacheRole, RateLimitingRole, RoutingRole}
	for k, v := range cfg {
		if k == "" {
			errs = append(errs, errors.New(fmt.Sprintf("endpoint name is empty")))
			continue
		}
		if v[networkKey] == "" {
			errs = append(errs, errors.New(fmt.Sprintf("network file name is empty for endpoint: %v", k)))
			continue
		}
		agent := opsAgent.registerCaseOfficer(k)
		netCfg, err := network.BuildConfig(v[networkKey], read)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		chain, errs1 := agent.BuildNetwork(netCfg, roles)
		if errs1 != nil {
			errs = append(errs, errs1...)
			continue
		}
		err = buildEndpoint(k, v, chain)
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return errs
}

// Http endpoints

const (
	ServiceEndpoint = "service"
	HealthEndpoint  = "health"
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

// Shutdown -
// TODO: need to shutdown all global assigned agents
func Shutdown() {
	opsAgent.Message(messaging.ShutdownMessage)
}
