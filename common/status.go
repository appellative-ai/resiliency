package common

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
)

func ConfigEmptyStatusError(agent messaging.Agent) *messaging.Status {
	return messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("config map is nil"), agent.Uri())
}

func ConfigContentStatusError(agent messaging.Agent, key string) *messaging.Status {
	return messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New(fmt.Sprintf("config map does not contain key: %v", key)), agent.Uri())
}
