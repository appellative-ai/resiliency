package test

import (
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/intermediary/cache"
	"github.com/behavioral-ai/traffic/analytics"
	"github.com/behavioral-ai/traffic/limiter"
	"github.com/behavioral-ai/traffic/redirect"
)

func NewRootEndpoint() host.ExchangeHandler {

	chain := httpx.BuildChain(host.AccessLogLink, host.AuthorizationLink, redirect.Agent,
		analytics.Agent, cache.Agent, limiter.Agent, RoutingLink)

	return host.NewEndpoint(chain)
}
