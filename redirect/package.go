package redirect

import "github.com/behavioral-ai/resiliency/common"

const (
	RegionKey     = common.RegionKey
	ZoneKey       = common.ZoneKey
	SubZoneKey    = common.SubZoneKey
	HostKey       = common.HostKey
	InstanceIdKey = common.InstanceIdKey
)

var (
	Agent = New()
)
