package endpoint

import (
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/intermediary/cache/cachetest"
	"github.com/behavioral-ai/intermediary/config"
	"github.com/behavioral-ai/intermediary/routing/routingtest"
	urn2 "github.com/behavioral-ai/intermediary/urn"
	"github.com/behavioral-ai/traffic/urn"
)

func NewRootEndpoint() *rest.Endpoint {
	// Overriding cache agent http exchange
	cache := exchange.Agent(urn2.CacheAgent)
	cache.Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	m := make(map[string]string)
	m[config.CacheHostKey] = "localhost:8082"
	cache.Message(messaging.NewConfigMapMessage(m))

	// Overriding routing agent http exchange
	routing := exchange.Agent(urn2.RoutingAgent)
	routing.Message(httpx.NewConfigExchangeMessage(routingtest.EchoExchange))
	m[config.AppHostKey] = "localhost:8080"
	//m[config.TimeoutKey] = "10ms"
	routing.Message(messaging.NewConfigMapMessage(m))

	return host.NewEndpoint(exchange.Agent(urn.RedirectAgent),
		exchange.Agent(urn2.CacheAgent),
		exchange.Agent(urn.LimiterAgent),
		exchange.Agent(urn2.RoutingAgent))
}
