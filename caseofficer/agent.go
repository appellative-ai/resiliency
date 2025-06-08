package caseofficer

import (
	"github.com/behavioral-ai/core/messaging"
)

const (
	LoggingRole       = "logging"
	AuthorizationRole = "authorization"
	CacheRole         = "role"
	RateLimiterRole   = "rate-limiter"
	RoutingRole       = "routing"
	RedirectRole      = "redirect"
	NameKey           = "name"
)

type Agent interface {
	messaging.Agent
	BuildNetwork(m map[string]map[string]string) ([]any, []error)
}
