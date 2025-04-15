package operations

import (
	"github.com/behavioral-ai/collective/timeseries"
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
)

func newOriginFromMessage(agent messaging.Agent, m *messaging.Message) (o access.Origin, ok bool) {
	a := agent
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(a), a.Uri())
		return
	}
	region := cfg[timeseries.RegionKey]
	if region == "" {
		return
	}
	o.Region = region
	o.Zone = cfg[timeseries.ZoneKey]
	if o.Zone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, timeseries.ZoneKey), a.Uri())
		return
	}
	o.SubZone = cfg[timeseries.SubZoneKey]
	if o.SubZone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, timeseries.SubZoneKey), a.Uri())
		return
	}
	o.Host = cfg[timeseries.HostKey]
	if o.Host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, timeseries.HostKey), a.Uri())
		return
	}
	o.InstanceId = cfg[timeseries.InstanceIdKey]
	if o.Host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, timeseries.InstanceIdKey), a.Uri())
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
	return o, true
}
