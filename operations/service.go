package operations

import (
	"errors"
	"fmt"
	"github.com/behavioral-ai/collective/exchange"
	"github.com/behavioral-ai/core/host"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/rest"
	"github.com/behavioral-ai/resiliency/caseofficer"
	"github.com/behavioral-ai/resiliency/module"
	"github.com/behavioral-ai/traffic/cache"
	"github.com/behavioral-ai/traffic/cache/cachetest"
	"github.com/behavioral-ai/traffic/routing"
	"github.com/behavioral-ai/traffic/routing/routingtest"
	"net/http"
	"strings"
)

const (
	eventKey = "event"
)

type service struct {
	pattern string
}

func newServiceEndpoint(pattern string) rest.Endpoint {
	o := new(service)
	o.pattern = pattern
	return o
}

func (s *service) Pattern() string {
	return s.pattern
}

func (s *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.pattern) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	values := r.URL.Query()
	if len(values) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: no query args"))
		return
	}
	event := values.Get(eventKey)
	if event == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error: event query key not found"))
		return
	}
	err := validateEvent(event)
	if err == nil {
		exchange.Message(messaging.NewMessage(messaging.ChannelControl, event).AddTo(module.NamespaceNameOps))
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func validateEvent(event string) error {
	switch event {
	case messaging.StartupEvent, messaging.ShutdownEvent, messaging.PauseEvent, messaging.ResumeEvent:
		return nil
	default:
		return errors.New(fmt.Sprintf("invalid event: %v", event))
	}
}

const (
	primaryPattern   = "/primary"
	secondaryPattern = "/secondary"
)

func buildEndpoint(name string, chain []any) error {
	switch name {
	case caseofficer.NamespaceNamePrimary:
		// In testing, need to override the Exchange for cache and routing
		m := rest.NewExchangeMessage(cachetest.Exchange)
		m.AddTo(cache.NamespaceName)
		exchange.Message(m)

		m = rest.NewExchangeMessage(routingtest.Exchange)
		m.AddTo(routing.NamespaceName)
		exchange.Message(m)

		Endpoint[PrimaryEndpoint] = host.NewEndpoint(primaryPattern, chain)
	case caseofficer.NamespaceNameSecondary:
		// In testing, need to override the Exchange for routing
		m := rest.NewExchangeMessage(routingtest.Exchange)
		m.AddTo(routing.NamespaceName)
		m.SetCareOf(caseofficer.NamespaceNameSecondary)
		exchange.Message(m)

		Endpoint[SecondaryEndpoint] = host.NewEndpoint(secondaryPattern, chain)
	default:
		return errors.New(fmt.Sprintf("agent not found for name: %v", name))
	}
	return nil
}
