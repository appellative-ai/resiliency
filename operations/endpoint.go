package operations

import (
	"errors"
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/core/host"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/rest"
	"github.com/appellative-ai/traffic/cache"
	"github.com/appellative-ai/traffic/cache/cachetest"
	"github.com/appellative-ai/traffic/routing"
	"github.com/appellative-ai/traffic/routing/routingtest"
)

func buildEndpoint(name string, m map[string]string, chain []any) error {
	if len(m) == 0 {
		return errors.New("nil map")
	}
	Endpoint[name] = host.NewEndpoint(m[patternKey], chain)

	/*
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
			// In testing, need to override the Exchange for routing
			m := rest.NewExchangeMessage(routingtest.Exchange)
			m.AddTo(routing.NamespaceName)
			m.SetCareOf(caseofficer.NamespaceNameSecondary)
			exchange.Message(m)

			Endpoint[SecondaryEndpoint] = host.NewEndpoint(secondaryPattern, chain)
		default:
			//return errors.New(fmt.Sprintf("agent not found for name: %v", name))
		}

	*/
	return nil
}

func setTestOverrides() {
	agent1 := opsAgent.Operative("core:common:agent/caseofficer/request/http/primary")
	agent2 := opsAgent.Operative("core:common:agent/caseofficer/request/http/secondary")

	// Cache overrides - global
	m := messaging.NewConfigMessage(rest.Exchange(cachetest.Exchange))
	m.AddTo(cache.NamespaceName)
	exchange.Message(m)

	// Cache overrides - Local
	agent1.Message(m)
	agent2.Message(m)

	// Routing overrides - global
	m = messaging.NewConfigMessage(rest.Exchange(routingtest.Exchange))
	m.AddTo(routing.NamespaceName)
	exchange.Message(m)

	// Routing overrides - local
	agent1.Message(m)
	agent2.Message(m)

}
