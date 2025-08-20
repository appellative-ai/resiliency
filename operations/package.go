package operations

import (
	"errors"
	"fmt"
	"github.com/appellative-ai/agency/logger"
	"github.com/appellative-ai/agency/logx"
	"github.com/appellative-ai/agency/network"
	"github.com/appellative-ai/collective/exchange"
	cops "github.com/appellative-ai/collective/operations"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/rest"
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
func ConfigureOrigin(origin map[string]string) error {
	return cops.ConfigOrigin(origin)
}

// ConfigureRegistry - map must provide region, zone, sub-zone, domain, collective, and service-name
func ConfigureRegistry(hosts []string) error {
	return cops.ConfigRegistry(hosts)
}

// ConfigureLogging -
func ConfigureLogging(ops []logx.Operator) error {
	if len(ops) == 0 {
		return nil
	}
	newOps, err := logx.InitOperators(ops)
	if err != nil {
		return err
	}
	m := messaging.NewConfigMessage(newOps).AddTo(logger.AgentName)
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
		caseOfficer := agent.registerCaseOfficer(m[endpointKey])
		if m[testKey] == "true" {
			setTestOverrides(caseOfficer)
		}
		netCfg, err := network.BuildConfig(m[networkKey], read)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		operatives, errs1 := caseOfficer.BuildNetwork(netCfg)
		if errs1 != nil {
			errs = append(errs, errs1...)
			continue
		}
		if len(operatives) == 0 {
			errs = append(errs, errors.New(fmt.Sprintf("no operatives configured for network: %v", m[networkKey])))
			continue
		}
		Endpoint[m[endpointKey]] = NewEndpoint(m[patternKey], operatives)
	}
	return errs
}

// Startup - application
func Startup() error {
	// Start the collective first
	cops.ConfigLogging(logger.Agent.LogEgress)
	err := cops.Startup()
	if err != nil {
		return err
	}
	agent.Message(messaging.StartupMessage)
	return nil
}

// Shutdown -
// TODO: need to shutdown all global assigned agents
func Shutdown() {
	cops.Shutdown()
	agent.Message(messaging.ShutdownMessage)
}
