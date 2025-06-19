package operations

import (
	"encoding/json"
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

func ReadAppConfig(read func() ([]byte, error)) (map[string]string, error) {
	var appCfg map[string]string

	buf, err := read()
	if err != nil {
		return nil, err //fmt.Printf("test: readFile(\"%v\") -> [bytes:%v] [err:%v]\n", subDir+appFileName, len(buf), err)
	}
	err = json.Unmarshal(buf, &appCfg)
	if err != nil {
		return nil, err //fmt.Printf("test: json.Unmarshal() -> [err:%v]\n", err)
	}
	return appCfg, nil
}

// Http endpoints

const (
	ServiceEndpoint   = "service"
	HealthEndpoint    = "health"
	PrimaryEndpoint   = "primary"
	SecondaryEndpoint = "secondary"
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
