package test

import (
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/core/host"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/rest"
	"github.com/appellative-ai/traffic/cache"
	"github.com/appellative-ai/traffic/cache/cachetest"
	"github.com/appellative-ai/traffic/limiter"
)

func NewRootEndpoint() rest.Endpoint {
	//_ = operations.Agent
	_ = cache.NamespaceName
	_ = limiter.NamespaceName

	cache := exchange.Agent(cache.NamespaceName)
	cache.Message(rest.NewExchangeMessage(cachetest.Exchange))
	m := make(map[string]string)
	m["host"] = "localhost:8082"
	cache.Message(messaging.NewMapMessage(m))
	/*
		chain := httpx.BuildChain(host.AccessLogLink, host.AuthorizationLink,
			exchange.Agent(urn.RedirectAgent),
			exchange.Agent(urn2.CacheAgent),
			exchange.Agent(urn.LimiterAgent), RoutingLink)

	*/

	return host.NewEndpoint("", []any{ //repository.Agent(redirect.NamespaceName),
		cache, //repository.Agent(cache.Nurn2.CacheAgent),
		exchange.Agent(limiter.NamespaceName), RoutingLink})
}
