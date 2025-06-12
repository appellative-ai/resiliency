package endpoint

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/resiliency/caseofficer"
	"github.com/behavioral-ai/traffic/cache/cachetest"
	"github.com/behavioral-ai/traffic/routing/routingtest"
)

const (
	operationsPattern = "/operations"
	primaryPattern    = "/"
	healthPattern     = "/health"
)

//Root       = newRootEndpoint(rootPattern)

var (
	Operations = newOperationsEndpoint(operationsPattern)
	Health     = newHealthEndpoint(healthPattern)

	Primary *rest.Endpoint
)

func Build(name string, chain []any) error {
	switch name {
	case caseofficer.NamespaceNamePrimary:
		// In testing, need to override the Exchange for cache and routing
		m := rest.NewExchangeMessage(cachetest.Exchange)
		m.AddTo(caseofficer.NamespaceNamePrimary)
		repository.Message(m)

		m = rest.NewExchangeMessage(routingtest.Exchange)
		m.AddTo(caseofficer.NamespaceNamePrimary)
		repository.Message(m)

		Primary = host.NewEndpoint(primaryPattern, chain)
	default:
		return errors.New(fmt.Sprintf("error: agent not found for name: %v", name))
	}
	return nil
}
