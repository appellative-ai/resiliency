package limit

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/event"
	http2 "github.com/behavioral-ai/core/http"
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
	NamespaceName = "resiliency:agent/behavioral-ai/resiliency/rate-limiting"
	minDuration   = time.Second * 10
	maxDuration   = time.Second * 15
	version       = 1
)

type agentT struct {
	running bool
	traffic string
	origin  common.Origin
	limiter *rateLimiter

	handler  messaging.Agent
	ticker   *messaging.Ticker
	emissary *messaging.Channel
	master   *messaging.Channel
}

func agentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", NamespaceName, strconv.Itoa(version), origin)
}

// New - create a new agent1 agent
func New() http2.Agent {
	return newAgent(common.Origin{}, nil)
}

func newAgent(origin common.Origin, handler messaging.Agent) *agentT {
	a := new(agentT)
	a.origin = origin
	if handler != nil {
		a.handler = handler
	} else {
		a.handler = event.Agent
	}
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
		if origin, ok := m.Body.(common.Origin); ok {
			a.origin = origin
		}
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
func (a *agentT) Run() {
	if a.running {
		return
	}
	go masterAttend(a, content.Resolver)
	go emissaryAttend(a, content.Resolver, nil)
	a.running = true
}

// Exchange - run the agent
func (a *agentT) Exchange(req *http.Request, next *http2.Frame) (resp *http.Response, err error) {
	// TODO: process rate limiting, and if not allowed, return
	if !a.limiter.Allow() {
		return &http.Response{StatusCode: http.StatusTooManyRequests}, nil
	}
	if next != nil {
		resp, err = next.Fn(req, next.Next)
		// TODO: need to update the response metrics
		a.Message(nil)
	} else {
		resp = &http.Response{StatusCode: http.StatusOK}
	}
	return
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
