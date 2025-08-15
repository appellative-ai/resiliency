package operations

import (
	"github.com/appellative-ai/collective/exchange"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/core/rest"
	"github.com/appellative-ai/traffic/cache"
	"github.com/appellative-ai/traffic/cache/cachetest"
	"github.com/appellative-ai/traffic/routing"
	"github.com/appellative-ai/traffic/routing/routingtest"
)

func setTestOverrides(a messaging.Agent) {
	// Cache overrides - global
	m := messaging.NewConfigMessage(rest.Exchange(cachetest.Exchange))
	m.AddTo(cache.AgentName)
	exchange.Message(m)

	// Cache overrides - Local
	a.Message(m)

	// Routing overrides - global
	m = messaging.NewConfigMessage(rest.Exchange(routingtest.Exchange))
	m.AddTo(routing.AgentName)
	exchange.Message(m)

	// Routing overrides - local
	a.Message(m)
}
