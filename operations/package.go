package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appellative-ai/agency/network"
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/core/host"
	"github.com/appellative-ai/core/logx"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/rest"
	"github.com/appellative-ai/core/std"
	"github.com/appellative-ai/traffic/logger"
)

const (
	serviceEndpoint = "service"
	healthEndpoint  = "health"

	endpointKey = "endpoint"
	patternKey  = "pattern"
	networkKey  = "network"
	testKey     = "test"
)

// Endpoint - HTTP endpoints
var (
	Endpoint = map[string]rest.Endpoint{
		serviceEndpoint: newServiceEndpoint("/operations/"),
		healthEndpoint:  newHealthEndpoint("/health/"),
	}
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
	//roles := []string{LoggingRole, AuthorizationRole, CacheRole, RateLimitingRole, RoutingRole}
	for _, m := range endpointCfg {
		if m[endpointKey] == "" {
			errs = append(errs, errors.New(fmt.Sprintf("endpoint name is empty")))
			continue
		}
		if m[networkKey] == "" {
			errs = append(errs, errors.New(fmt.Sprintf("network file name is empty for endpoint: %v", m[endpointKey])))
			continue
		}
		if m[patternKey] == "" {
			errs = append(errs, errors.New(fmt.Sprintf("pattern is empty for endpoint: %v", m[endpointKey])))
			continue
		}
		agent := opsAgent.registerCaseOfficer(m[endpointKey])
		if m[testKey] == "true" {
			setTestOverrides(agent)
		}
		netCfg, err := network.BuildConfig(m[networkKey], read)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		operatives, errs1 := agent.BuildNetwork(netCfg)
		if errs1 != nil {
			errs = append(errs, errs1...)
			continue
		}
		if len(operatives) == 0 {
			errs = append(errs, errors.New(fmt.Sprintf("no operatives configured for network: %v", m[networkKey])))
			continue
		}
		Endpoint[m[endpointKey]] = host.NewEndpoint(m[patternKey], operatives)
	}
	return errs
}

// ReadEndpointConfig -
func ReadEndpointConfig(read func() ([]byte, error)) ([]map[string]string, error) {
	return network.ReadEndpointConfig(read)
}

// Startup - application
func Startup() {
	opsAgent.Message(messaging.StartupMessage)
}

// Shutdown -
// TODO: need to shutdown all global assigned agents
func Shutdown() {
	opsAgent.Message(messaging.ShutdownMessage)
}
