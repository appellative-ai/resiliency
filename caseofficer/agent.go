package caseofficer

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
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
	Startup(m map[string]map[string]string) (*rest.Endpoint, []error)
}
