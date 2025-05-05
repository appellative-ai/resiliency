package operations

import (
	"errors"
	access "github.com/behavioral-ai/core/access2"
	"github.com/behavioral-ai/traffic/timeseries"
)

func originFromMap(m map[string]string) (o access.Origin, err error) {
	o.Region = m[timeseries.RegionKey]
	if o.Region == "" {
		return o, errors.New("invalid argument: origin region is empty")
	}
	o.Zone = m[timeseries.ZoneKey]
	if o.Zone == "" {
		return o, errors.New("invalid argument: origin zone is empty")
	}

	// SubZone is optional
	o.SubZone = m[timeseries.SubZoneKey]

	o.Host = m[timeseries.HostKey]
	if o.Host == "" {
		return o, errors.New("invalid argument: origin host is empty")
	}

	// InstanceId is optional
	o.InstanceId = m[timeseries.InstanceIdKey]
	return o, nil
}
