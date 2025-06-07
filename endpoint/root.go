package endpoint

import (
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/cache/cachetest"
	"github.com/behavioral-ai/intermediary/routing"
	"github.com/behavioral-ai/intermediary/routing/routingtest"

	intermediary "github.com/behavioral-ai/intermediary/module"
	link "github.com/behavioral-ai/resiliency/link"
	traffic "github.com/behavioral-ai/traffic/module"
	//_ "github.com/behavioral-ai/traffic/limiter"
	//_ "github.com/behavioral-ai/intermediary/cache"
	//_ "github.com/behavioral-ai/traffic/redirect"
)

func newRootEndpoint() *rest.Endpoint {
	cache.ConstructorOverride(nil, cachetest.Exchange, operations.Serve)
	routing.ConstructorOverride(nil, routingtest.Exchange, operations.Serve)
	// Overriding agent http exchange
	//repository.Agent(intermediary.CacheNamespaceName).Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	//repository.Agent(intermediary.RoutingNamespaceName).Message(httpx.NewConfigExchangeMessage(routingtest.Exchange))

	return host.NewEndpoint([]any{link.Logger,
		repository.Agent(traffic.RedirectNamespaceName),
		repository.Agent(intermediary.CacheNamespaceName),
		repository.Agent(traffic.LimiterNamespaceName),
		repository.Agent(intermediary.RoutingNamespaceName)})
}
