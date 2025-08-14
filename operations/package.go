package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appellative-ai/agency/network"
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/core/logx"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/rest"
	"github.com/appellative-ai/core/std"
	"github.com/appellative-ai/traffic/logger"
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
	status := std.SetOrigin(m2)
	if !status.OK() {
		return status.Err
	}
	return nil
}

// ConfigureLogging -
func ConfigureLogging(read func() ([]byte, error)) error {
	if read == nil {
		return errors.New("logging operators read function is nil")
	}
	buf, err := read()
	if err != nil {
		return err
	}
	var ops []logx.Operator

	err = json.Unmarshal(buf, &ops)
	if err != nil {
		return err
	}
	m := messaging.NewConfigMessage(ops).AddTo(logger.AgentName)
	exchange.Message(m)
	return nil
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
	if len(errs) == 0 {
		setTestOverrides()
	}
	return errs
}

// ReadEndpointConfig -
func ReadEndpointConfig(read func() ([]byte, error)) ([]map[string]string, error) {
	return network.ReadEndpointConfig(read)
}

// Http endpoints

const (
	ServiceEndpoint = "service"
	HealthEndpoint  = "health"
)

var (
	Endpoint = map[string]rest.Endpoint{
		ServiceEndpoint: newServiceEndpoint("/operations/"),
		HealthEndpoint:  newHealthEndpoint("/health/"),
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
