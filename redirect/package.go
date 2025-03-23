package redirect

import "github.com/behavioral-ai/resiliency/common"

const (
	AppHostKey    = "app-host"
	RegionKey     = common.RegionKey
	ZoneKey       = common.ZoneKey
	SubZoneKey    = common.SubZoneKey
	HostKey       = common.HostKey
	InstanceIdKey = common.InstanceIdKey
)

var (
	Agent = New()
)
