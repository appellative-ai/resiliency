package caseofficer

import (
	"github.com/behavioral-ai/core/messaging"
)

type Agent interface {
	messaging.Agent
	BuildNetwork(m map[string]map[string]string) ([]any, []error)
}
