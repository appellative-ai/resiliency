package operations

import (
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
)

const (
	RegionKey     = "region"
	ZoneKey       = "zone"
	SubZoneKey    = "sub-zone"
	HostKey       = "host"
	InstanceIdKey = "id"
)

func newOriginFromMessage(agent messaging.Agent, m *messaging.Message) (o access.Origin, ok bool) {
	a := agent
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(a), a.Uri())
		return
	}
	region := cfg[RegionKey]
	if region == "" {
		return
	}
	o.Region = region
	o.Zone = cfg[ZoneKey]
	if o.Zone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, ZoneKey), a.Uri())
		return
	}
	o.SubZone = cfg[SubZoneKey]
	if o.SubZone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, SubZoneKey), a.Uri())
		return
	}
	o.Host = cfg[HostKey]
	if o.Host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, HostKey), a.Uri())
		return
	}
	o.InstanceId = cfg[InstanceIdKey]
	if o.Host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, InstanceIdKey), a.Uri())
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
	return o, true
}
