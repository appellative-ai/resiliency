package redirect

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/httpx"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"net/http"
	"strconv"
	"time"
)

// Namespace ID Namespace Specific String
// NID + NSS
// NamespaceName
const (
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/redirect"
	minDuration   = time.Second * 10
	maxDuration   = time.Second * 15
	version       = 1
)

type agentT struct {
	running bool
	traffic string
	origin  common.Origin

	handler  messaging.Agent
	ticker   *messaging.Ticker
	emissary *messaging.Channel
	master   *messaging.Channel
}

func agentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", NamespaceName, strconv.Itoa(version), origin)
}

// New - create a new agent
func New() httpx.Agent {
	return newAgent(common.Origin{}, nil)
}

func newAgent(origin common.Origin, handler messaging.Agent) *agentT {
	a := new(agentT)
	a.origin = origin

	a.ticker = messaging.NewTicker(messaging.Emissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()
	a.master = messaging.NewMasterChannel()

	return a
}

// String - identity
func (a *agentT) String() string { return a.Uri() }

// Uri - agent identifier
func (a *agentT) Uri() string { return agentUri(a.origin) }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	if m.Event() == messaging.ConfigEvent {
		a.configure(m)
		return
	}
	if m.Event() == messaging.StartupEvent {
		a.run()
		return
	}
	if !a.running {
		return
	}
	switch m.Channel() {
	case messaging.Emissary:
		a.emissary.C <- m
	case messaging.Master:
		a.master.C <- m
	case messaging.Control:
		a.emissary.C <- m
		a.master.C <- m
	default:
		a.emissary.C <- m
	}
}

// Run - run the agent
func (a *agentT) run() {
	if a.running {
		return
	}
	go masterAttend(a, content.Resolver)
	go emissaryAttend(a, content.Resolver, nil)
	a.running = true
}

// Exchange - chainable exchange
func (a *agentT) Exchange(next httpx.Exchange) httpx.Exchange {
	return func(req *http.Request) (resp *http.Response, err error) {
		// TODO: if a redirect is configured, then process and ignore rest of pipeline
		if next != nil {
			resp, err = next(req)
		} else {
			resp = &http.Response{StatusCode: http.StatusOK}
		}
		return
	}
}

func (a *agentT) dispatch(channel any, event1 string) {
	a.handler.Message(event.NewDispatchMessage(a, channel, event1))
}

func (a *agentT) reviseTicker(resolver *content.Resolution, s messaging.Spanner) {
}

func (a *agentT) emissaryFinalize() {
	a.emissary.Close()
	a.ticker.Stop()
}

func (a *agentT) masterFinalize() {
	a.master.Close()
}

func (a *agentT) configure(m *messaging.Message) {
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(a), a.Uri())
		return
	}
	a.origin.Region = cfg[RegionKey]
	if a.origin.Region == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, RegionKey), a.Uri())
		return
	}
	a.origin.Zone = cfg[ZoneKey]
	if a.origin.Zone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, ZoneKey), a.Uri())
		return
	}
	a.origin.SubZone = cfg[SubZoneKey]
	if a.origin.SubZone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, SubZoneKey), a.Uri())
		return
	}
	a.origin.Host = cfg[HostKey]
	if a.origin.Host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, HostKey), a.Uri())
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
}
