package caseofficer

import (
	"github.com/behavioral-ai/collective/operations"
	"github.com/behavioral-ai/collective/repository"
	"github.com/behavioral-ai/core/messaging"
)

func init() {
	repository.RegisterConstructor(NamespaceNamePrimary, func() messaging.Agent {
		return NewPrimaryAgent(operations.Serve)
	})

}
