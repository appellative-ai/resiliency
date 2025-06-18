package caseofficer

import (
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/messaging"
)

const (
	LoggingRole       = "logging"
	AuthorizationRole = "authorization"
	CacheRole         = "cache"
	RateLimitingRole  = "rate-limiting"
	RoutingRole       = "routing"
	NameKey           = "name"
)

func init() {
	repository.RegisterConstructor(NamespaceNamePrimary, func() messaging.Agent {
		return NewPrimaryAgent(operations.Serve)
	})

}
