package endpoint

import (
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/rest"
	_ "github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/cache/cachetest"
	intermediary "github.com/behavioral-ai/intermediary/module"
	//_ "github.com/behavioral-ai/intermediary/routing"
	"github.com/behavioral-ai/intermediary/routing/routingtest"
	link "github.com/behavioral-ai/resiliency/link"
	//_ "github.com/behavioral-ai/traffic/limiter"
	traffic "github.com/behavioral-ai/traffic/module"
	//_ "github.com/behavioral-ai/traffic/redirect"
)

func newRootEndpoint() *rest.Endpoint {
	// Overriding agent http exchange
	repository.Agent(intermediary.CacheNamespaceName).Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	repository.Agent(intermediary.RoutingNamespaceName).Message(httpx.NewConfigExchangeMessage(routingtest.Exchange))

	return host.NewEndpoint(link.Logger,
		repository.Agent(traffic.RedirectNamespaceName),
		repository.Agent(intermediary.CacheNamespaceName),
		repository.Agent(traffic.LimiterNamespaceName),
		repository.Agent(intermediary.RoutingNamespaceName))
}
