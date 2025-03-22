package limit

import "github.com/behavioral-ai/resiliency/common"

// Configuration map keys

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
