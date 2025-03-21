package redirect

import "github.com/behavioral-ai/resiliency/common"

const (
	RegionKey  = common.RegionKey
	ZoneKey    = common.ZoneKey
	SubZoneKey = common.SubZoneKey
	HostKey    = common.HostKey
)

var (
	Agent = New()
)
