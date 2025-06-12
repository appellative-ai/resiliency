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
	OperationsPattern = "/operations"
	RootPattern       = "/"
	HealthPattern     = "/health"
)

var (
	Operations = newOperationsEndpoint()
	Root       = newRootEndpoint()
	Health     = newHealthEndpoint()
	Primary    *rest.Endpoint
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

		//cache.ConstructorOverride(nil, cachetest.Exchange, operations.Serve)
		//routing.ConstructorOverride(nil, routingtest.Exchange, operations.Serve)
		Primary = host.NewEndpoint(chain)
	default:
		return errors.New(fmt.Sprintf("error: agent not found for name: %v", name))
	}
	return nil
}
