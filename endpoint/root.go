package endpoint

import (
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/intermediary/cache/cachetest"
	"github.com/behavioral-ai/traffic/analytics"
	"github.com/behavioral-ai/traffic/config"
	"github.com/behavioral-ai/traffic/limiter"
	"github.com/behavioral-ai/traffic/redirect"
)

func NewRootEndpoint() host.ExchangeHandler {
	// configuring to not process upstream http requests for the cache and routing agents
	cache.Agent.Message(httpx.NewConfigExchangeMessage(cachetest.Exchange))
	m := make(map[string]string)
	m[config.CacheHostKey] = "localhost:8082"
	cache.Agent.Message(messaging.NewConfigMapMessage(m))

	chain := httpx.BuildChain(host.AccessLogLink, host.AuthorizationLink, redirect.Agent,
		analytics.Agent, cache.Agent, limiter.Agent, RoutingLink)

	return host.NewEndpoint(chain)
}
