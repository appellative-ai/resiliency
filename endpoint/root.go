package endpoint

import (
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/traffic/cache"
	"github.com/behavioral-ai/traffic/cache/cachetest"
	"github.com/behavioral-ai/traffic/routing"
	"github.com/behavioral-ai/traffic/routing/routingtest"
)

func newRootEndpoint(pattern string) *rest.Endpoint {
	cache.ConstructorOverride(nil, cachetest.Exchange, operations.Serve)
	routing.ConstructorOverride(nil, routingtest.Exchange, operations.Serve)
	// Overriding agent http exchange
	//repository.Agent(intermediary.CacheNamespaceName).Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	//repository.Agent(intermediary.RoutingNamespaceName).Message(httpx.NewConfigExchangeMessage(routingtest.Exchange))

	return nil
	/*host.NewEndpoint(pattern,
	//[]any{link.Logger,
		repository.Agent(traffic.CacheNamespaceName),
		repository.Agent(traffic.LimiterNamespaceName),
		repository.Agent(traffic.RoutingNamespaceName)})

	*/
}
