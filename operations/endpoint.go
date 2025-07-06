package operations

import (
	"errors"
	"github.com/behavioral-ai/core/host"
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
