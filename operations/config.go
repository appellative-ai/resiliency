package operations

import (
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/messaging"
)

func configure(m *messaging.Message) (o access.Origin) {
	a := Agent
	cfg := messaging.ConfigMapContent(m)
	if cfg == nil {
		messaging.Reply(m, messaging.ConfigEmptyStatusError(a), a.Uri())
		return
	}
	o.Region = cfg[access.RegionKey]
	if o.Region == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, access.RegionKey), a.Uri())
		return
	}
	o.Zone = cfg[access.ZoneKey]
	if o.Zone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, access.ZoneKey), a.Uri())
		return
	}
	o.SubZone = cfg[access.SubZoneKey]
	if o.SubZone == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, access.SubZoneKey), a.Uri())
		return
	}
	o.Host = cfg[access.HostKey]
	if o.Host == "" {
		messaging.Reply(m, messaging.ConfigContentStatusError(a, access.HostKey), a.Uri())
		return
	}
	messaging.Reply(m, messaging.StatusOK(), a.Uri())
	return
}
