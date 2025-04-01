package operations

import (
	"github.com/behavioral-ai/collective/eventing"
	"github.com/behavioral-ai/core/messaging"
)

const (
	contentTypeNotifier = "application/notifier"
)

func newConfigNotifier(notifier eventing.NotifyFunc) *messaging.Message {
	m := messaging.NewMessage(messaging.Control, messaging.ConfigEvent)
	m.SetContent(contentTypeNotifier, notifier)
	return m
}

func configNotifierContent(m *messaging.Message) (eventing.NotifyFunc, bool) {
	if m.Event() != messaging.ConfigEvent || m.ContentType() != contentTypeNotifier {
		return nil, false
	}
	if cfg, ok := m.Body.(eventing.NotifyFunc); ok {
		return cfg, true
	}
	return nil, false
}
