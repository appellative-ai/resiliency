package endpoint

import (
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/cache/cachetest"
	"github.com/behavioral-ai/intermediary/routing"
	"github.com/behavioral-ai/intermediary/routing/routingtest"
	"github.com/behavioral-ai/traffic/analytics"
	"github.com/behavioral-ai/traffic/config"
	"github.com/behavioral-ai/traffic/limiter"
	"github.com/behavioral-ai/traffic/redirect"
)

func NewRootEndpoint() host.ExchangeHandler {
	// overriding cache agent http exchange
	cache.Agent.Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	m := make(map[string]string)
	m[config.CacheHostKey] = "localhost:8082"
	cache.Agent.Message(messaging.NewConfigMapMessage(m))

	// overriding routing agent http exchange
	routing.Agent.Message(httpx.NewConfigExchangeMessage(routingtest.Exchange))
	m = make(map[string]string)
	m[config.AppHostKey] = "localhost:8080"
	//m[config.TimeoutKey] = "10ms"
	routing.Agent.Message(messaging.NewConfigMapMessage(m))

	chain := httpx.BuildChain(host.AccessLogLink, host.AuthorizationLink, redirect.Agent,
		analytics.Agent, cache.Agent, limiter.Agent, routing.Agent)

	return host.NewEndpoint(chain)
}
