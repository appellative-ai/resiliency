package operations

import (
	"github.com/behavioral-ai/collective/event"
	"github.com/behavioral-ai/core/messaging"
)

const (
	NamespaceName = "resiliency:agent/behavioral-ai/operations/operative"
)

// TODO : need host name
type agentT struct {
	running bool
	uri     string

	emissary   *messaging.Channel
	agents     *messaging.Exchange
	activity   messaging.ActivityFunc
	notifier   messaging.NotifyFunc
	dispatcher messaging.Dispatcher
}

// New - create a new operative
func New() messaging.Agent {
	return newAgent(nil, nil, nil)
}

func newAgent(activity messaging.ActivityFunc, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) *agentT {
	a := new(agentT)
	a.uri = NamespaceName

	a.agents = messaging.NewExchange()
	a.emissary = messaging.NewEmissaryChannel()
	a.activity = activity
	a.notifier = notifier
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
	a.emissary.C <- m
}

// Run - run the agent
func (a *agentT) Run() {
	if a.running {
		return
	}
	go emissaryAttend(a)
	a.running = true
}

// Shutdown - shutdown the agent
func (a *agentT) Shutdown() {
	if !a.emissary.IsClosed() {
		a.emissary.C <- messaging.Shutdown
	}
	a.running = false
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

func (a *agentT) addActivity(e messaging.ActivityItem) {
	if a.activity != nil {
		a.activity(e)
	} else {
		event.Agent.Message(messaging.NewActivityMessage(e))
	}
}

func (a *agentT) dispatch(channel any, event string) {
	messaging.Dispatch(a, a.dispatcher, channel, event)
}

func (a *agentT) finalize() {
	a.emissary.Close()
	a.agents.Shutdown()
}
