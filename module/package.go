package module

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/operations"
)

const (
	OperationsPath = "/operations"
)

func Startup(hostName string) {
	//test.Startup()
	operations.Message(messaging.StartupEvent)
}
