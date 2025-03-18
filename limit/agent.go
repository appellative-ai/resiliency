package limit

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/common"
	"github.com/behavioral-ai/domain/metrics1"
	"github.com/behavioral-ai/domain/timeseries1"
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
)

type agentT struct {
	running bool
	uri     string
	traffic string
	origin  common.Origin

	ticker     *messaging.Ticker
	emissary   *messaging.Channel
	master     *messaging.Channel
	resolver   content.Resolution
	notifier   messaging.NotifyFunc
	dispatcher messaging.Dispatcher
	activity   messaging.ActivityFunc
}

func agentUri(origin common.Origin) string {
	return fmt.Sprintf("%v%v#%v", NamespaceName, strconv.Itoa(version), origin)
}

// New - create a new agent1 agent
func New(origin common.Origin, activity messaging.ActivityFunc, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, activity, notifier, dispatcher)
}

func newAgent(origin common.Origin, activity messaging.ActivityFunc, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) *agentT {
	a := new(agentT)
	a.origin = origin
	a.uri = agentUri(origin)

	a.ticker = messaging.NewTicker(messaging.Emissary, maxDuration)
	a.emissary = messaging.NewEmissaryChannel()
	a.master = messaging.NewMasterChannel()

	a.notifier = notifier
	a.activity = activity
	a.dispatcher = dispatcher
	return a
}

// String - identity
func (a *agentT) String() string { return a.Uri() }

// Uri - agent identifier
func (a *agentT) Uri() string { return a.uri }

// Name - agent urn
func (a *agentT) Name() string { return NamespaceName }

// Message - message the agent
func (a *agentT) Message(m *messaging.Message) {
	if m == nil {
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
	go emissaryAttend(a, timeseries1.Observations, content.Resolver, nil)
	a.running = true
}

// Shutdown - shutdown the agent
func (a *agentT) Shutdown() {
	if !a.emissary.IsClosed() {
		a.emissary.C <- messaging.Shutdown
	}
	if !a.master.IsClosed() {
		a.master.C <- messaging.Shutdown
	}
}

func (a *agentT) notify(e messaging.NotifyItem) {
	if e == nil {
		return
	}
	if a.notifier != nil {
		a.notifier(e)
	} else {
		event.Agent.Message(messaging.NewNotifyMessage(e))
	}
}

func (a *agentT) addActivity(e *messaging.ActivityItem) {
	if e == nil {
		return
	}
	if a.activity != nil {
		a.activity(*e)
	} else {
		event.Agent.Message(messaging.NewActivityMessage(*e))
	}
}

func (a *agentT) dispatch(channel any, event string) {
	messaging.Dispatch(a, a.dispatcher, channel, event)
}

func (a *agentT) reviseTicker(resolver *content.Resolution, s messaging.Spanner) {
	if s != nil {
		dur := s.Span()
		a.ticker.Start(dur)
		a.notify(messaging.NewStatusMessage(http.StatusOK, fmt.Sprintf("revised ticker -> traffic: %v duration: %v", a.traffic, dur), a.uri))
		return
	}
	p, status := content.Resolve[metrics1.TrafficProfile](metrics1.ProfileName, 1, resolver)
	if !status.OK() {
		a.ticker.Start(maxDuration)
		if status.NotFound() {
			status.SetAgent(a.Uri())
		}
		a.notify(status)
		return
	}
	traffic := p.Now()
	if p.IsMedium(traffic) || traffic == a.traffic {
		return
	}
	var dur time.Duration
	if p.IsLow(traffic) {
		dur = maxDuration
	} else {
		dur = minDuration
	}
	a.ticker.Start(dur)
	a.traffic = traffic
	a.notify(messaging.NewStatusMessage(http.StatusOK, fmt.Sprintf("revised ticker -> traffic: %v duration: %v", a.traffic, dur), a.uri))
}

func (a *agentT) emissaryFinalize() {
	a.emissary.Close()
	a.ticker.Stop()
}

func (a *agentT) masterFinalize() {
	a.master.Close()
}
