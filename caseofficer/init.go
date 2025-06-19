package caseofficer

import (
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/core/messaging"
)

const (
	LoggingRole       = "logging"
	AuthorizationRole = "authorization"
	CacheRole         = "cache"
	RateLimitingRole  = "rate-limiting"
	RoutingRole       = "routing"
	NameKey           = "name"
	AssignmentKey     = "assignment"
	AssignmentGlobal  = "global"
	AssignmentLocal   = "local"
)

func init() {
	exchange.RegisterConstructor(NamespaceNamePrimary, func() messaging.Agent {
		return NewPrimaryAgent(operations.Serve)
	})
	exchange.RegisterConstructor(NamespaceNameSecondary, func() messaging.Agent {
		return NewSecondaryAgent(operations.Serve)
	})
}
