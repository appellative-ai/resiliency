package endpoint

import (
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/cache/cachetest"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/routing"
	"github.com/behavioral-ai/intermediary/routing/routingtest"
	urn2 "github.com/behavioral-ai/intermediary/urn"
	acc2 "github.com/behavioral-ai/resiliency/access"
	"github.com/behavioral-ai/traffic/limiter"
	"github.com/behavioral-ai/traffic/redirect"
	"github.com/behavioral-ai/traffic/urn"
)

func NewRootEndpoint() *rest.Endpoint {
	// Overriding cache agent http exchange
	_ = cache.Route
	cacheAgent := exchange.Agent(urn2.CacheAgent)
	cacheAgent.Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	m := make(map[string]string)
	m[config.CacheHostKey] = "localhost:8082"
	cacheAgent.Message(messaging.NewConfigMapMessage(m))

	// Overriding routing agent http exchange
	_ = routing.Route
	routingAgent := exchange.Agent(urn2.RoutingAgent)
	routingAgent.Message(httpx.NewConfigExchangeMessage(routingtest.Exchange))
	m[config.AppHostKey] = "localhost:8080"
	//m[config.TimeoutKey] = "10ms"
	routingAgent.Message(messaging.NewConfigMapMessage(m))

	_ = limiter.NamespaceName
	_ = redirect.NamespaceName

	return host.NewEndpoint(acc2.Agent,
		exchange.Agent(urn.RedirectAgent),
		exchange.Agent(urn2.CacheAgent),
		exchange.Agent(urn.LimiterAgent),
		exchange.Agent(urn2.RoutingAgent))
}
