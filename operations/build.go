package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/resiliency/caseofficer"
	"github.com/behavioral-ai/traffic/cache"
	"github.com/behavioral-ai/traffic/cache/cachetest"
	"github.com/behavioral-ai/traffic/routing"
	"github.com/behavioral-ai/traffic/routing/routingtest"
)

const (
	primaryPattern   = "/primary"
	secondaryPattern = "/secondary"
)

func buildEndpoint(name string, chain []any) error {
	switch name {
	case caseofficer.NamespaceNamePrimary:
		// In testing, need to override the Exchange for cache and routing
		m := rest.NewExchangeMessage(cachetest.Exchange)
		m.AddTo(cache.NamespaceName)
		exchange.Message(m)

		m = rest.NewExchangeMessage(routingtest.Exchange)
		m.AddTo(routing.NamespaceName)
		exchange.Message(m)

		Endpoint[PrimaryEndpoint] = host.NewEndpoint(primaryPattern, chain)
	case caseofficer.NamespaceNameSecondary:
		// In testing, need to override the Exchange for cache and routing
		//m := rest.NewExchangeMessage(cachetest.Exchange)
		//m.AddTo(caseofficer.NamespaceNamePrimary)
		//repository.Message(m)

		m := rest.NewExchangeMessage(routingtest.Exchange)
		m.AddTo(caseofficer.NamespaceNameSecondary)
		exchange.Message(m)

		Endpoint[SecondaryEndpoint] = host.NewEndpoint(secondaryPattern, chain)
	default:
		return errors.New(fmt.Sprintf("agent not found for name: %v", name))
	}
	return nil
}
